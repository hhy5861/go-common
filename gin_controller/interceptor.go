package gin_controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hhy5861/go-common/jwt"
	"github.com/hhy5861/go-common/logger"
)

const (
	tokenName          = "x-auth-token"
	StandardClaimsName = "standardClaims"
)

func (ctl *Controller) JwtInterceptor(ctx *gin.Context) {
	token := ctx.GetHeader(tokenName)
	claims, err := jwt.NewJwtPackage(ctl.JwtConf).ParseWithClaims(token)
	if err != nil {
		logger.Error(err)
		ctl.UnauthorizedExecption(ctx)
		ctx.Abort()
		return
	}

	ctx.Set(StandardClaimsName, claims)
}
