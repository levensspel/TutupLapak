package httpServer

import (
	"math/rand/v2"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type FileController struct {
	Service FileService
}

func NewFileController(service FileService) FileController {
	return FileController{
		Service: service,
	}
}

func (f *FileController) Upload(ctx *fiber.Ctx) error {
	if i := randRange(1, 10); i <= 5 {
		var code int
		switch i {
		case 1, 2:
			code = http.StatusUnauthorized
		case 3:
			code = http.StatusConflict
		case 4:
			code = http.StatusBadRequest
		case 5:
			code = http.StatusInternalServerError
		}
		return &fiber.Error{Code: code, Message: "Got an error"}
	}
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"code": http.StatusOK, "message": "OK"})
}

func randRange(min, max int) int {
	return rand.IntN(max-min) + min
}
