package purchaseService

import (
	"context"
	"fmt"
	"time"

	functionCallerInfo "github.com/TimDebug/FitByte/src/logger/helper"
	loggerZap "github.com/TimDebug/FitByte/src/logger/zap"
	"github.com/TimDebug/FitByte/src/model/dtos/request"
	Entity "github.com/TimDebug/FitByte/src/model/entities/purchase"
	purchaseRepository "github.com/TimDebug/FitByte/src/repositories/purchase"
	purchaseCartRepository "github.com/TimDebug/FitByte/src/repositories/purchaseCart"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/do/v2"
)

type PurchaseService struct {
	logger                 loggerZap.LoggerInterface
	purchaseCartRepository purchaseCartRepository.IPuchaseCartRepository
	purchaseRepository     purchaseRepository.IPurchaseRepository
	db                     *pgxpool.Pool
}

func NewInject(i do.Injector) (*PurchaseService, error) {
	_logger := do.MustInvoke[loggerZap.LoggerInterface](i)
	_db := do.MustInvoke[*pgxpool.Pool](i)
	_pcr := do.MustInvoke[purchaseCartRepository.IPuchaseCartRepository](i)
	_pr := do.MustInvoke[purchaseRepository.IPurchaseRepository](i)
	return &PurchaseService{
		logger:                 _logger,
		db:                     _db,
		purchaseCartRepository: _pcr,
		purchaseRepository:     _pr,
	}, nil
}

// formely returned (*response.PurchaseResponseDTO, error)
// Service yang menggunakan pool
func (this PurchaseService) SaveCart(c *fiber.Ctx, entity request.CartDto) (*string, error) {
	requestId := uuid.New()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Dapatkan koneksi dari pool (dengan menunggu jika pool penuh)
	start := time.Now()

	conn, err := this.db.Acquire(ctx)
	if err != nil {
		this.logger.Error(err.Error(), functionCallerInfo.PurchaserServiceSaveCart, "Acquire Connection", fmt.Sprintf("RequestID:%s|WaitTime:%v", requestId, time.Since(start)))
		return nil, err
	}
	defer conn.Release()

	// Gunakan koneksi dari pool
	tx, err := conn.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		this.logger.Error(err.Error(), functionCallerInfo.PurchaserServiceSaveCart, "Begin Transaction", fmt.Sprintf("RequestID:%s", requestId))
		return nil, err
	}
	defer tx.Rollback(ctx)

	// Eksekusi query
	senderDetail := Entity.Purchase{
		SenderName:          entity.SenderName,
		SenderContactDetail: entity.SenderContactDetail,
		SenderContactType:   entity.SenderContactType,
	}
	insertedId, err := this.purchaseRepository.InsertInto(tx, ctx, senderDetail)
	if err != nil {
		this.logger.Error(err.Error(), functionCallerInfo.PurchaserServiceSaveCart, senderDetail, fmt.Sprintf("RequestID:%s", requestId))
		tx.Rollback(ctx)
		return nil, err
	}

	err = this.purchaseCartRepository.InsertInto(tx, ctx, insertedId, entity.PurchasedItems)
	if err != nil {
		this.logger.Error(err.Error(), functionCallerInfo.PurchaserServiceSaveCart, entity.PurchasedItems, fmt.Sprintf("RequestID:%s", requestId))
		tx.Rollback(ctx)
		return nil, err
	}

	// Commit transaksi jika sukses
	if err := tx.Commit(ctx); err != nil {
		this.logger.Error(err.Error(), functionCallerInfo.PurchaserServiceSaveCart, fmt.Sprintf("RequestID:%s", requestId))
		tx.Rollback(ctx)
		return nil, err
	}

	// Return inserted ID
	return &insertedId, nil
}
