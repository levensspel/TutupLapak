package logger

import (
	"time"

	"github.com/TimDebug/TutupLapak/File/src/config"
	"github.com/gofiber/fiber/v2"
)

func ResponseLogger(ctx *fiber.Ctx) error {
	err := ctx.Next()
	requestID, _ := ctx.Locals(config.RequestKey).(string)
	startTime, ok := ctx.Locals(config.ElapsedKey).(time.Time)
	var elapsed time.Duration = 0
	if ok {
		elapsed = time.Since(startTime)
	}
	logger := Logger.Info().
		Str("request_id", requestID).
		Str("method", ctx.Method()).
		Str("path", ctx.Path()).
		Int("status", ctx.Response().StatusCode()).
		Dur("elapsed", elapsed)
	if err != nil {
		logger.Err(err).Msg("request completed with error")
	} else {
		logger.Msg("request processed")
	}
	return err
}
