package repository

import (
	"context"

	"github.com/TIM-DEBUG-ProjectSprintBatch3/go-fiber-template/src/model/entity"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/do/v2"
)

type ProductRepository struct {
}

func New() ProductRepoInterface {
	return &ProductRepository{}
}

func NewInject(i do.Injector) (ProductRepoInterface, error) {
	return New(), nil
}

func (repository *ProductRepository) Create(ctx context.Context, pool *pgxpool.Pool, product entity.Product) (productId string, err error) {

	query := `INSERT INTO products (user_id, name, category, qty, price, sku, file_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`

	row := pool.QueryRow(ctx, query, product.UserId, product.Name, product.Category, product.Qty, product.Price, product.Sku, product.FileId, product.CreatedAt, product.UpdateAt)
	err = row.Scan(&productId)

	if err != nil {
		return "", err
	}

	return productId, nil
}
