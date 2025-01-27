package protoUserController

import (
	"context"

	loggerZap "github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/logger/zap"
	fileService "github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/services/external/file"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/services/proto/user"
	userService "github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/services/user"
	"github.com/samber/do/v2"
)

type ProtoUserController struct {
	userService userService.UserServiceInterface
	fileService fileService.FileServiceInterface
	logger      loggerZap.LoggerInterface

	// Embed UnimplementedUserServiceServer to satisfy gRPC interface
	user.UnimplementedUserServiceServer
}

func NewUserController(userService userService.UserServiceInterface, fileService fileService.FileServiceInterface, logger loggerZap.LoggerInterface) ProtoUserControllerInterface {
	return &ProtoUserController{userService: userService, logger: logger}
	// return ProtoUserController{userService: userService, logger: logger}
}

func NewUserControllerInject(i do.Injector) (*ProtoUserController, error) {
	_userService := do.MustInvoke[userService.UserServiceInterface](i)
	_fileService := do.MustInvoke[fileService.FileServiceInterface](i)
	_logger := do.MustInvoke[loggerZap.LoggerInterface](i)
	// kalau pakai interface, dibawah ini
	// return NewUserController(_userService, _fileService, _logger), nil

	return &ProtoUserController{
		userService: _userService,
		fileService: _fileService,
		logger:      _logger,
	}, nil
}

// GetUserDetails implements user.UserServiceServer.
func (puc ProtoUserController) GetUserDetails(ctx context.Context, request *user.UserRequest) (*user.UsersResponse, error) {
	result := user.UsersResponse{}

	return &result, nil
}

// mustEmbedUnimplementedUserServiceServer implements user.UserServiceServer.
func (puc ProtoUserController) mustEmbedUnimplementedUserServiceServer() {
	return
}
