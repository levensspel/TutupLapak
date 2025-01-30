package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/TimDebug/TutupLapak/File/src/config"
	"github.com/TimDebug/TutupLapak/File/src/grpc/proto/model/file"
	"github.com/TimDebug/TutupLapak/File/src/logger"
	"github.com/TimDebug/TutupLapak/File/src/repo"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

var (
	appConfig *config.Configuration = config.GetConfig()
)

type GrpcServer struct {
	DB *pgxpool.Pool
}

func (g *GrpcServer) Listen() error {
	serverConn := fmt.Sprintf("%s:%s", "0.0.0.0", appConfig.GRPCPort)
	lis, err := net.Listen("tcp", serverConn)
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	fileRepo := repo.NewFileRepository(g.DB)
	fileService := NewFileService(&fileRepo)
	file.RegisterFileServiceServer(server, fileService)

	go func() {
		<-context.Background().Done()
		server.GracefulStop()
	}()

	logger.Logger.Info().Msg(fmt.Sprintf("grpc server listens to: %s", serverConn))
	return server.Serve(lis)
}
