package authJwt

import (
	"github.com/TimDebug/FitByte/src/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/samber/do/v2"
)

type JwtService struct {
}

func NewJwtService() JwtServiceInterface {
	return &JwtService{}
}

func NewJwtServiceInject(i do.Injector) (JwtServiceInterface, error) {
	return NewJwtService(), nil
}

func (js *JwtService) GenerateToken(userId string) (string, error) {
	claim := jwt.MapClaims{}
	claim["userId"] = userId

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	signedToken, err := token.SignedString([]byte(config.GetSecretKey()))
	if err != nil {
		return signedToken, err
	}

	return signedToken, nil
}
