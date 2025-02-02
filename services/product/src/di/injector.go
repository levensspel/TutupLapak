package di

import (
	authJwt "github.com/TIM-DEBUG-ProjectSprintBatch3/go-fiber-template/src/auth/jwt"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/go-fiber-template/src/database/postgre"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/go-fiber-template/src/http/controller"
	loggerZap "github.com/TIM-DEBUG-ProjectSprintBatch3/go-fiber-template/src/logger/zap"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/go-fiber-template/src/repository"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/go-fiber-template/src/service"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/go-fiber-template/src/validation"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/do/v2"
)

var Injector *do.RootScope

func init() {
	Injector = do.New()

	//? Setup Database Connection
	do.Provide[*pgxpool.Pool](Injector, postgre.NewPgxConnectInject)

	//? Setup Validation
	//? Validator
	do.Provide[*validator.Validate](Injector, validation.NewValidatorInject)

	//? Logger
	//? Zap
	do.Provide[loggerZap.LoggerInterface](Injector, loggerZap.NewInject)

	//? Setup Auth
	//? JWT Service
	do.Provide[authJwt.JwtServiceInterface](Injector, authJwt.NewJwtServiceInject)

	//? Setup Repositories
	//? Product Repository
	do.Provide[repository.ProductRepoInterface](Injector, repository.NewInject)

	//? Setup Services
	//? Product Service
	do.Provide[service.ProductServiceInterface](Injector, service.NewInject)

	//? Setup Controller/Handler
	//? Product Controller
	do.Provide[controller.ProductControllerInterface](Injector, controller.NewInject)
}
