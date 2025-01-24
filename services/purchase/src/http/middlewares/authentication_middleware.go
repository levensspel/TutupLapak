package middlewares

import (
	"errors"
	"fmt"
	"os"
	"strings"

	response "github.com/TimDebug/FitByte/src/model/web"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(c *fiber.Ctx) error {
	authorizationHeader := c.Get("Authorization")
	if !strings.Contains(authorizationHeader, "Bearer") && !strings.Contains(authorizationHeader, "bearer") {
		return c.Status(fiber.StatusUnauthorized).JSON(response.Web{
			Message: "UNAUTHORIZED",
			Data:    "Token Invalid",
		})
	}

	bearerToken := ""
	if strings.Contains(authorizationHeader, "Bearer") {
		bearerToken = strings.Replace(authorizationHeader, "Bearer ", "", -1)
	}
	if strings.Contains(authorizationHeader, "bearer") {
		bearerToken = strings.Replace(authorizationHeader, "bearer ", "", -1)
	}

	token, err := jwt.Parse(bearerToken, CheckTokenJWT)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(response.Web{
			Message: "UNAUTHORIZED",
			Data:    "Token Invalid",
		})
	}

	if !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(response.Web{
			Message: "UNAUTHORIZED",
			Data:    "Token Invalid",
		})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(response.Web{
			Message: "UNAUTHORIZED",
			Data:    "Token Invalid",
		})
	}

	userId := claims["userId"].(string)
	fmt.Printf("userId: %s", userId)
	c.Locals("userId", userId)

	return c.Next()
}

func CheckTokenJWT(t *jwt.Token) (interface{}, error) {
	if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
		return nil, errors.New("invalid token signing method")
	}

	return []byte(os.Getenv("JWT_SECRET_KEY")), nil
}
