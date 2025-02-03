package service

import (
	"context"

	"github.com/TIM-DEBUG-ProjectSprintBatch3/go-fiber-template/src/model/dtos/request"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/go-fiber-template/src/model/dtos/response"
)

type ProductServiceInterface interface {
	Create(ctx context.Context, payload request.ProductCreate) (response.ProductCreate, error)
	DeletedById(ctx context.Context, productId string, userId string) error
	UpdateById(ctx context.Context, payload request.ProductUpdate) (response.ProductCreate, error)
	GetAll(ctx context.Context, filter request.ProductFilter) ([]response.ProductCreate, error)
}
