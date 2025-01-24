package purchaseRoute

import (
	appController "github.com/TimDebug/FitByte/src/http/controllers/purchase"
	"github.com/gofiber/fiber/v2"
)

func SetRoutePurchase(router fiber.Router, controller appController.IPurchaseController) {
	router.Post("/purchase", controller.Cart)
	// router.Post("/purchase/:purchaseId", controller.Payment)
}
