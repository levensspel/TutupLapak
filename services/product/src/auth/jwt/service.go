package authJwt

import (
	"github.com/TIM-DEBUG-ProjectSprintBatch3/go-fiber-template/src/config"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/go-fiber-template/src/exceptions"
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

func (js *JwtService) ValidateToken(token string) (*jwt.Token, error) {

	jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, exceptions.NewUnauthorizedError("Invalid token signin method")
		}
		return []byte(config.GetSecretKey()), nil
	})

	if err != nil {
		return nil, err
	}

	return jwtToken, nil

}
