package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/hhy5861/go-common/common"
)

type (
	JwtPackage struct {
		key     []byte
		expired int64
	}

	StandardClaims struct {
		Id      uint   `json:"id"`
		Uuid    string `json:"userUuid"`
		UnionId string `json:"unionId"`
		OpenId  string `json:"openId"`
		jwt.StandardClaims
	}

	JwtConfig struct {
		Secret  string `yaml:"secret"`
		Expired int64  `yaml:"expired"`
	}
)

var (
	tools *common.Tools
)

func init() {
	tools = common.NewTools()
}

func NewJwtPackage(config *JwtConfig) *JwtPackage {
	proKey := jwt.EncodeSegment([]byte(config.Secret))

	return &JwtPackage{key: []byte(proKey), expired: config.Expired}
}

func (pkg *JwtPackage) CreateToken(standardClaims *StandardClaims) (string, error) {
	standardClaims.Issuer = standardClaims.Uuid
	standardClaims.Subject = standardClaims.UnionId
	standardClaims.IssuedAt = tools.GetNowMillisecond()
	standardClaims.ExpiresAt = tools.GetNowMillisecond() + (pkg.expired * 1000)

	return jwt.NewWithClaims(jwt.SigningMethodHS256, standardClaims).SignedString(pkg.key)
}

func (pkg *JwtPackage) ParseWithClaims(jwtToken string) (*StandardClaims, error) {
	var standardClaims StandardClaims

	claims, err := jwt.ParseWithClaims(jwtToken, &standardClaims, func(token *jwt.Token) (interface{}, error) {
		return pkg.key, nil
	})

	if err != nil {
		return nil, err
	}

	return claims.Claims.(*StandardClaims), nil
}
