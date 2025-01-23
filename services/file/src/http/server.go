package httpServer

import (
	"fmt"

	"github.com/TimDebug/TutupLapak/File/src/config"
	"github.com/TimDebug/TutupLapak/File/src/database/postgres"
	"github.com/TimDebug/TutupLapak/File/src/http/middlewares"
	log "github.com/TimDebug/TutupLapak/File/src/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type HttpServer struct{}

func (s *HttpServer) Listen() {
	app := fiber.New(fiber.Config{
		ServerHeader: "TIM-DEBUG",
		ErrorHandler: middlewares.ErrorHandler,
	})
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowCredentials: false,
		AllowMethods:     "POST",
	}))

	db, err := postgres.NewPgxConnect()
	if err != nil {
		log.Logger.Fatal().Err(err).Msg("unable to establish database connection")
	}
	appConfig := config.Config
	var storageClient StorageClient
	if appConfig.IsProduction {
		storageClient = &S3StorageClient{Config: appConfig}
	} else {
		storageClient = nil
	}
	repo := NewFileRepository(db)
	service := NewFileService(repo, db, storageClient)
	controller := NewFileController(service)

	routes := app.Group("/v1")
	routes.Post("/file", controller.Upload)

	app.Listen(fmt.Sprintf("%s:%s", "0.0.0.0", config.GetPort()))
}
