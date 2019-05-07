package google_grpc

import (
	"github.com/hhy5861/go-common/common"
	"github.com/hhy5861/go-common/logger"
	"google.golang.org/grpc"
	"sync"
)

type (
	GrpcClientService struct {
		cfg   *GrpcConfig
		tools *common.Tools
	}

	GrpcConfig struct {
		Cfg map[string][]string
	}
)

var connMap sync.Map

func GetGrpcConn(keys ...string) *grpc.ClientConn {
	key := "default"
	if len(keys) >= 1 {
		key = keys[0]
	}

	conn, ok := connMap.Load(key)
	if ok {
		return conn.(*grpc.ClientConn)
	}

	return nil
}

func NewGrpcClientService(cfg *GrpcConfig) *GrpcClientService {
	return &GrpcClientService{
		cfg:   cfg,
		tools: common.NewTools(),
	}
}

func (svc *GrpcClientService) GrpcDial() {
	for key, value := range svc.cfg.Cfg {
		num := svc.tools.GenerateRangeNum(0, len(value))
		conn, err := grpc.Dial(value[num], grpc.WithInsecure())
		if err != nil {
			logger.Fatal(err)
		}

		connMap.Store(key, conn)
	}
}
