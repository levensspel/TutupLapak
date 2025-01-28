package userGrpc

import (
	"fmt"
	"log"
	"net"

	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/config"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/di"
	protoUserController "github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/http/controllers/user/proto"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/services/proto/user"
	"github.com/samber/do/v2"
	"google.golang.org/grpc"
)

type UserGrpcServer struct {
}

func StartGrpcServer() {
	// Start gRPC server
	go func() {
		_PORT := config.GetGRPCPort()

		lis, err := net.Listen("tcp", fmt.Sprintf(":%s", _PORT))
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}

		grpcServer := grpc.NewServer()

		puc := do.MustInvoke[*protoUserController.ProtoUserController](di.Injector)
		user.RegisterUserServiceServer(grpcServer, puc)

		fmt.Printf("> gRPC server listening on :%s", _PORT)

		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()
}
