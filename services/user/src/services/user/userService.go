package userService

import (
	"context"
	"fmt"

	authJwt "github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/auth/jwt"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/exceptions"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/helper"
	functionCallerInfo "github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/logger/helper"
	loggerZap "github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/logger/zap"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/model/dtos/repository"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/model/dtos/request"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/model/dtos/response"
	userRepository "github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/repositories/user"
	fileService "github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/services/external/file"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/services/user/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/do/v2"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceInterface interface {
	RegisterByEmail(ctx context.Context, input request.AuthByEmailRequest) (response.AuthResponse, error)
	RegisterByPhone(ctx context.Context, input request.AuthByPhoneRequest) (response.AuthResponse, error)
	LoginByEmail(ctx context.Context, input request.AuthByEmailRequest) (response.AuthResponse, error)
	LoginByPhone(ctx context.Context, input request.AuthByPhoneRequest) (response.AuthResponse, error)

	LinkEmail(ctx context.Context, input request.LinkEmailRequest, userId string) (response.UserResponse, error)
	LinkPhone(ctx context.Context, input request.LinkPhoneRequest, userId string) (response.UserResponse, error)
	GetUserProfile(ctx context.Context, userId string) (response.UserResponse, error)
	UpdateUserProfile(ctx context.Context, input request.UpdateUserProfileRequest, userId string) (response.UserResponse, error) //for grpc

	GetUserProfiles(ctx context.Context, userId []string) ([]response.UserResponse, error)
	GetUserProfilesWithId(ctx context.Context, userId []string) ([]response.UserWithIdResponse, error)
}

type userService struct {
	UserRepository userRepository.UserRepositoryInterface
	Db             *pgxpool.Pool
	jwtService     authJwt.JwtServiceInterface
	fileService    fileService.FileServiceInterface
	Logger         loggerZap.LoggerInterface
}

func NewUserService(
	userRepo userRepository.UserRepositoryInterface,
	db *pgxpool.Pool,
	jwtService authJwt.JwtServiceInterface,
	fileService fileService.FileServiceInterface,
	logger loggerZap.LoggerInterface,
) UserServiceInterface {
	return &userService{
		UserRepository: userRepo,
		Db:             db,
		jwtService:     jwtService,
		fileService:    fileService,
		Logger:         logger,
	}
}

func NewUserServiceInject(i do.Injector) (UserServiceInterface, error) {
	_db := do.MustInvoke[*pgxpool.Pool](i)
	_userRepo := do.MustInvoke[userRepository.UserRepositoryInterface](i)
	_jwtService := do.MustInvoke[authJwt.JwtServiceInterface](i)
	_fileService := do.MustInvoke[fileService.FileServiceInterface](i)
	_logger := do.MustInvoke[loggerZap.LoggerInterface](i)

	return NewUserService(_userRepo, _db, _jwtService, _fileService, _logger), nil
}

func (us *userService) RegisterByEmail(ctx context.Context, input request.AuthByEmailRequest) (response.AuthResponse, error) {
	err := validator.ValidateStructFields(input)
	if err != nil {
		return response.AuthResponse{}, exceptions.ErrBadRequest(err.Error())
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		us.Logger.Error(err.Error(), functionCallerInfo.UserServiceRegisterByEmail)
		return response.AuthResponse{}, exceptions.ErrBadRequest(err.Error())
	}

	passwordHash := string(hash)

	userId, err := us.UserRepository.CreateUserByEmail(ctx, us.Db, input.Email, passwordHash)
	if err != nil {
		us.Logger.Error(err.Error(), functionCallerInfo.UserRepositoryCreateUserByEmail, err)

		statusCode, message := helper.MapPgxError(err)
		return response.AuthResponse{}, exceptions.NewErrorResponse(statusCode, message)
	}

	token, err := us.jwtService.GenerateToken(userId)
	if err != nil {
		us.Logger.Error(err.Error(), functionCallerInfo.UserServiceRegisterByEmail)
		return response.AuthResponse{}, exceptions.ErrServer(err.Error())
	}

	return response.AuthResponse{
		Email: input.Email,
		Phone: "",
		Token: token,
	}, nil
}

func (us *userService) RegisterByPhone(ctx context.Context, input request.AuthByPhoneRequest) (response.AuthResponse, error) {
	err := validator.ValidateStructFields(input)
	if err != nil {
		return response.AuthResponse{}, exceptions.ErrBadRequest(err.Error())
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		us.Logger.Error(err.Error(), functionCallerInfo.UserServiceRegisterByEmail)
		return response.AuthResponse{}, exceptions.ErrBadRequest(err.Error())
	}

	passwordHash := string(hash)

	userId, err := us.UserRepository.CreateUserByPhone(context.Background(), us.Db, input.Phone, passwordHash)
	if err != nil {
		us.Logger.Error(err.Error(), functionCallerInfo.UserRepositoryCreateUserByPhone, err)

		statusCode, message := helper.MapPgxError(err)
		return response.AuthResponse{}, exceptions.NewErrorResponse(statusCode, message)
	}

	token, err := us.jwtService.GenerateToken(userId)
	if err != nil {
		us.Logger.Error(err.Error(), functionCallerInfo.UserServiceRegisterByPhone)
		return response.AuthResponse{}, exceptions.ErrServer(err.Error())
	}

	return response.AuthResponse{
		Email: "",
		Phone: input.Phone,
		Token: token,
	}, nil
}

func (us *userService) LoginByEmail(ctx context.Context, input request.AuthByEmailRequest) (response.AuthResponse, error) {
	err := validator.ValidateStructFields(input)
	if err != nil {
		return response.AuthResponse{}, exceptions.ErrBadRequest(err.Error())
	}

	auth, err := us.UserRepository.GetAuthByEmail(context.Background(), us.Db, input.Email)
	if err != nil {
		us.Logger.Error(err.Error(), functionCallerInfo.UserRepositoryGetAuthByEmail, err)

		statusCode, message := helper.MapPgxError(err)
		return response.AuthResponse{}, exceptions.NewErrorResponse(statusCode, message)
	}

	if auth.UserId == "" || auth.HashPassword == "" {
		return response.AuthResponse{}, exceptions.ErrNotFound("User not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(auth.HashPassword), []byte(input.Password))
	if err != nil {
		us.Logger.Error(err.Error(), functionCallerInfo.UserServiceLoginByEmail)
		return response.AuthResponse{}, exceptions.ErrBadRequest("Wrong password")
	}

	token, err := us.jwtService.GenerateToken(auth.UserId)
	if err != nil {
		us.Logger.Error(err.Error(), functionCallerInfo.UserServiceLoginByEmail)
		return response.AuthResponse{}, exceptions.ErrServer(err.Error())
	}

	return response.AuthResponse{
		Email: input.Email,
		Phone: auth.Phone,
		Token: token,
	}, nil
}

func (us *userService) LoginByPhone(ctx context.Context, input request.AuthByPhoneRequest) (response.AuthResponse, error) {
	err := validator.ValidateStructFields(input)
	if err != nil {
		return response.AuthResponse{}, exceptions.ErrBadRequest(err.Error())
	}

	auth, err := us.UserRepository.GetAuthByPhone(context.Background(), us.Db, input.Phone)
	if err != nil {
		us.Logger.Error(err.Error(), functionCallerInfo.UserRepositoryGetAuthByPhone, err)

		statusCode, message := helper.MapPgxError(err)
		return response.AuthResponse{}, exceptions.NewErrorResponse(statusCode, message)
	}

	if auth.UserId == "" || auth.HashPassword == "" {
		return response.AuthResponse{}, exceptions.ErrNotFound("User not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(auth.HashPassword), []byte(input.Password))
	if err != nil {
		us.Logger.Error(err.Error(), functionCallerInfo.UserServiceLoginByPhone)
		return response.AuthResponse{}, exceptions.ErrBadRequest("Wrong password")
	}

	token, err := us.jwtService.GenerateToken(auth.UserId)
	if err != nil {
		us.Logger.Error(err.Error(), functionCallerInfo.UserServiceLoginByPhone)
		return response.AuthResponse{}, exceptions.ErrServer(err.Error())
	}

	return response.AuthResponse{
		Email: auth.Email,
		Phone: input.Phone,
		Token: token,
	}, nil
}

func (us *userService) LinkEmail(ctx context.Context, input request.LinkEmailRequest, userId string) (response.UserResponse, error) {
	err := validator.ValidateStructFields(input)
	if err != nil {
		return response.UserResponse{}, exceptions.ErrBadRequest(err.Error())
	}

	user, err := us.UserRepository.UpdateEmail(context.Background(), us.Db, input.Email, userId)
	if err != nil {
		us.Logger.Error(err.Error(), functionCallerInfo.UserRepositoryUpdateEmail, err)

		statusCode, message := helper.MapPgxError(err)
		return response.UserResponse{}, exceptions.NewErrorResponse(statusCode, message)
	}

	if user == nil {
		us.Logger.Error("User not found", functionCallerInfo.UserServiceLinkEmail)
		return response.UserResponse{}, exceptions.NewNotFoundError("User not found", fiber.StatusNotFound)
	}

	response := helper.ConvertUserToResponse(user)
	response.Email = input.Email

	return response, nil
}

func (us *userService) LinkPhone(ctx context.Context, input request.LinkPhoneRequest, userId string) (response.UserResponse, error) {
	err := validator.ValidateStructFields(input)
	if err != nil {
		return response.UserResponse{}, exceptions.ErrBadRequest(err.Error())
	}

	user, err := us.UserRepository.UpdatePhone(context.Background(), us.Db, input.Phone, userId)
	if err != nil {
		us.Logger.Error(err.Error(), functionCallerInfo.UserRepositoryUpdatePhone, err)

		statusCode, message := helper.MapPgxError(err)
		return response.UserResponse{}, exceptions.NewErrorResponse(statusCode, message)
	}

	if user == nil {
		us.Logger.Error("User not found", functionCallerInfo.UserServiceLinkPhone)
		return response.UserResponse{}, exceptions.NewNotFoundError("User not found", fiber.StatusNotFound)
	}

	response := helper.ConvertUserToResponse(user)
	response.Phone = input.Phone

	return response, nil
}

func (us *userService) GetUserProfile(ctx context.Context, userId string) (response.UserResponse, error) {
	user, err := us.UserRepository.GetUserProfile(context.Background(), us.Db, userId)
	if err != nil {
		us.Logger.Error(err.Error(), functionCallerInfo.UserRepositoryGetUserProfile, err)

		statusCode, message := helper.MapPgxError(err)
		return response.UserResponse{}, exceptions.NewErrorResponse(statusCode, message)
	}

	if user == nil {
		us.Logger.Error("Bad request", functionCallerInfo.UserServiceGetUserProfile)
		return response.UserResponse{}, exceptions.NewBadRequestError("Bad request", fiber.StatusBadRequest)
	}

	response := helper.ConvertUserToResponse(user)

	return response, nil
}

func (us *userService) UpdateUserProfile(ctx context.Context, input request.UpdateUserProfileRequest, userId string) (response.UserResponse, error) {
	err := validator.ValidateStructFields(input)
	if err != nil {
		return response.UserResponse{}, exceptions.ErrBadRequest(err.Error())
	}

	updateUser := repository.UpdateUser{
		BankAccountName:   &input.BankAccountName,
		BankAccountHolder: &input.BankAccountHolder,
		BankAccountNumber: &input.BankAccountNumber,
	}

	if input.FileId != "" {
		file, statusCode := us.fileService.GetFile(ctx, input.FileId)
		fmt.Println("status", statusCode)
		if statusCode != fiber.StatusOK {
			if statusCode != fiber.StatusBadRequest {
				statusCode = fiber.StatusInternalServerError
			}
			// *Logging is done in the fileService implementation
			return response.UserResponse{}, exceptions.NewErrorResponse(int16(statusCode), "File error")
		}

		updateUser.FileId = &input.FileId
		updateUser.FileUri = &file.FileUri
		updateUser.FileThumbnailUri = &file.FileThumbnailUri
	}

	user, err := us.UserRepository.UpdateUserProfile(context.Background(), us.Db, updateUser, userId)
	if err != nil {
		us.Logger.Error(err.Error(), functionCallerInfo.UserRepositoryUpdateUserProfile, err)

		statusCode, message := helper.MapPgxError(err)
		return response.UserResponse{}, exceptions.NewErrorResponse(statusCode, message)
	}

	if user == nil {
		us.Logger.Error("Bad request", functionCallerInfo.UserServiceUpdateUserProfile)
		return response.UserResponse{}, exceptions.NewBadRequestError("Bad request", fiber.StatusBadRequest)
	}

	response := helper.ConvertUserToResponse(user)

	return response, nil
}

func (us *userService) GetUserProfiles(ctx context.Context, userId []string) ([]response.UserResponse, error) {
	users, err := us.UserRepository.GetUserProfiles(ctx, us.Db, userId)

	if err != nil {
		statusCode, message := helper.MapPgxError(err)
		us.Logger.Error(err.Error(), functionCallerInfo.UserRepositoryGetUserProfile, statusCode, message, err)
		// kalau ada masalah apa-pun, kembalikan badrequest
		return nil, exceptions.NewBadRequestError(fiber.ErrBadRequest.Message, fiber.StatusBadRequest)
	}

	if users == nil {
		us.Logger.Error(fiber.ErrBadRequest.Message, functionCallerInfo.UserServiceGetUserProfile)
		return nil, exceptions.NewBadRequestError(fiber.ErrBadRequest.Message, fiber.StatusBadRequest)
	}

	return users, nil
}

func (us *userService) GetUserProfilesWithId(ctx context.Context, userId []string) ([]response.UserWithIdResponse, error) {
	users, err := us.UserRepository.GetUserProfilesWithId(ctx, us.Db, userId)

	if err != nil {
		statusCode, message := helper.MapPgxError(err)
		us.Logger.Error(err.Error(), functionCallerInfo.UserRepositoryGetUserProfile, statusCode, message, err)
		// kalau ada masalah apa-pun, kembalikan badrequest
		return nil, exceptions.NewBadRequestError(fiber.ErrBadRequest.Message, fiber.StatusBadRequest)
	}

	if users == nil {
		us.Logger.Error(fiber.ErrBadRequest.Message, functionCallerInfo.UserServiceGetUserProfile)
		return nil, exceptions.NewBadRequestError(fiber.ErrBadRequest.Message, fiber.StatusBadRequest)
	}

	return users, nil
}
