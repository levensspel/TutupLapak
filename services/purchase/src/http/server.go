package httpServer

import (
	"fmt"

	"github.com/TimDebug/FitByte/src/config"
	"github.com/TimDebug/FitByte/src/di"
	appController "github.com/TimDebug/FitByte/src/http/controllers/purchase"
	"github.com/TimDebug/FitByte/src/http/middlewares"
	"github.com/TimDebug/FitByte/src/http/routes"
	swaggerRoutes "github.com/TimDebug/FitByte/src/http/routes/apidocumentation"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/do/v2"
)

type HttpServer struct{}

func (s *HttpServer) Listen() {
	fmt.Printf("New Fiber\n")
	app := fiber.New(fiber.Config{
		ServerHeader: "TIM-DEBUG",
	})

	app.Use(middlewares.LoggerMiddleware())

	fmt.Printf("Inject Controllers\n")
	do.MustInvoke[appController.IPurchaseController](di.Injector)

	routes := routes.SetRoutes(app)
	swaggerRoutes.SetRouteSwagger(routes)

	fmt.Printf("Start Listener\n")
	app.Listen(fmt.Sprintf("%s:%s", "0.0.0.0", config.GetPort()))
}
