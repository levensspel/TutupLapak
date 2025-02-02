package controller

import "github.com/gofiber/fiber/v2"

type ProductControllerInterface interface {
	Create(c *fiber.Ctx) error
}
