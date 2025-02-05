package grpc

import (
	"log"
	"net"

	"github.com/TIM-DEBUG-ProjectSprintBatch3/go-fiber-template/src/config"
	"google.golang.org/grpc"
)

type GrpcServer struct{}

func (gs *GrpcServer) Listen() {
	// create tcp server
	lis, err := net.Listen("tcp", ":"+config.GetPortGrpc())

	if err != nil {
		log.Fatal("Failed to Listen on Port: %v", err)
	}

	// create new grpc server handler
	server := grpc.NewServer()

	// run server
	if err := server.Serve(lis); err != nil {
		log.Fatal(err.Error())
	}
}
