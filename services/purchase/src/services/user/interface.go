package userService

import (
	"github.com/TimDebug/FitByte/src/model/dtos/request"
	"github.com/TimDebug/FitByte/src/model/dtos/response"
	"github.com/gofiber/fiber/v2"
)

type UserServiceInterface interface {
	Login(ctx *fiber.Ctx, input request.UserRegister) (response.UserRegister, error)
	Register(ctx *fiber.Ctx, input request.UserRegister) (response.UserRegister, error)
	Update(ctx *fiber.Ctx, id string, input request.UpdateProfile) (*response.UpdateProfile, error)
}
