package appController

import (
	"fmt"
	"regexp"
	"strings"

	helper "github.com/TimDebug/FitByte/src/helper/validator"
	functionCallerInfo "github.com/TimDebug/FitByte/src/logger/helper"
	loggerZap "github.com/TimDebug/FitByte/src/logger/zap"
	"github.com/TimDebug/FitByte/src/model/dtos/request"
	"github.com/TimDebug/FitByte/src/model/dtos/response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/do/v2"
)

type PurchaseController struct {
	logger    loggerZap.LoggerInterface
	validator helper.XValidator
}

func NewPurchaseController(logger loggerZap.LoggerInterface) IPurchaseController {
	xValidator := helper.XValidator{Validator: validator.New()}
	xValidator.Validator.RegisterValidation("sender_email_or_phone", func(fl validator.FieldLevel) bool {
		contactType := fl.Parent().FieldByName("SenderContactType").String()
		value := fl.Field().String()

		switch contactType {
		case "email":
			emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
			return regexp.MustCompile(emailRegex).MatchString(value)
		case "phone":
			phoneRegex := `^\+?[1-9]\d{1,14}$` // E.164 format
			return regexp.MustCompile(phoneRegex).MatchString(value)
		default:
			return false
		}
	})
	return &PurchaseController{logger: logger, validator: xValidator}
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
		return err
	}

	// Validation
	if errs := pc.validator.Validate(requestBody); len(errs) > 0 && errs[0].Error {
		errMsgs := make([]string, 0)

		for _, err := range errs {
			errMsgs = append(errMsgs, fmt.Sprintf(
				"[%s]: '%v' | Needs to implement '%s'",
				err.FailedField,
				err.Value,
				err.Tag,
			))
		}
		return fiber.NewError(fiber.StatusBadRequest, strings.Join(errMsgs, " and "))
	}
	return c.Status(fiber.StatusOK).JSON(response.PurchaseResponseDTO{})
}
