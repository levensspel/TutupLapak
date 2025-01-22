package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/TimDebug/TutupLapak/File/src/database/migrations"
	"github.com/TimDebug/TutupLapak/File/src/di"
	httpServer "github.com/TimDebug/TutupLapak/File/src/http"
	loggerZap "github.com/TimDebug/TutupLapak/File/src/logger/zap"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/samber/do/v2"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	di.HealthCheck()

	log := do.MustInvoke[loggerZap.LoggerInterface](di.Injector)
	log.Info("Done loading env and dependency injection", "main", nil)

	// Handle graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sig
		di.Injector.Shutdown()
		os.Exit(0)
	}()

	//? Auto Migrate
	migrations.Migrate()
	log.Info("Done database migration", "main", nil)

	server := httpServer.HttpServer{}
	server.Listen()
}
