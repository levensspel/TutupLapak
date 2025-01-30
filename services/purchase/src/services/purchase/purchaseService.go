package purchaseService

import (
	"context"
	"time"

	functionCallerInfo "github.com/TimDebug/FitByte/src/logger/helper"
	loggerZap "github.com/TimDebug/FitByte/src/logger/zap"
	"github.com/TimDebug/FitByte/src/model/dtos/request"
	Entity "github.com/TimDebug/FitByte/src/model/entities/purchase"
	purchaseRepository "github.com/TimDebug/FitByte/src/repositories/purchase"
	purchaseCartRepository "github.com/TimDebug/FitByte/src/repositories/purchaseCart"
	"github.com/gofiber/fiber/v2"
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
func (this PurchaseService) SaveCart(c *fiber.Ctx, entity request.CartDto) (*string, error) {
	// buka koneksi aja
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// buka transaksi db
	tx, err := this.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	// untuk handle panic() yang tidak tertangkap, return err pun enggak akan tertangkap oleh defer ini! Use it wisely!
	defer tx.Rollback(ctx)

	// Lanjut eksekusi query...
	senderDetail := Entity.Purchase{
		SenderName:          entity.SenderName,
		SenderContactDetail: entity.SenderContactDetail,
		SenderContactType:   entity.SenderContactType,
	}
	insertedId, err := this.purchaseRepository.InsertInto(tx, ctx, senderDetail)
	if err != nil {
		tx.Rollback(ctx)
		this.logger.Error(err.Error(), functionCallerInfo.PurchaserServiceSaveCart, senderDetail)
		return nil, err
	}
	err = this.purchaseCartRepository.InsertInto(tx, ctx, insertedId, entity.PurchasedItems)
	if err != nil {
		tx.Rollback(ctx)
		this.logger.Error(err.Error(), functionCallerInfo.PurchaserServiceSaveCart, entity.PurchasedItems)
		return nil, err
	}
	// commit transaksi
	if err := tx.Commit(ctx); err != nil {
		tx.Rollback(ctx)
		this.logger.Error(err.Error(), functionCallerInfo.PurchaserServiceSaveCart)
		return nil, err
	}

	// return euy
	return &insertedId, nil
}
