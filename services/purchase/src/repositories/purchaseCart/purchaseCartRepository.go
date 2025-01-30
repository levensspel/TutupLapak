package purchaseCartRepository

import (
	"context"

	functionCallerInfo "github.com/TimDebug/FitByte/src/logger/helper"
	loggerZap "github.com/TimDebug/FitByte/src/logger/zap"
	"github.com/TimDebug/FitByte/src/model/dtos/request"
	"github.com/jackc/pgx/v5"
	"github.com/samber/do/v2"
)

type PuchaseCartRepository struct {
	logger loggerZap.LoggerInterface
}

func NewPurchaseCartRepository(logger loggerZap.LoggerInterface) IPuchaseCartRepository {
	return &PuchaseCartRepository{logger}
}

func NewPurhcaseCartRepositoryInject(i do.Injector) (IPuchaseCartRepository, error) {
	_logger := do.MustInvoke[loggerZap.LoggerInterface](i)
	return NewPurchaseCartRepository(_logger), nil
}

func (pr *PuchaseCartRepository) InsertInto(tx pgx.Tx, ctx context.Context, purchaseId string, entities []request.PurchasedItem) error {
	// Insert purchased items ke tabel terkait
	query := `
	INSERT INTO purchase_cart (purchase_id, product_id, quantity) 
	VALUES ($1, $2, $3)
	`
	for _, item := range entities {
		_, err := tx.Exec(ctx, query, purchaseId, item.ProductId, item.Qty)
		if err != nil {
			pr.logger.Error("failed to insert purchased item", functionCallerInfo.PurchaseRepositoryInsertInto, err.Error())
			return err
		}
	}
	return nil
}
