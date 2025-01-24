package appController

import "github.com/gofiber/fiber/v2"

type IPurchaseController interface {
	Cart(c *fiber.Ctx) error
	// Create(C *fiber.Ctx) error
	// Update(C *fiber.Ctx) error
	// Delete(C *fiber.Ctx) error
}
