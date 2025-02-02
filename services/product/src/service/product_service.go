package service

import (
	"context"
	"time"

	"github.com/TIM-DEBUG-ProjectSprintBatch3/go-fiber-template/src/exceptions"
	loggerZap "github.com/TIM-DEBUG-ProjectSprintBatch3/go-fiber-template/src/logger/zap"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/go-fiber-template/src/model/dtos/request"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/go-fiber-template/src/model/dtos/response"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/go-fiber-template/src/model/entity"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/go-fiber-template/src/repository"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/do/v2"
)

type ProductService struct {
	DB          *pgxpool.Pool
	ProductRepo repository.ProductRepoInterface
	Logger      loggerZap.LoggerInterface
	Validation  *validator.Validate
}

func New(db *pgxpool.Pool, productRepo repository.ProductRepoInterface, logger loggerZap.LoggerInterface, validation *validator.Validate) ProductServiceInterface {
	return &ProductService{
		DB:          db,
		ProductRepo: productRepo,
		Logger:      logger,
		Validation:  validation,
	}
}

func NewInject(i do.Injector) (ProductServiceInterface, error) {
	_db := do.MustInvoke[*pgxpool.Pool](i)
	_productRepo := do.MustInvoke[repository.ProductRepoInterface](i)
	_logger := do.MustInvoke[loggerZap.LoggerInterface](i)
	_validation := do.MustInvoke[*validator.Validate](i)

	return New(_db, _productRepo, _logger, _validation), nil
}

func (ps *ProductService) Create(ctx context.Context, payload request.ProductCreate) (response.ProductCreate, error) {
	err := ps.Validation.Struct(payload)
	if err != nil {
		return response.ProductCreate{}, exceptions.NewBadRequestError(err.Error())
	}

	// TODO FindFileId

	time := time.Now()
	product := entity.Product{
		UserId:    payload.UserId,
		Name:      *payload.Name,
		Category:  *payload.Category,
		Qty:       *payload.Qty,
		Price:     *payload.Price,
		Sku:       *payload.Sku,
		FileId:    *payload.FileId,
		CreatedAt: time,
		UpdateAt:  time,
	}

	id, err := ps.ProductRepo.Create(ctx, ps.DB, product)

	if err != nil {
		return response.ProductCreate{}, err
	}

	return response.ProductCreate{
		ProductId:        id,
		Name:             product.Name,
		Qty:              product.Qty,
		Price:            product.Price,
		Sku:              product.Sku,
		FileId:           product.FileId,
		Category:         product.Category,
		FileUri:          "",
		FileThumbnailUri: "",
		CreatedAt:        product.CreatedAt,
		UpdatedAt:        product.UpdateAt,
	}, nil
}
