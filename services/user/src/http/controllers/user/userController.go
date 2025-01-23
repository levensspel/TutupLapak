package userController

import (
	"net/http"

	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/exceptions"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/helper"
	functionCallerInfo "github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/logger/helper"
	loggerZap "github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/logger/zap"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/model/dtos/request"
	userService "github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/services/user"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/do/v2"
)

type UserControllerInterface interface {
	RegisterByEmail(C *fiber.Ctx) error
	RegisterByPhone(C *fiber.Ctx) error
	LoginByEmail(C *fiber.Ctx) error
	LoginByPhone(C *fiber.Ctx) error
}

type UserController struct {
	userService userService.UserServiceInterface
	logger      loggerZap.LoggerInterface
}

func NewUserController(userService userService.UserServiceInterface, logger loggerZap.LoggerInterface) UserControllerInterface {
	return &UserController{userService: userService, logger: logger}
}

func NewUserControllerInject(i do.Injector) (UserControllerInterface, error) {
	_userService := do.MustInvoke[userService.UserServiceInterface](i)
	_logger := do.MustInvoke[loggerZap.LoggerInterface](i)
	return NewUserController(_userService, _logger), nil
}

func (uc *UserController) RegisterByEmail(c *fiber.Ctx) error {
	userRequestParse := request.AuthByEmailRequest{}

	if err := c.BodyParser(&userRequestParse); err != nil {
		uc.logger.Error(err.Error(), functionCallerInfo.UserControllerRegisterByEmail)
		return c.Status(http.StatusBadRequest).JSON(exceptions.ErrBadRequest(err.Error()))
	}

	response, err := uc.userService.RegisterByEmail(c, userRequestParse)
	if err != nil {
		uc.logger.Error(err.Error(), functionCallerInfo.UserControllerRegisterByEmail, userRequestParse)
		return c.Status(int(err.(exceptions.ErrorResponse).StatusCode)).
			JSON(err)
	}

	c.Set(helper.X_AUTHOR_HEADER_KEY, helper.X_AUTHOR_HEADER_VALUE)
	return c.Status(fiber.StatusCreated).JSON(response)
}

func (uc *UserController) RegisterByPhone(c *fiber.Ctx) error {
	userRequestParse := request.AuthByPhoneRequest{}

	if err := c.BodyParser(&userRequestParse); err != nil {
		uc.logger.Error(err.Error(), functionCallerInfo.UserControllerRegisterByPhone)
		return c.Status(http.StatusBadRequest).JSON(exceptions.ErrBadRequest(err.Error()))
	}

	response, err := uc.userService.RegisterByPhone(c, userRequestParse)
	if err != nil {
		uc.logger.Error(err.Error(), functionCallerInfo.UserControllerRegisterByPhone, userRequestParse)
		return c.Status(int(err.(exceptions.ErrorResponse).StatusCode)).
			JSON(err)
	}

	c.Set(helper.X_AUTHOR_HEADER_KEY, helper.X_AUTHOR_HEADER_VALUE)
	return c.Status(fiber.StatusCreated).JSON(response)
}

func (uc *UserController) LoginByEmail(c *fiber.Ctx) error {
	userRequestParse := request.AuthByEmailRequest{}

	if err := c.BodyParser(&userRequestParse); err != nil {
		uc.logger.Error(err.Error(), functionCallerInfo.UserControllerLoginByEmail)
		return c.Status(http.StatusBadRequest).JSON(exceptions.ErrBadRequest(err.Error()))
	}

	response, err := uc.userService.LoginByEmail(c, userRequestParse)
	if err != nil {
		uc.logger.Error(err.Error(), functionCallerInfo.UserControllerLoginByEmail, userRequestParse)
		return c.Status(int(err.(exceptions.ErrorResponse).StatusCode)).JSON(err)
	}

	c.Set(helper.X_AUTHOR_HEADER_KEY, helper.X_AUTHOR_HEADER_VALUE)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (uc *UserController) LoginByPhone(c *fiber.Ctx) error {
	userRequestParse := request.AuthByPhoneRequest{}

	if err := c.BodyParser(&userRequestParse); err != nil {
		uc.logger.Error(err.Error(), functionCallerInfo.UserControllerRegisterByPhone)
		return c.Status(http.StatusBadRequest).JSON(exceptions.ErrBadRequest(err.Error()))
	}

	response, err := uc.userService.LoginByPhone(c, userRequestParse)
	if err != nil {
		uc.logger.Error(err.Error(), functionCallerInfo.UserControllerLoginByPhone, userRequestParse)
		return c.Status(int(err.(exceptions.ErrorResponse).StatusCode)).JSON(err)
	}

	c.Set(helper.X_AUTHOR_HEADER_KEY, helper.X_AUTHOR_HEADER_VALUE)
	return c.Status(fiber.StatusOK).JSON(response)
}
