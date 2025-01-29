package main

import (
	"context"
	"fmt"
	"time"

	"github.com/TimDebug/TutupLapak/File/src/grpc/proto/model/file"
	"github.com/TimDebug/TutupLapak/File/src/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient(
		"localhost:5001",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Logger.Fatal().Err(err).Msg("fatal initiating new grpc client")
		return
	}
	defer conn.Close()

	fileClient := file.NewFileServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := fileClient.CheckExist(ctx, &file.FileRequest{FileId: "5a54fd34-57c9-42b4-b439-ee16903d68fc"})
	if err != nil {
		logger.Logger.Error().Err(err).Msg("check procedure returns error")
		return
	}
	fmt.Printf("success: %+v\n", response)
}
