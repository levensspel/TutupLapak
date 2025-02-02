package service

import (
	"context"

	"github.com/TIM-DEBUG-ProjectSprintBatch3/go-fiber-template/src/model/dtos/request"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/go-fiber-template/src/model/dtos/response"
)

type ProductServiceInterface interface {
	Create(ctx context.Context, payload request.ProductCreate) (response.ProductCreate, error)
}
