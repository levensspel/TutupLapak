package httpServer

import (
	"io"
	"net/http"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
	if ctx.Request().Header.ContentLength() > 100*1024 {
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "File size must be under 100 KiB",
		}
	}
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Unable to read uploaded file",
		}
	}
	file, err := fileHeader.Open()
	if err != nil {
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Unable to open uploaded file",
		}
	}
	defer file.Close()
	buff := make([]byte, 512)
	n, err := file.Read(buff)
	if err != nil && err != io.EOF {
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Unable to read file buffer",
		}
	}
	contentType := http.DetectContentType(buff[:n])
	allowedTypes := []string{"image/jpeg", "image/png"}
	if !contains(allowedTypes, contentType) {
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "File format not supported",
		}
	}
	if seeker, ok := file.(io.Seeker); ok {
		if _, err := seeker.Seek(0, 0); err != nil {
			return &fiber.Error{
				Code:    fiber.StatusInternalServerError,
				Message: "Cannot reset file pointer",
			}
		}
	} else {
		return &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "File pointer is not seekable",
		}
	}
	uniqueFileName := uuid.New().String() + filepath.Ext(fileHeader.Filename)
	mainEntity, err := f.Service.UploadFile(ctx, uniqueFileName, file, contentType)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return ctx.Status(fiber.StatusOK).
		JSON(mainEntity)
}

func contains(slice []string, item string) bool {
	for _, value := range slice {
		if value == item {
			return true
		}
	}
	return false
}
