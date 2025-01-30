package di

import (
	authJwt "github.com/TimDebug/FitByte/src/auth/jwt"
	"github.com/TimDebug/FitByte/src/database/postgre"
	purchaseGrpc "github.com/TimDebug/FitByte/src/grpc"
	appController "github.com/TimDebug/FitByte/src/http/controllers/purchase"
	loggerZap "github.com/TimDebug/FitByte/src/logger/zap"
	purchaseRepository "github.com/TimDebug/FitByte/src/repositories/purchase"
	purchaseCartRepository "github.com/TimDebug/FitByte/src/repositories/purchaseCart"
	purchaseService "github.com/TimDebug/FitByte/src/services/purchase"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/do/v2"
)

var Injector *do.RootScope

func init() {
	Injector = do.New()

	//? Setup Database Connection
	do.Provide[*pgxpool.Pool](Injector, postgre.NewPgxConnectInject)

	//? Logger
	//? Zap
	do.Provide[loggerZap.LoggerInterface](Injector, loggerZap.NewLogHandlerInject)

	//? GRPCs
	//? UserService
	do.Provide[*purchaseGrpc.ProtoUserController](Injector, purchaseGrpc.NewGRPCClientInject)

	//? Setup Auth
	//? JWT Service
	do.Provide[authJwt.JwtServiceInterface](Injector, authJwt.NewJwtServiceInject)

	// Repositories
	do.Provide[purchaseRepository.IPurchaseRepository](Injector, purchaseRepository.NewPurhcaseRepositoryInject)
	do.Provide[purchaseCartRepository.IPuchaseCartRepository](Injector, purchaseCartRepository.NewPurhcaseCartRepositoryInject)
	// Services
	do.Provide[*purchaseService.PurchaseService](Injector, purchaseService.NewInject)
	// Controllers
	do.Provide[appController.IPurchaseController](Injector, appController.NewPurchaseControllerInject)
}
