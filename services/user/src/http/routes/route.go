package routes

import (
	"strings"
	"time"

	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func SetRoutes(app *fiber.App) fiber.Router {
	if strings.ToUpper(config.MODE) == config.MODE_DEBUG {
		app.Get("/monitor", monitor.New(monitor.Config{
			Title:   "User Service Metrics",
			Refresh: 5 * time.Second,
			APIOnly: false,
		}))
	}

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	route := app.Group("/v1")

	return route
}
