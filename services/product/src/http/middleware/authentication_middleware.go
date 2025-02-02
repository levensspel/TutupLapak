package middleware

import (
	"strings"

	authJwt "github.com/TIM-DEBUG-ProjectSprintBatch3/go-fiber-template/src/auth/jwt"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/go-fiber-template/src/di"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/go-fiber-template/src/model/dtos/response"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/samber/do/v2"
)

func AuthMiddleware(c *fiber.Ctx) error {
	tokenStr := c.Get("Authorization")
	tokenStr = strings.Replace(tokenStr, "Bearer ", "", 1)
	if tokenStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(response.Web{
			Message: "UNAUTHORIZED",
		})
	}

	jwtService := do.MustInvoke[authJwt.JwtServiceInterface](di.Injector)

	token, err := jwtService.ValidateToken(tokenStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(response.Web{
			Message: err.Error(),
		})
	}

	if !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(response.Web{
			Message: "UNAUTHORIZED",
		})
	}

	c.Locals("userId", token.Claims.(jwt.MapClaims)["userId"])
	return c.Next()
}
