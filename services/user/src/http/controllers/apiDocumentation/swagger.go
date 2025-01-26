package swaggerRoutes

import (
	_ "github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/docs"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func SetRouteSwagger(router fiber.Router) {
	router.Get("/docs/swagger/*", swagger.HandlerDefault)
}
