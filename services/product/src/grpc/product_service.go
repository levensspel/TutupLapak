package grpc

import (
	"context"

	"github.com/TIM-DEBUG-ProjectSprintBatch3/go-fiber-template/src/model/dtos/response"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/go-fiber-template/src/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductService struct {
	DB          *pgxpool.Pool
	ProductRepo repository.ProductRepoInterface
}

func (ps *ProductService) GetDetailById(ctx context.Context, productId string) (response.Product, error) {

}
