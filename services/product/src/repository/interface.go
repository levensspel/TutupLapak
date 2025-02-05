package repository

import (
	"context"
	"time"

	"github.com/TIM-DEBUG-ProjectSprintBatch3/go-fiber-template/src/model/dtos/request"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/go-fiber-template/src/model/entity"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepoInterface interface {
	Create(ctx context.Context, pool *pgxpool.Pool, product entity.Product) (productId string, err error)
	DeleteById(ctx context.Context, pool *pgxpool.Pool, productId string, userId string) error
	UpdateById(ctx context.Context, pool *pgxpool.Pool, product entity.Product) (time.Time, error)
	GetAll(ctx context.Context, pool *pgxpool.Pool, filter request.ProductFilter) ([]entity.Product, error)
}
