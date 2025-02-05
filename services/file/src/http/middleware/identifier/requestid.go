package identifier

import (
	"time"

	"github.com/TimDebug/TutupLapak/File/src/config"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func RequestID(c *fiber.Ctx) error {
	requestID := uuid.NewString()
	c.Locals(config.RequestKey, requestID)
	c.Locals(config.ElapsedKey, time.Now())
	return c.Next()
}
