package userService

import (
	"context"

	authJwt "github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/auth/jwt"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/exceptions"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/helper"
	functionCallerInfo "github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/logger/helper"
	loggerZap "github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/logger/zap"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/model/dtos/request"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/model/dtos/response"
	userRepository "github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/repositories/user"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/services/user/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/do/v2"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceInterface interface {
	RegisterByEmail(ctx *fiber.Ctx, input request.AuthByEmailRequest) (response.AuthResponse, error)
	RegisterByPhone(ctx *fiber.Ctx, input request.AuthByPhoneRequest) (response.AuthResponse, error)
	LoginByEmail(ctx *fiber.Ctx, input request.AuthByEmailRequest) (response.AuthResponse, error)
	LoginByPhone(ctx *fiber.Ctx, input request.AuthByPhoneRequest) (response.AuthResponse, error)
}

type userService struct {
	UserRepository userRepository.UserRepositoryInterface
	Db             *pgxpool.Pool
	jwtService     authJwt.JwtServiceInterface
	Logger         loggerZap.LoggerInterface
}

func NewUserService(userRepo userRepository.UserRepositoryInterface, db *pgxpool.Pool, jwtService authJwt.JwtServiceInterface, logger loggerZap.LoggerInterface) UserServiceInterface {
	return &userService{UserRepository: userRepo, Db: db, jwtService: jwtService, Logger: logger}
}

func NewUserServiceInject(i do.Injector) (UserServiceInterface, error) {
	_db := do.MustInvoke[*pgxpool.Pool](i)
	_userRepo := do.MustInvoke[userRepository.UserRepositoryInterface](i)
	_jwtService := do.MustInvoke[authJwt.JwtServiceInterface](i)
	_logger := do.MustInvoke[loggerZap.LoggerInterface](i)

	return NewUserService(_userRepo, _db, _jwtService, _logger), nil
}

func (us *userService) RegisterByEmail(ctx *fiber.Ctx, input request.AuthByEmailRequest) (response.AuthResponse, error) {
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

	userId, err := us.UserRepository.CreateUserByEmail(context.Background(), us.Db, input.Email, passwordHash)
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

func (us *userService) RegisterByPhone(ctx *fiber.Ctx, input request.AuthByPhoneRequest) (response.AuthResponse, error) {
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

func (us *userService) LoginByEmail(ctx *fiber.Ctx, input request.AuthByEmailRequest) (response.AuthResponse, error) {
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

func (us *userService) LoginByPhone(ctx *fiber.Ctx, input request.AuthByPhoneRequest) (response.AuthResponse, error) {
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
