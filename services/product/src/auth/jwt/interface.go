package authJwt

import "github.com/golang-jwt/jwt/v5"

type JwtServiceInterface interface {
	GenerateToken(userId string) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}
