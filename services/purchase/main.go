package main

import (
	"fmt"

	"github.com/TimDebug/FitByte/src/database/migrations"
	"github.com/TimDebug/FitByte/src/di"
	httpServer "github.com/TimDebug/FitByte/src/http"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	fmt.Printf("Load ENV\n")
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	fmt.Printf("DI Healthcheck\n")
	di.HealthCheck()

	//? Auto Migrate
	fmt.Printf("Migrate\n")
	migrations.Migrate()

	fmt.Printf("Start Server\n")
	server := httpServer.HttpServer{}
	server.Listen()

}
