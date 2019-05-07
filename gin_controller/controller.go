package gin_controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hhy5861/go-common/common"
	"github.com/hhy5861/go-common/jwt"
	"github.com/hhy5861/go-common/logger"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"net/http"
	"strconv"
)

type (
	Controller struct {
		JwtConf   *jwt.JwtConfig
		localizer *i18n.Localizer
	}

	Response struct {
		Success int         `json:"success"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}
)

var (
	Tools *common.Tools
)

func init() {
	Tools = common.NewTools()
}

func NewController(ctl *Controller) *Controller {
	return ctl
}

func (ctl *Controller) ResponseList(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, data)
}

func (ctl *Controller) Response(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, Response{
		Success: 0,
		Message: "success",
		Data:    data,
	})
}

func (ctl *Controller) ParamsExecption(ctx *gin.Context, err error) {
	logger.Error(err)

	ctx.JSON(http.StatusOK, Response{
		Success: -100000,
		Message: ctl.getLocalize(ctx).i18nLocalize(-100000),
		Data:    nil,
	})
}

func (ctl *Controller) ServiceExecption(ctx *gin.Context, err error) {
	logger.Error(err)

	ctx.JSON(http.StatusOK, Response{
		Success: -100001,
		Message: ctl.getLocalize(ctx).i18nLocalize(-100001),
		Data:    nil,
	})
}

func (ctl *Controller) ServiceCodeExecption(ctx *gin.Context, code int, err error) {
	logger.Error(err)

	ctx.JSON(http.StatusOK, Response{
		Success: code,
		Message: ctl.getLocalize(ctx).i18nLocalize(code),
		Data:    nil,
	})
}

func (ctl *Controller) UnauthorizedExecption(ctx *gin.Context) {
	ctx.JSON(http.StatusUnauthorized, Response{
		Success: http.StatusUnauthorized,
		Message: http.StatusText(http.StatusUnauthorized),
		Data:    nil,
	})
}

func (ctl *Controller) Health(ctx *gin.Context) {
	ctl.Response(ctx, nil)
}

func (ctl *Controller) i18nLocalize(code int) string {

	return ctl.localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: strconv.Itoa(code),
		DefaultMessage: &i18n.Message{
			ID:    strconv.Itoa(code),
			Other: "Internal error in the service",
		},
	})
}

func (ctl *Controller) getLocalize(ctx *gin.Context) *Controller {
	localizer, ok := ctx.Get("Localizer")
	if ok && localizer != nil {
		ctl.localizer = localizer.(*i18n.Localizer)
	}

	return ctl
}
