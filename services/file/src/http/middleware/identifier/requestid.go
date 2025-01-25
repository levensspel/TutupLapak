package identifier

import (
	"github.com/TimDebug/TutupLapak/File/src/config"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func RequestID(c *fiber.Ctx) error {
	requestID := uuid.NewString()
	c.Locals(config.RequestKey, requestID)
	return c.Next()
}
