package httpServer

import (
	"fmt"
	"strconv"
	"time"

	"github.com/TimDebug/FitByte/src/config"
	"github.com/TimDebug/FitByte/src/di"
	appController "github.com/TimDebug/FitByte/src/http/controllers/purchase"
	"github.com/TimDebug/FitByte/src/http/routes"
	swaggerRoutes "github.com/TimDebug/FitByte/src/http/routes/apidocumentation"
	purchaseRoute "github.com/TimDebug/FitByte/src/http/routes/purchase"
	response "github.com/TimDebug/FitByte/src/model/web"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/samber/do/v2"
)

type HttpServer struct{}

func (s *HttpServer) Listen() {
	fmt.Printf("New Fiber\n")
	app := fiber.New(fiber.Config{
		ServerHeader: "TIM-DEBUG",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusBadRequest).JSON(response.GlobalErrorHandlerResp{
				Success: false,
				Message: err.Error(),
			})
		},
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
	})

	// Setup Middlewares
	app.Use(logger.New(logger.Config{
		Format:     "${time} ${status} - ${method} ${path} - Internal Latency: ${latency}\n",
		TimeFormat: "2006-01-02 15:04:05",
	}))

	// app.Use(middlewares.RequestLogger)
	app.Use(cache.New(cache.Config{
		ExpirationGenerator: func(c *fiber.Ctx, cfg *cache.Config) time.Duration {
			newCacheTime, _ := strconv.Atoi(c.GetRespHeader("Cache-Time", "600"))
			return time.Second * time.Duration(newCacheTime)
		},
		KeyGenerator: func(c *fiber.Ctx) string {
			return utils.CopyString(c.Path())
		},
		CacheControl: true,
		Expiration:   10 * time.Second,
	}))

	fmt.Printf("Inject Controllers\n")
	pc := do.MustInvoke[appController.IPurchaseController](di.Injector)
	fmt.Printf("Prepare Routes\n")
	routes := routes.SetRoutes(app)
	swaggerRoutes.SetRouteSwagger(routes)
	purchaseRoute.SetRoutePurchase(routes, pc)

	fmt.Printf("Start Listener\n")
	app.Listen(fmt.Sprintf("%s:%s", "0.0.0.0", config.GetPort()))
}
