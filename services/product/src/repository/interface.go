package repository

import (
	"context"

	"github.com/TIM-DEBUG-ProjectSprintBatch3/go-fiber-template/src/model/entity"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepoInterface interface {
	Create(ctx context.Context, pool *pgxpool.Pool, product entity.Product) (productId string, err error)
}
