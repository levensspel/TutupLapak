package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/TimDebug/TutupLapak/File/src/config"
	"github.com/rs/zerolog"
)

var (
	LogFile *os.File
	Logger  zerolog.Logger
	conf    *config.Configuration = config.GetConfig()
)

func Init() error {
	if conf == nil {
		return fmt.Errorf("unable to open the env file")
	}
	err := initMultiWriter()
	if err != nil {
		return err
	}
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	if conf.IsProduction {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	}
	return nil
}

func initMultiWriter() error {
	if _, err := os.Stat("log"); os.IsNotExist(err) {
		err := os.Mkdir("log", os.ModePerm)
		if err != nil {
			return err
		}
	}
	LogFile, err := os.OpenFile("log/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	multi := zerolog.MultiLevelWriter(LogFile, os.Stdout)
	if conf.IsProduction {
		multi = zerolog.MultiLevelWriter(LogFile)
	}
	Logger = zerolog.New(multi).With().Timestamp().Logger()
	return nil
}

func Cleanup() {
	if LogFile != nil {
		LogFile.Close()
	}
}
