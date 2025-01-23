package activityService

import (
	"github.com/TimDebug/FitByte/src/model/dtos/request"
	"github.com/TimDebug/FitByte/src/model/dtos/response"
	"github.com/gofiber/fiber/v2"
)

type ActivityServiceInterface interface {
	Create(ctx *fiber.Ctx, req request.RequestActivity) (response.ResponseActivity, error)
	Update(ctx *fiber.Ctx, req request.RequestActivity, userId, activityId string) (response.ResponseActivity, error)
	Delete(ctx *fiber.Ctx, userId, activityId string) error
}
