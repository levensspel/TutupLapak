package userService

import (
	"strings"
	"time"

	authJwt "github.com/TimDebug/FitByte/src/auth/jwt"
	"github.com/TimDebug/FitByte/src/exceptions"
	functionCallerInfo "github.com/TimDebug/FitByte/src/logger/helper"
	loggerZap "github.com/TimDebug/FitByte/src/logger/zap"
	"github.com/TimDebug/FitByte/src/model/dtos/request"
	"github.com/TimDebug/FitByte/src/model/dtos/response"
	Entity "github.com/TimDebug/FitByte/src/model/entities/user"
	userRepository "github.com/TimDebug/FitByte/src/repositories/user"
	"github.com/TimDebug/FitByte/src/services/user/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/do/v2"
	"golang.org/x/crypto/bcrypt"
)

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

func (us *userService) Login(ctx *fiber.Ctx, input request.UserRegister) (response.UserRegister, error) {
	err := validator.ValidateAuthParams(input)
	if err != nil {
		return response.UserRegister{}, exceptions.ErrBadRequest(err.Error())
	}

	user := Entity.User{}
	user.Email = input.Email
	users, err := us.UserRepository.Login(ctx, us.Db, &user)
	if err != nil {
		us.Logger.Error(err.Error(), functionCallerInfo.UserServiceLogin)
		return response.UserRegister{}, exceptions.ErrServer(err.Error())
	}
	if len(users) == 0 {
		return response.UserRegister{}, exceptions.ErrNotFound("User not found")
	}

	passwordHash := users[0].PasswordHash
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(input.Password))
	if err != nil {
		return response.UserRegister{}, exceptions.ErrBadRequest(err.Error())
	}

	token, err := us.jwtService.GenerateToken(users[0].Id)
	if err != nil {
		us.Logger.Error(err.Error(), functionCallerInfo.UserServiceLogin, err)
		return response.UserRegister{}, exceptions.ErrServer(err.Error())
	}
	return response.UserRegister{Email: users[0].Email, Token: token}, nil
}

func (us *userService) Register(ctx *fiber.Ctx, input request.UserRegister) (response.UserRegister, error) {
	err := validator.ValidateAuthParams(input)
	if err != nil {
		return response.UserRegister{}, exceptions.ErrBadRequest(err.Error())
	}

	user := Entity.User{}
	timeNow := time.Now()
	user.CreatedAt = timeNow
	user.UpdatedAt = timeNow
	user.Email = input.Email

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return response.UserRegister{}, exceptions.ErrBadRequest(err.Error())
	}

	user.PasswordHash = string(passwordHash)
	user.Id, err = us.UserRepository.CreateUser(ctx, us.Db, user)

	if err != nil {
		if strings.Contains(err.Error(), "23505") {
			return response.UserRegister{}, exceptions.ErrConflict(err.Error())
		} else {
			us.Logger.Error(err.Error(), functionCallerInfo.UserServiceRegister)
			return response.UserRegister{}, exceptions.ErrBadRequest(err.Error())
		}
	}

	token, err := us.jwtService.GenerateToken(user.Id)
	if err != nil {
		return response.UserRegister{}, exceptions.ErrServer(err.Error())
	}

	return response.UserRegister{
		Email: user.Email,
		Token: token,
	}, nil
}

func (us *userService) Update(ctx *fiber.Ctx, id string, req request.UpdateProfile) (*response.UpdateProfile, error) {
	profile, err := us.UserRepository.FindById(ctx, id)
	if err != nil {
		us.Logger.Error(err.Error(), functionCallerInfo.UserServiceUpdate, err)
		return nil, err
	}

	profile.Id = id
	profile.Preference = &req.Preference
	profile.WeightUnit = &req.WeightUnit
	profile.HeightUnit = &req.HeightUnit
	profile.Weight = &req.Weight
	profile.Height = &req.Height
	if req.Name != nil {
		profile.Name = req.Name
	}
	if req.ImageUri != nil {
		profile.ImageUri = req.ImageUri
	}

	_, err = us.UserRepository.Update(ctx, *profile)
	if err != nil {
		us.Logger.Error(err.Error(), functionCallerInfo.UserServiceUpdate, err)
		return nil, err
	}

	return &response.UpdateProfile{
		Preference: *profile.Preference,
		WeightUnit: *profile.WeightUnit,
		HeightUnit: *profile.HeightUnit,
		Weight:     *profile.Weight,
		Height:     *profile.Height,
		Name:       profile.Name,
		ImageUri:   profile.ImageUri,
	}, nil
}
