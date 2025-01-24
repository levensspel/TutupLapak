package appController

import (
	functionCallerInfo "github.com/TimDebug/FitByte/src/logger/helper"
	loggerZap "github.com/TimDebug/FitByte/src/logger/zap"
	"github.com/TimDebug/FitByte/src/model/dtos/request"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/do/v2"
)

type PurchaseController struct {
	logger loggerZap.LoggerInterface
}

func NewPurchaseController(logger loggerZap.LoggerInterface) IPurchaseController {
	return &PurchaseController{logger: logger}
}

func NewPurchaseControllerInject(i do.Injector) (IPurchaseController, error) {
	_logger := do.MustInvoke[loggerZap.LoggerInterface](i)
	return NewPurchaseController(_logger), nil
}

// Purchase godoc
// @Summary Add items to the cart
// @Description Pembeli akan memasukkan detail produk dan jumlah yang akan dibeli, kemudian mengembalikan daftar detail produk beserta dengan daftar detail bank dari masing-masing penjual
// @Tags Purchase
// @Accept json
// @Produce json
// @Param request body request.CartDto true "Cart Data"
// @Success 201 {object} response.PurchaseResponseDTO "success response"
// @Failure 400 {object} map[string]interface{} "bad request"
// @Failure 500 {object} map[string]interface{} "internal server error"
// @Router /v1/purchase [post]
func (pc *PurchaseController) Cart(c *fiber.Ctx) error {
	requestBody := new(request.CartDto)
	if err := c.BodyParser(requestBody); err != nil {
		pc.logger.Error(err.Error(), functionCallerInfo.PurhcaseControllerPutCart)
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return nil
}
