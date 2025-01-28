package protoUserController

import (
	"context"
	"fmt"

	functionCallerInfo "github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/logger/helper"
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

func NewProtoUserControllerInject(i do.Injector) (*ProtoUserController, error) {
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

	// todo; check cache ada profile user tidak?
	// belum tau kalian mau pakai cache middleware atau internal jadi saya kasih catatan saja-ad1ee

	response, err := puc.userService.GetUserProfiles(ctx, request.UserIds)
	if err != nil {
		fmt.Sprintf("Error %s", err.Error())
		puc.logger.Error(err.Error(), functionCallerInfo.UserControllerGetUserProfile)
		return nil, err
	}

	// Populate the response with the data from user profiles.
	for _, profile := range response {
		userResponse := &user.UserResponse{
			Email:             profile.Email,
			Phone:             profile.Phone,
			BankAccountName:   profile.BankAccountName,
			BankAccountHolder: profile.BankAccountHolder,
			BankAccountNumber: profile.BankAccountNumber,
		}
		result.Users = append(result.Users, userResponse)
	}
	return &result, nil
}

// GetUserDetails implements user.UserServiceServer.
func (puc ProtoUserController) GetUserDetailsWithId(ctx context.Context, request *user.UserRequest) (*user.UsersWithIdResponse, error) {
	fmt.Printf("Masuk request\n")
	result := user.UsersWithIdResponse{}

	// todo; check cache ada profile user tidak?
	// belum tau kalian mau pakai cache middleware atau internal jadi saya kasih catatan saja-ad1ee

	response, err := puc.userService.GetUserProfilesWithId(ctx, request.UserIds)
	if err != nil {
		fmt.Sprintf("Error %s", err.Error())
		puc.logger.Error(err.Error(), functionCallerInfo.UserControllerGetUserProfile)
		return nil, err
	}

	// Populate the response with the data from user profiles.
	for _, profile := range response {
		userResponse := &user.UserWithIdResponse{
			UserId:            profile.UserId,
			Email:             profile.Email,
			Phone:             profile.Phone,
			BankAccountName:   profile.BankAccountName,
			BankAccountHolder: profile.BankAccountHolder,
			BankAccountNumber: profile.BankAccountNumber,
		}
		result.Users = append(result.Users, userResponse)
	}
	fmt.Printf("Returning result: %s\n", result)
	return &result, nil
}

// mustEmbedUnimplementedUserServiceServer implements user.UserServiceServer.
func (puc ProtoUserController) mustEmbedUnimplementedUserServiceServer() {
	return
}
