package userController

import (
	"net/http"
	"net/url"
	"unicode/utf8"

	"github.com/TimDebug/FitByte/src/exceptions"
	functionCallerInfo "github.com/TimDebug/FitByte/src/logger/helper"
	loggerZap "github.com/TimDebug/FitByte/src/logger/zap"
	"github.com/TimDebug/FitByte/src/model/dtos/request"
	userService "github.com/TimDebug/FitByte/src/services/user"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/do/v2"
)

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

func (uc *UserController) Login(c *fiber.Ctx) error {
	bodyParsed := request.UserRegister{}
	if err := c.BodyParser(&bodyParsed); err != nil {
		uc.logger.Error(err.Error(), functionCallerInfo.UserControllerLogin)
		return c.Status(http.StatusBadRequest).JSON(exceptions.ErrBadRequest(err.Error()))
	}

	response, err := uc.userService.Login(c, bodyParsed)
	if err != nil {
		uc.logger.Error(err.Error(), functionCallerInfo.UserControllerLogin, bodyParsed)
		return c.Status(int(err.(exceptions.ErrorResponse).StatusCode)).
			JSON(err)
	}

	c.Set("X-Author", "TIM-DEBUG")
	return c.Status(http.StatusOK).JSON(response)
}

func (uc *UserController) Register(c *fiber.Ctx) error {
	userRequestParse := request.UserRegister{}

	if err := c.BodyParser(&userRequestParse); err != nil {
		uc.logger.Error(err.Error(), functionCallerInfo.UserControllerRegister)
		return c.Status(http.StatusBadRequest).JSON(exceptions.ErrBadRequest(err.Error()))
	}

	response, err := uc.userService.Register(c, userRequestParse)
	if err != nil {
		uc.logger.Error(err.Error(), functionCallerInfo.UserControllerRegister, userRequestParse)
		return c.Status(int(err.(exceptions.ErrorResponse).StatusCode)).
			JSON(err)
	}

	c.Set("X-Author", "TIM-DEBUG")
	return c.Status(201).JSON(response)
}

func (uc *UserController) Update(c *fiber.Ctx) error {
	userId, ok := c.Locals("userId").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(
			exceptions.NewUnauthorizedError(
				fiber.ErrUnauthorized.Error(),
				fiber.StatusUnauthorized,
			),
		)
	}

	updateRequest := request.UpdateProfile{}

	if err := c.BodyParser(&updateRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
		// panic(exceptions.NewBadRequestError(err.Error(), 400))
	}

	if updateRequest.Name == nil {
		return c.Status(fiber.StatusBadRequest).JSON("name is nil")
	}
	nameLength := utf8.RuneCountInString(*updateRequest.Name)
	if nameLength < 2 || nameLength > 62 {
		return fiber.NewError(fiber.StatusBadRequest, "name must be between 2 and 62 characters")
	}

	// Manually Validate URI
	// Validate ImageUri
	if updateRequest.ImageUri == nil || *updateRequest.ImageUri == "" {
		return c.Status(fiber.StatusBadRequest).JSON("imageuri is required and cannot be empty")
	}
	// Parse and validate the URI
	parsedUri, err := url.ParseRequestURI(*updateRequest.ImageUri)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("imageuri is not a valid URI")
	}
	// Optional: Check for schemes like http or https
	if parsedUri.Scheme != "http" && parsedUri.Scheme != "https" {
		return c.Status(fiber.StatusBadRequest).JSON("imageuri must have http or https scheme")
	}
	// Ensure there's a path (e.g., "/image.jpg")
	if parsedUri.Path == "" || parsedUri.Path == "/" {
		return c.Status(fiber.StatusBadRequest).JSON("must have an path and in image format")
	}
	//end Validate URI

	response, err := uc.userService.Update(c, userId, updateRequest)
	if err != nil {
		uc.logger.Error(err.Error(), functionCallerInfo.ProfileControllerUpdate, userId, updateRequest)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	c.Set("X-Author", "TIM-DEBUG")
	return c.Status(fiber.StatusOK).JSON(response)
}
