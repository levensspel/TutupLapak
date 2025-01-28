package purchaseGrpc

import (
	"fmt"

	"github.com/TimDebug/FitByte/src/config"
	functionCallerInfo "github.com/TimDebug/FitByte/src/logger/helper"
	loggerZap "github.com/TimDebug/FitByte/src/logger/zap"
	"github.com/TimDebug/FitByte/src/services/proto/user"
	"github.com/samber/do/v2"
	"google.golang.org/grpc"
)

type ProtoUserController struct {
	logger      loggerZap.LoggerInterface
	UserService user.UserServiceClient
}

func NewGRPCClientInject(i do.Injector) (*ProtoUserController, error) {
	_logger := do.MustInvoke[loggerZap.LoggerInterface](i)

	_userServiceAddress := fmt.Sprintf("%s:%s", config.GetUserGRPCHost(), config.GetUserGRPCPort()) // known as 50051

	connection, err := grpc.Dial(_userServiceAddress, grpc.WithInsecure())
	if err != nil {
		_logger.Error(err.Error(), functionCallerInfo.PurhcaseControllerPutCart, "Grpc Connection")
		fmt.Printf("Grpc Connection Error: %s\n", err.Error())
		return nil, err
	}
	// defer {
	// 	connection.Close()
	// }

	// Create new userService client
	_userServiceClient := user.NewUserServiceClient(connection)

	fmt.Printf("GRPC Client>> Listening to %s\n", _userServiceAddress)
	return &ProtoUserController{
		logger:      _logger,
		UserService: _userServiceClient,
	}, nil
}
