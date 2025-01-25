package logger

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/TimDebug/TutupLapak/File/src/config"
	"github.com/gofiber/fiber/v2"
)

func ZerologWriter(ctx *fiber.Ctx, logString []byte) {
	logLine := string(logString)
	parts := strings.Split(logLine, "|")
	if len(parts) != 7 {
		Logger.Error().Str("raw_log", logLine).Msg("invalid log format")
	}

	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}

	status, err := strconv.Atoi(parts[1])
	if err != nil {
		Logger.Error().Str("raw_log", logLine).Err(err).Msg("invalid status code")
	}

	entryLog := Logger.Info()
	if status != 200 {
		entryLog = Logger.Error().Err(fmt.Errorf(parts[6]))
	}
	entryLog.
		Str("time", parts[0]).
		Int("status", status).
		Str("latency", parts[2]).
		Str("ip", parts[3]).
		Str("method", parts[4]).
		Str("path", parts[5]).
		Str("request_id", ctx.Locals(config.RequestKey).(string)).
		Msg("parsed log")
}
