package activityRepository

import (
	Entity "github.com/TimDebug/FitByte/src/model/entities/activity"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ActivityRepositoryInterface interface {
	Create(ctx *fiber.Ctx, pool *pgxpool.Pool, activity Entity.Activity) (activityId string, err error)
	GetValidCaloriesFactors(ctx *fiber.Ctx, pool *pgxpool.Pool, activityId, userId string) (*Entity.CaloriesFactor, error)
	GetActivityByUserId(ctx *fiber.Ctx, pool *pgxpool.Pool, activityId, userId string) (string, error)
	Update(ctx *fiber.Ctx, pool *pgxpool.Pool, activity Entity.Activity) error
	Delete(ctx *fiber.Ctx, pool *pgxpool.Pool, activityId, userId string) error
}
