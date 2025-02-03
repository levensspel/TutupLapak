package route

import (
	"github.com/TIM-DEBUG-ProjectSprintBatch3/go-fiber-template/src/http/controller"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/go-fiber-template/src/http/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetRouteProduct(router fiber.Router, pc controller.ProductControllerInterface) {
	router.Post("/product", middleware.AuthMiddleware, pc.Create)
	router.Delete("/product/:productId", middleware.AuthMiddleware, pc.DeleteById)
	router.Put("/product/:productId", middleware.AuthMiddleware, pc.UpdateById)
	router.Get("/product", pc.GetAll)
}
