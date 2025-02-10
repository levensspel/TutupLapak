package main

import (
	"fmt"

	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/config"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/database/migrations"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/di"
	userGrpc "github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/grpc"
	httpServer "github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/http"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	fmt.Printf("Load ENV\n")
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	err = config.SetupReusableEnv()
	if err != nil {
		panic(err)
	}

	fmt.Printf("DI Healthcheck\n")
	di.HealthCheck()

	//? Auto Migrate
	fmt.Printf("Migrate\n")
	migrations.Migrate()

	fmt.Printf("Start gRPC Server\n")
	userGrpc.StartGrpcServer()

	fmt.Printf("Start Server\n")
	server := httpServer.HttpServer{}
	server.Listen()

}
