package httpServer

import (
	"fmt"
	_ "net/http/pprof" // Import pprof untuk otomatis registrasi ke http server
	"strings"
	"time"

	"github.com/TimDebug/FitByte/src/config"
	"github.com/TimDebug/FitByte/src/di"
	appController "github.com/TimDebug/FitByte/src/http/controllers/purchase"
	"github.com/TimDebug/FitByte/src/http/routes"
	swaggerRoutes "github.com/TimDebug/FitByte/src/http/routes/apidocumentation"
	purchaseRoute "github.com/TimDebug/FitByte/src/http/routes/purchase"
	response "github.com/TimDebug/FitByte/src/model/web"
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
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
	fmt.Printf("Setup middlewares\n")

	// Or extend your config for customization
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // 1
	}))
	// server profiling

	// Prometheus
	fmt.Printf("Setup prometheus\n")
	prometheus := fiberprometheus.New("purhcase-service")
	prometheus.RegisterAt(app, "/metrics")

	app.Use(prometheus.Middleware)

	// app.Use(middlewares.RequestLogger)
	if strings.ToUpper(config.GetMode()) == config.MODE_DEBUG {
		// app.Use(logger.New(logger.Config{
		// 	Format:     "${time} ${status} - ${method} ${path} - Internal Latency: ${latency}\n",
		// 	TimeFormat: "2006-01-02 15:04:05",
		// }))
		app.Use(func(c *fiber.Ctx) error {
			start := time.Now()
			err := c.Next()                                         // Lanjutkan ke handler berikutnya
			c.Set("X-Internal-Latency", time.Since(start).String()) // Simpan latency di response header
			return err
		})
		// Middleware PProf
		// go func() {
		// 	fmt.Println("Starting pprof")
		// 	app.Use(pprof.New())
		// }()
	}

	fmt.Printf("Inject Controllers\n")
	pc := do.MustInvoke[appController.IPurchaseController](di.Injector)
	fmt.Printf("Prepare Routes\n")
	routes := routes.SetRoutes(app)
	swaggerRoutes.SetRouteSwagger(routes)
	purchaseRoute.SetRoutePurchase(routes, pc)

	fmt.Printf("Start Listener\n")
	app.Listen(fmt.Sprintf("%s:%s", "0.0.0.0", config.GetPort()))
}
