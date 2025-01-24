package swaggerRoutes

import (
	_ "github.com/TimDebug/FitByte/docs"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func SetRouteSwagger(router fiber.Router) {
	router.Get("/docs/swagger/*", swagger.HandlerDefault)
}
