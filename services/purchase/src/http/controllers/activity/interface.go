package activityController

import "github.com/gofiber/fiber/v2"

type ActivityControllerInterface interface {
	Create(C *fiber.Ctx) error
	Update(C *fiber.Ctx) error
	Delete(C *fiber.Ctx) error
}
