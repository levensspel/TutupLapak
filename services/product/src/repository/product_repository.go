package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/TIM-DEBUG-ProjectSprintBatch3/go-fiber-template/src/exceptions"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/go-fiber-template/src/model/dtos/request"
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

func (pr *ProductRepository) Create(ctx context.Context, pool *pgxpool.Pool, product entity.Product) (productId string, err error) {

	query := `INSERT INTO products (user_id, name, category, qty, price, sku, file_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`

	row := pool.QueryRow(ctx, query, product.UserId, product.Name, product.Category, product.Qty, product.Price, product.Sku, product.FileId, product.CreatedAt, product.UpdatedAt)
	err = row.Scan(&productId)

	if err != nil {
		return "", err
	}

	return productId, nil
}

func (pr *ProductRepository) DeleteById(ctx context.Context, pool *pgxpool.Pool, productId string, userId string) error {
	query := `DELETE FROM products WHERE id = $1 AND user_id = $2 RETURNING id`

	row := pool.QueryRow(ctx, query, productId, userId)

	var deletedId string
	err := row.Scan(&deletedId)
	if err != nil {
		return exceptions.NewNotFoundError(productId + "is not found")
	}

	return nil
}

func (pr *ProductRepository) UpdateById(ctx context.Context, pool *pgxpool.Pool, product entity.Product) (time.Time, error) {

	query := `UPDATE products SET name = $1, category = $2 , qty = $3, price = $4, sku = $5, file_id = $6 , updated_at = $7 WHERE id = $8 AND user_id = $9 RETURNING id, created_at`

	row := pool.QueryRow(ctx, query, product.Name, product.Category, product.Qty, product.Price, product.Sku, product.FileId, product.UpdatedAt, product.Id, product.UserId)

	var (
		productId string
		createdAt time.Time
	)
	err := row.Scan(&productId, &createdAt)
	if err != nil {
		return product.UpdatedAt, exceptions.NewNotFoundError(product.Id + "is not found")
	}

	return createdAt, nil
}

func (pr *ProductRepository) GetAll(ctx context.Context, pool *pgxpool.Pool, filter request.ProductFilter) ([]entity.Product, error) {
	var args []interface{}
	argCounter := 1

	// Query dasar
	query := `SELECT id, name, category, qty, price, sku, file_id, created_at, updated_at FROM products WHERE 1=1`

	if filter.ProductId != "" {
		query += fmt.Sprintf(" AND id = $%d", argCounter)
		args = append(args, filter.ProductId)
		argCounter++
	}

	if filter.Sku != "" {
		query += fmt.Sprintf(" AND sku = $%d", argCounter)
		args = append(args, filter.Sku)
		argCounter++
	}

	if filter.Category != "" {
		query += fmt.Sprintf(" AND category = $%d", argCounter)
		args = append(args, filter.Category)
		argCounter++
	}

	// Sorting berdasarkan SortBy
	if filter.SortBy != "" {
		if filter.SortBy == "newest" {
			query += fmt.Sprintf(" ORDER BY updated_at DESC, created_at DESC")
		}
		if filter.SortBy == "cheapest" {
			query += fmt.Sprintf(" ORDER BY price ASC")
		}
		if strings.HasPrefix(filter.SortBy, "sold-") {
			// Ambil angka setelah sold-
			parts := strings.Split(filter.SortBy, "-")
			fmt.Print(parts)
		}
	}

	// Offset
	if filter.Offset >= 0 {
		query += fmt.Sprintf(" OFFSET $%d", argCounter)
		args = append(args, filter.Offset)
		argCounter++
	}

	// Limit
	if filter.Limit >= 0 {
		query += fmt.Sprintf(" LIMIT $%d", argCounter)
		args = append(args, filter.Limit)
		argCounter++
	}

	// Debugging query
	fmt.Println(query)

	// Eksekusi query
	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []entity.Product
	// Iterasi hasil query
	for rows.Next() {
		var product entity.Product
		if err := rows.Scan(&product.Id, &product.Name, &product.Category, &product.Qty, &product.Price, &product.Sku, &product.FileId, &product.CreatedAt, &product.UpdatedAt); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	// Cek apakah ada error dalam iterasi rows
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}
