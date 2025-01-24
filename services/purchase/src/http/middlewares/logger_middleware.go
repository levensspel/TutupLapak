package middlewares

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

// LoggerMiddleware adalah middleware untuk logging
func LoggerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Catat waktu sebelum request diproses
		start := time.Now()

		// Lanjutkan ke handler berikutnya
		err := c.Next()

		// Catat waktu setelah request diproses
		duration := time.Since(start)

		// Log detail request dan respons
		log.Printf(
			"%s - %s %s %d - %s",
			c.IP(),                    // IP Address
			c.Method(),                // HTTP Method
			c.Path(),                  // Request Path
			c.Response().StatusCode(), // HTTP Status Code
			duration,                  // Durasi pemrosesan
		)

		return err
	}
}
