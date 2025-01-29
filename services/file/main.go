package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"

	"github.com/TimDebug/TutupLapak/File/src/config"
	"github.com/TimDebug/TutupLapak/File/src/database/migrations"
	"github.com/TimDebug/TutupLapak/File/src/grpc"
	httpServer "github.com/TimDebug/TutupLapak/File/src/http"
	log "github.com/TimDebug/TutupLapak/File/src/logger"
)

func main() {
	var wg sync.WaitGroup

	// Initialize logger
	err := log.Init()
	if err != nil {
		log.Logger.Fatal().Err(err).Msg("unable to init basic logger")
	}
	log.Logger.Info().Msg("configured basic logger")
	log.Logger.Info().Msg(fmt.Sprintf("NumCPU: %d", runtime.NumCPU()))

	// Initialize app configurations
	err = config.New()
	if err != nil {
		log.Logger.Fatal().Err(err).Msg("failed to load the env file")
	}
	log.Logger.Info().Msg("loaded basic configuration")

	// Reinitialize logger based on app config
	log.Add(config.Config)
	log.Logger.Info().Msg("configured logger from configuration")

	// Auto migrate
	err = migrations.Migrate()
	if err != nil {
		log.Logger.Fatal().Err(err).Msg("unable to run migration files")
	}
	log.Logger.Info().Msg("successfully run migration files")

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go func() {
		oscall := <-quit
		log.Logger.Warn().Msgf("system call:%+v", oscall)
		log.Cleanup()
		os.Exit(0)
	}()

	// Run http server
	wg.Add(1)
	go func() {
		defer wg.Done()
		server := httpServer.HttpServer{}
		server.Listen()
	}()

	// run grpc server
	wg.Add(1)
	go func() {
		defer wg.Done()
		server := grpc.GrpcServer{}
		server.Listen()
	}()

	wg.Wait()
}
