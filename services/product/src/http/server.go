package httpServer

import (
	"fmt"

	"github.com/TIM-DEBUG-ProjectSprintBatch3/go-fiber-template/src/config"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/go-fiber-template/src/di"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/go-fiber-template/src/http/controller"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/go-fiber-template/src/http/middleware"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/go-fiber-template/src/http/route"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/samber/do/v2"
)

type HttpServer struct{}

func (s *HttpServer) Listen() {
	fmt.Printf("New Fiber\n")
	app := fiber.New(fiber.Config{
		ServerHeader: "TIM-DEBUG",
		ErrorHandler: middleware.ErrorHandle,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowCredentials: false,
		AllowMethods:     "GET,POST,PUT,DELETE",
	}))

	fmt.Printf("Inject Controllers\n")
	//? Depedency Injection
	//? ProductController
	pc := do.MustInvoke[controller.ProductControllerInterface](di.Injector)

	routes := route.SetRoutes(app)
	route.SetRouteProduct(routes, pc)

	fmt.Printf("Start Lister\n")
	app.Listen(fmt.Sprintf("%s:%s", "0.0.0.0", config.GetPort()))
}
