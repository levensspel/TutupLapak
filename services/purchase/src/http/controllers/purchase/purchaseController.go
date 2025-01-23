package appController

import (
	loggerZap "github.com/TimDebug/FitByte/src/logger/zap"
	"github.com/samber/do/v2"
)

type PurchaseController struct {
	logger loggerZap.LoggerInterface
}

func NewActivityController(logger loggerZap.LoggerInterface) IPurchaseController {
	return &PurchaseController{logger: logger}
}

func NewActivityControllerInject(i do.Injector) (IPurchaseController, error) {
	_logger := do.MustInvoke[loggerZap.LoggerInterface](i)
	return NewActivityController(_logger), nil
}
