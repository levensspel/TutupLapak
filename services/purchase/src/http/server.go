package httpServer

import (
	"fmt"

	"github.com/TimDebug/FitByte/src/config"
	"github.com/TimDebug/FitByte/src/di"
	activityController "github.com/TimDebug/FitByte/src/http/controllers/activity"
	userController "github.com/TimDebug/FitByte/src/http/controllers/user"
	"github.com/TimDebug/FitByte/src/http/routes"
	activityroutes "github.com/TimDebug/FitByte/src/http/routes/activity"
	userroutes "github.com/TimDebug/FitByte/src/http/routes/user"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/do/v2"
)

type HttpServer struct{}

func (s *HttpServer) Listen() {
	fmt.Printf("New Fiber\n")
	app := fiber.New(fiber.Config{
		ServerHeader: "TIM-DEBUG",
	})

	fmt.Printf("Inject Controllers\n")
	//? Dependency Injection
	//? UserController
	uc := do.MustInvoke[userController.UserControllerInterface](di.Injector)
	ac := do.MustInvoke[activityController.ActivityControllerInterface](di.Injector)

	routes := routes.SetRoutes(app)
	userroutes.SetRouteUsers(routes, uc)
	activityroutes.SetRouteActivities(routes, ac)

	fmt.Printf("Start Listener\n")
	app.Listen(fmt.Sprintf("%s:%s", "0.0.0.0", config.GetPort()))
}
