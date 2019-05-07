package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/hhy5861/go-common/common"
	"github.com/hhy5861/go-common/gin_controller"
	"github.com/hhy5861/go-common/gin_validator"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/soheilhy/cmux"
	"golang.org/x/sync/errgroup"
	"golang.org/x/text/language"
	"net"
	"net/http"
	"reflect"
)

type (
	GinServer struct {
		ServerCfg     *Config
		Engine        *gin.Engine
		tools         *common.Tools
		RouterGroup   *gin.RouterGroup
		I18nBundle    *i18n.Bundle
		LoopCall      func(structs ...interface{})
		RegisterRoute func()
		httpListener  net.Listener
		GrpcServer    *GrpcServer
	}

	Config struct {
		ContextPath string `json:"contextPath" yaml:"contextPath"`
		Host        string `json:"host" yaml:"host"`
		Port        int    `json:"port" yaml:"port"`
		Mode        string `json:"mode" yaml:"mode"`
		GrpcEnabled bool   `json:"grpcEnabled" yaml:"grpcEnabled"`
	}
)

const (
	Localizer = "Localizer"
)

func NewGinServer(cfg *Config) *GinServer {
	gin.SetMode(cfg.Mode)
	binding.Validator = new(gin_validator.DefaultValidator)

	return &GinServer{
		tools:      common.NewTools(),
		ServerCfg:  cfg,
		Engine:     gin.Default(),
		I18nBundle: &i18n.Bundle{DefaultLanguage: language.English},
		GrpcServer: NewGrpcServer(),
		LoopCall: func(structs ...interface{}) {
			for _, v := range structs {
				classType := reflect.TypeOf(v)
				classValue := reflect.ValueOf(v)

				for i := 0; i < classType.NumMethod(); i++ {
					m := classValue.MethodByName(classType.Method(i).Name)
					if m.IsValid() {
						var params []reflect.Value
						m.Call(params)
					}
				}
			}
		},
	}
}

func (svc *GinServer) Run() {
	addr := fmt.Sprintf("%s:%d", svc.ServerCfg.Host, svc.ServerCfg.Port)
	fmt.Println("Listening on: ", addr)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	m := cmux.New(listener)
	g := new(errgroup.Group)

	if svc.ServerCfg.GrpcEnabled {
		svc.GrpcServer.Listener = m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
		g.Go(func() error { return svc.GrpcServer.RunGrpcServe() })
	}

	svc.httpListener = m.Match(cmux.HTTP1Fast())
	svc.LoopCall()

	g.Go(func() error { return svc.RunHttpServe() })
	g.Go(func() error { return m.Serve() })

	fmt.Println("run server: ", g.Wait())
}

func (svc *GinServer) RunHttpServe() error {
	svc.newRoute().RegisterRoute()

	s := &http.Server{Handler: svc.Engine}
	return s.Serve(svc.httpListener)
}

func (svc *GinServer) newRoute() *GinServer {
	svc.RouterGroup = svc.Engine.Group(svc.ServerCfg.ContextPath, func(ctx *gin.Context) {
		localizer := i18n.NewLocalizer(svc.I18nBundle, ctx.Request.FormValue("lang"), ctx.GetHeader("Accept-Language"))

		ctx.Set(Localizer, localizer)
		ctx.Next()
	})

	svc.RouterGroup.GET("/health", new(gin_controller.Controller).Health)
	return svc
}
