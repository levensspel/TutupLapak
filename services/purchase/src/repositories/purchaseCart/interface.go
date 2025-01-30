package purchaseCartRepository

import (
	"context"

	"github.com/TimDebug/FitByte/src/model/dtos/request"
	"github.com/jackc/pgx/v5"
)

type IPuchaseCartRepository interface {
	InsertInto(tx pgx.Tx, ctx context.Context, purchaseId string, entities []request.PurchasedItem) error
	// Create(ctx *fiber.Ctx, pool *pgxpool.Pool, activity Entity.Activity) (activityId string, err error)
	// GetValidCaloriesFactors(ctx *fiber.Ctx, pool *pgxpool.Pool, activityId, userId string) (*Entity.CaloriesFactor, error)
	// GetActivityByUserId(ctx *fiber.Ctx, pool *pgxpool.Pool, activityId, userId string) (string, error)
	// Update(ctx *fiber.Ctx, pool *pgxpool.Pool, activity Entity.Activity) error
	// Delete(ctx *fiber.Ctx, pool *pgxpool.Pool, activityId, userId string) error
}
