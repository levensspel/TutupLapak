package httpServer

import (
	"fmt"

	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/config"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/di"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/helper"
	userController "github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/http/controllers/user"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/http/routes"
	userroutes "github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/http/routes/user"
	"github.com/gofiber/fiber/v2"
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
	})

	app.Use(recover.New())

	fmt.Printf("Inject Controllers\n")
	//? Dependency Injection
	//? UserController
	uc := do.MustInvoke[userController.UserControllerInterface](di.Injector)

	routes := routes.SetRoutes(app)
	userroutes.SetRouteUsers(routes, uc)

	fmt.Printf("Start Listener\n")
	app.Listen(fmt.Sprintf("%s:%s", "0.0.0.0", config.GetPort()))
}
