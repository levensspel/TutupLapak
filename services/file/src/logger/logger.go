package logger

import (
	"io"
	"os"
	"time"

	"github.com/TimDebug/TutupLapak/File/src/config"
	"github.com/rs/zerolog"
)

var (
	logFile *os.File
	Logger  zerolog.Logger
)

func Init() error {
	err := initMultiWriter()
	if err != nil {
		return err
	}
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	return nil
}

func Add(conf config.Configuration) {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if conf.IsProduction {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	} else {
		zerolog.TimeFieldFormat = time.RFC3339
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}

func initMultiWriter() error {
	logFile, err := os.OpenFile("log/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	multi := io.MultiWriter(logFile, os.Stdout)
	Logger = zerolog.New(multi).With().Timestamp().Logger()
	return nil
}

func Cleanup() {
	if logFile != nil {
		logFile.Close()
	}
}
