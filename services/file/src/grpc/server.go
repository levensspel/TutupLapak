package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/TimDebug/TutupLapak/File/src/config"
	"github.com/TimDebug/TutupLapak/File/src/database/postgres"
	"github.com/TimDebug/TutupLapak/File/src/grpc/proto/model/file"
	"github.com/TimDebug/TutupLapak/File/src/repo"
	"google.golang.org/grpc"
)

var (
	appConfig config.Configuration = config.Config
)

type GrpcServer struct{}

func (g *GrpcServer) Listen() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", appConfig.GRPCPort))
	if err != nil {
		return err
	}

	db, err := postgres.NewPgxConnect()
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	fileRepo := repo.NewFileRepository(db)
	fileService := NewFileService(&fileRepo)
	file.RegisterFileServiceServer(server, fileService)

	go func() {
		<-context.Background().Done()
		server.GracefulStop()
	}()

	return server.Serve(lis)
}
