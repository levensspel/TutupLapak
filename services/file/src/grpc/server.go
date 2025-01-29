package grpc

import (
	"context"
	"net"

	"github.com/TimDebug/TutupLapak/File/src/database/postgres"
	"github.com/TimDebug/TutupLapak/File/src/grpc/proto/model/file"
	"github.com/TimDebug/TutupLapak/File/src/logger"
	"github.com/TimDebug/TutupLapak/File/src/repo"
	"google.golang.org/grpc"
)

type GrpcServer struct{}

func (g *GrpcServer) Listen() error {
	lis, err := net.Listen("tcp", ":5001")
	if err != nil {
		logger.Logger.Fatal().Err(err).Msg("fatal error")
	}

	db, err := postgres.NewPgxConnect()
	if err != nil {
		logger.Logger.Fatal().Err(err).Msg("unable to establish database connection")
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
