package httpServer

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/config"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/di"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/helper"
	swaggerRoutes "github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/http/controllers/apiDocumentation"
	userController "github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/http/controllers/user"
	protoUserController "github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/http/controllers/user/proto"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/http/routes"
	userroutes "github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/http/routes/user"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/model/dtos/response"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/services/proto/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/samber/do/v2"
	"google.golang.org/grpc"
)

type ServerInterface interface {
	Listen()
}

type HttpServer struct{}

func (s *HttpServer) Listen() {
	// Start gRPC server
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}

		grpcServer := grpc.NewServer()
		// pb.RegisterExampleServiceServer(grpcServer, &exampleServiceServer{})
		puc := do.MustInvoke[*protoUserController.ProtoUserController](di.Injector)
		user.RegisterUserServiceServer(grpcServer, puc)

		log.Println("gRPC server listening on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	config.FILE_SERVICE_BASE_URL = config.GetFileServiceBaseURL()
	if config.FILE_SERVICE_BASE_URL == "" {
		panic("FILE_SERVICE_BASE_URL value requires to be set")
	}
	config.MODE = config.GetMode()
	if config.MODE == "" {
		panic("MODE value requires to be set")
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

	app.Use(recover.New())

	if strings.ToUpper(config.MODE) == config.MODE_DEBUG {
		// terminal logger
		app.Use(logger.New(logger.Config{
			Format:     "${time} ${status} - ${method} ${path} - Internal Latency: ${latency}\n",
			TimeFormat: "2006-01-02 15:04:05",
		}))

		// resource monitoring
		app.Get("/monitor", monitor.New(monitor.Config{
			Title:   "User Service Metrics",
			Refresh: 5 * time.Second,
			APIOnly: false,
		}))
	}

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
