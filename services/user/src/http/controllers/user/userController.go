package userController

import (
	"net/http"

	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/exceptions"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/helper"
	functionCallerInfo "github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/logger/helper"
	loggerZap "github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/logger/zap"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/model/dtos/request"
	fileService "github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/services/external/file"
	userService "github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/services/user"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/do/v2"
)

type UserControllerInterface interface {
	RegisterByEmail(C *fiber.Ctx) error
	RegisterByPhone(C *fiber.Ctx) error
	LoginByEmail(C *fiber.Ctx) error
	LoginByPhone(C *fiber.Ctx) error

	LinkEmail(C *fiber.Ctx) error
	LinkPhone(C *fiber.Ctx) error
	GetUserProfile(C *fiber.Ctx) error
	UpdateUserProfile(C *fiber.Ctx) error
}

type UserController struct {
	userService userService.UserServiceInterface
	fileService fileService.FileServiceInterface
	logger      loggerZap.LoggerInterface
}

func NewUserController(userService userService.UserServiceInterface, fileService fileService.FileServiceInterface, logger loggerZap.LoggerInterface) UserControllerInterface {
	return &UserController{userService: userService, logger: logger}
}

func NewUserControllerInject(i do.Injector) (UserControllerInterface, error) {
	_userService := do.MustInvoke[userService.UserServiceInterface](i)
	_fileService := do.MustInvoke[fileService.FileServiceInterface](i)
	_logger := do.MustInvoke[loggerZap.LoggerInterface](i)
	return NewUserController(_userService, _fileService, _logger), nil
}

func (uc *UserController) RegisterByEmail(ctx *fiber.Ctx) error {
	userRequestParse := request.AuthByEmailRequest{}

	if err := ctx.BodyParser(&userRequestParse); err != nil {
		uc.logger.Error(err.Error(), functionCallerInfo.UserControllerRegisterByEmail)
		return ctx.Status(http.StatusBadRequest).JSON(exceptions.ErrBadRequest(err.Error()))
	}

	response, err := uc.userService.RegisterByEmail(ctx.Context(), userRequestParse)
	if err != nil {
		uc.logger.Error(err.Error(), functionCallerInfo.UserControllerRegisterByEmail, userRequestParse)
		return ctx.Status(int(err.(exceptions.ErrorResponse).StatusCode)).
			JSON(err)
	}

	ctx.Set(helper.X_AUTHOR_HEADER_KEY, helper.X_AUTHOR_HEADER_VALUE)
	return ctx.Status(fiber.StatusCreated).JSON(response)
}

func (uc *UserController) RegisterByPhone(ctx *fiber.Ctx) error {
	userRequestParse := request.AuthByPhoneRequest{}

	if err := ctx.BodyParser(&userRequestParse); err != nil {
		uc.logger.Error(err.Error(), functionCallerInfo.UserControllerRegisterByPhone)
		return ctx.Status(http.StatusBadRequest).JSON(exceptions.ErrBadRequest(err.Error()))
	}

	response, err := uc.userService.RegisterByPhone(ctx.Context(), userRequestParse)
	if err != nil {
		uc.logger.Error(err.Error(), functionCallerInfo.UserControllerRegisterByPhone, userRequestParse)
		return ctx.Status(int(err.(exceptions.ErrorResponse).StatusCode)).
			JSON(err)
	}

	ctx.Set(helper.X_AUTHOR_HEADER_KEY, helper.X_AUTHOR_HEADER_VALUE)
	return ctx.Status(fiber.StatusCreated).JSON(response)
}

func (uc *UserController) LoginByEmail(ctx *fiber.Ctx) error {
	userRequestParse := request.AuthByEmailRequest{}

	if err := ctx.BodyParser(&userRequestParse); err != nil {
		uc.logger.Error(err.Error(), functionCallerInfo.UserControllerLoginByEmail)
		return ctx.Status(http.StatusBadRequest).JSON(exceptions.ErrBadRequest(err.Error()))
	}

	response, err := uc.userService.LoginByEmail(ctx.Context(), userRequestParse)
	if err != nil {
		uc.logger.Error(err.Error(), functionCallerInfo.UserControllerLoginByEmail, userRequestParse)
		return ctx.Status(int(err.(exceptions.ErrorResponse).StatusCode)).JSON(err)
	}

	ctx.Set(helper.X_AUTHOR_HEADER_KEY, helper.X_AUTHOR_HEADER_VALUE)
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (uc *UserController) LoginByPhone(c *fiber.Ctx) error {
	userRequestParse := request.AuthByPhoneRequest{}

	if err := c.BodyParser(&userRequestParse); err != nil {
		uc.logger.Error(err.Error(), functionCallerInfo.UserControllerRegisterByPhone)
		return c.Status(http.StatusBadRequest).JSON(exceptions.ErrBadRequest(err.Error()))
	}

	response, err := uc.userService.LoginByPhone(c.Context(), userRequestParse)
	if err != nil {
		uc.logger.Error(err.Error(), functionCallerInfo.UserControllerLoginByPhone, userRequestParse)
		return c.Status(int(err.(exceptions.ErrorResponse).StatusCode)).JSON(err)
	}

	c.Set(helper.X_AUTHOR_HEADER_KEY, helper.X_AUTHOR_HEADER_VALUE)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (uc *UserController) LinkEmail(ctx *fiber.Ctx) error {
	userId, ok := ctx.Locals("userId").(string)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(exceptions.ErrUnauthorized)
	}

	userRequestParse := request.LinkEmailRequest{}

	if err := ctx.BodyParser(&userRequestParse); err != nil {
		uc.logger.Error(err.Error(), functionCallerInfo.UserControllerLinkEmail)
		return ctx.Status(http.StatusBadRequest).JSON(exceptions.ErrBadRequest(err.Error()))
	}

	response, err := uc.userService.LinkEmail(ctx.Context(), userRequestParse, userId)
	if err != nil {
		uc.logger.Error(err.Error(), functionCallerInfo.UserControllerLinkEmail, userRequestParse)
		return ctx.Status(int(err.(exceptions.ErrorResponse).StatusCode)).JSON(err)
	}

	ctx.Set(helper.X_AUTHOR_HEADER_KEY, helper.X_AUTHOR_HEADER_VALUE)
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (uc *UserController) LinkPhone(ctx *fiber.Ctx) error {
	userId, ok := ctx.Locals("userId").(string)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(exceptions.ErrUnauthorized)
	}

	userRequestParse := request.LinkPhoneRequest{}

	if err := ctx.BodyParser(&userRequestParse); err != nil {
		uc.logger.Error(err.Error(), functionCallerInfo.UserControllerLinkPhone)
		return ctx.Status(http.StatusBadRequest).JSON(exceptions.ErrBadRequest(err.Error()))
	}

	response, err := uc.userService.LinkPhone(ctx.Context(), userRequestParse, userId)
	if err != nil {
		uc.logger.Error(err.Error(), functionCallerInfo.UserControllerLinkPhone, userRequestParse)
		return ctx.Status(int(err.(exceptions.ErrorResponse).StatusCode)).JSON(err)
	}

	ctx.Set(helper.X_AUTHOR_HEADER_KEY, helper.X_AUTHOR_HEADER_VALUE)
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (uc *UserController) GetUserProfile(ctx *fiber.Ctx) error {
	userId, ok := ctx.Locals("userId").(string)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(exceptions.ErrUnauthorized)
	}

	response, err := uc.userService.GetUserProfile(ctx.Context(), userId)
	if err != nil {
		uc.logger.Error(err.Error(), functionCallerInfo.UserControllerGetUserProfile)
		return ctx.Status(int(err.(exceptions.ErrorResponse).StatusCode)).JSON(err)
	}

	ctx.Set(helper.X_AUTHOR_HEADER_KEY, helper.X_AUTHOR_HEADER_VALUE)
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (uc *UserController) UpdateUserProfile(ctx *fiber.Ctx) error {
	userId, ok := ctx.Locals("userId").(string)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(exceptions.ErrUnauthorized)
	}

	userRequestParse := request.UpdateUserProfileRequest{}

	if err := ctx.BodyParser(&userRequestParse); err != nil {
		uc.logger.Error(err.Error(), functionCallerInfo.UserControllerUpdateUserProfile)
		return ctx.Status(http.StatusBadRequest).JSON(exceptions.ErrBadRequest(err.Error()))
	}

	response, err := uc.userService.UpdateUserProfile(ctx.Context(), userRequestParse, userId)
	if err != nil {
		uc.logger.Error(err.Error(), functionCallerInfo.UserControllerUpdateUserProfile)
		return ctx.Status(int(err.(exceptions.ErrorResponse).StatusCode)).JSON(err)
	}

	ctx.Set(helper.X_AUTHOR_HEADER_KEY, helper.X_AUTHOR_HEADER_VALUE)
	return ctx.Status(fiber.StatusOK).JSON(response)
}
