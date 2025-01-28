package userGrpc

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/config"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/di"
	protoUserController "github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/http/controllers/user/proto"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/services/proto/user"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/samber/do/v2"
	"google.golang.org/grpc"
)

type UserGrpcServer struct {
}

func StartGrpcServer() {
	// Start gRPC server
	// Buat registry untuk Prometheus
	registry := prometheus.NewRegistry()
	grpcMetrics := grpc_prometheus.NewServerMetrics()

	// Registrasikan grpcMetrics ke registry
	registry.MustRegister(grpcMetrics)

	_PORT := config.GetGRPCPort()
	_GRPC_METRIC_PORT := config.GetMetricGRPCPort()

	// Buat gRPC server dengan interceptor Prometheus
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpcMetrics.UnaryServerInterceptor()),
		grpc.StreamInterceptor(grpcMetrics.StreamServerInterceptor()),
	)

	// Registrasikan service gRPC Anda
	puc := do.MustInvoke[*protoUserController.ProtoUserController](di.Injector)
	user.RegisterUserServiceServer(grpcServer, puc)

	// Inisialisasi metrik Prometheus untuk gRPC
	grpcMetrics.InitializeMetrics(grpcServer)

	// Jalankan HTTP server untuk endpoint metrik
	httpServer := &http.Server{
		Handler: promhttp.HandlerFor(registry, promhttp.HandlerOpts{}),
		Addr:    fmt.Sprintf("0.0.0.0:%s", _GRPC_METRIC_PORT), // Port untuk endpoint Prometheus
	}
	go func() {
		fmt.Printf("> Starting Prometheus metrics server at 0.0.0.0:%s\n", _GRPC_METRIC_PORT)
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Failed to start Prometheus metrics server: %v", err)
		}
	}()
	go func() {
		// Jalankan gRPC server
		lis, err := net.Listen("tcp", fmt.Sprintf(":%s", _PORT))
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}

		fmt.Printf("> gRPC server listening on :%s\n", _PORT)

		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()
}
