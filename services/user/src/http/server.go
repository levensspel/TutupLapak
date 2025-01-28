package httpServer

import (
	"fmt"

	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/config"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/di"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/helper"
	swaggerRoutes "github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/http/controllers/apiDocumentation"
	userController "github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/http/controllers/user"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/http/routes"
	userroutes "github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/http/routes/user"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/model/dtos/response"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/samber/do/v2"
)

type ServerInterface interface {
	Listen()
}

type HttpServer struct{}

func (s *HttpServer) Listen() {
	config.FILE_SERVICE_BASE_URL = config.GetFileServiceBaseURL()
	if config.FILE_SERVICE_BASE_URL == "" {
		panic("FILE_SERVICE_BASE_URL value requires to be set")
	}

	fmt.Printf("New Fiber\n")
	app := fiber.New(fiber.Config{
		ServerHeader: helper.X_AUTHOR_HEADER_VALUE,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusBadRequest).JSON(response.GlobalErrorHandlerResp{
				Success: false,
				Message: err.Error(),
			})
		},
	})

	// terminal logger
	app.Use(logger.New(logger.Config{
		Format:     "${time} ${status} - ${method} ${path} - Internal Latency: ${latency}\n",
		TimeFormat: "2006-01-02 15:04:05",
	}))

	app.Use(recover.New())

	fmt.Printf("Inject Controllers\n")
	//? Dependency Injection
	//? UserController
	uc := do.MustInvoke[userController.UserControllerInterface](di.Injector)

	routes := routes.SetRoutes(app)
	swaggerRoutes.SetRouteSwagger(routes)
	userroutes.SetRouteUsers(routes, uc)

	fmt.Printf("Start Listener\n")
	app.Listen(fmt.Sprintf("%s:%s", "0.0.0.0", config.GetPort()))
}
