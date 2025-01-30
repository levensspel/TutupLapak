package purchaseRepository

import (
	"context"

	"github.com/TimDebug/FitByte/src/helper"
	functionCallerInfo "github.com/TimDebug/FitByte/src/logger/helper"
	loggerZap "github.com/TimDebug/FitByte/src/logger/zap"
	Entity "github.com/TimDebug/FitByte/src/model/entities/purchase"
	"github.com/jackc/pgx/v5"
	"github.com/samber/do/v2"
)

type PurchaseRepository struct {
	logger loggerZap.LoggerInterface
}

func NewPurchaseRepository(logger loggerZap.LoggerInterface) IPurchaseRepository {
	return &PurchaseRepository{logger}
}

func NewPurhcaseRepositoryInject(i do.Injector) (IPurchaseRepository, error) {
	_logger := do.MustInvoke[loggerZap.LoggerInterface](i)
	return NewPurchaseRepository(_logger), nil
}

func (pr *PurchaseRepository) InsertInto(tx pgx.Tx, ctx context.Context, entity Entity.Purchase) (string, error) {
	var id string
	query := `
	INSERT INTO 
	purchase(sender_name, sender_contact_detail, sender_contact_type) 
	VALUES($1, $2, $3) 
	RETURNING id
	`
	err := tx.QueryRow(ctx, query, entity.SenderName, entity.SenderContactDetail, entity.SenderContactType).Scan(&id)
	if err != nil {
		statusCode, message := helper.MapPgxError(err)
		pr.logger.Error(message, functionCallerInfo.PurchaseRepositoryInsertInto, message, statusCode)
		return "", err
	}
	return id, nil
}
