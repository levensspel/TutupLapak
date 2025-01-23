package userController

import "github.com/gofiber/fiber/v2"

type UserControllerInterface interface {
	Login(C *fiber.Ctx) error
	Register(C *fiber.Ctx) error
	Update(C *fiber.Ctx) error
}
