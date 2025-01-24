package activityRepository

import (
	loggerZap "github.com/TimDebug/FitByte/src/logger/zap"
	"github.com/samber/do/v2"
)

type PuchaseRepository struct {
	logger loggerZap.LoggerInterface
}

func NewActivityRepository(logger loggerZap.LoggerInterface) IPuchaseRepository {
	return &PuchaseRepository{logger}
}

func NewActivityRepositoryInject(i do.Injector) (IPuchaseRepository, error) {
	_logger := do.MustInvoke[loggerZap.LoggerInterface](i)
	return NewActivityRepository(_logger), nil
}
