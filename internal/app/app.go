package app

import (
	"fmt"
	"net"
	"net/http"
	"notification-service/internal/config"
	"notification-service/internal/delivery"
	"notification-service/internal/server"
	"notification-service/internal/service"
	"notification-service/pkg/tracing"
	"notification-service/pkg/dialog"
	"notification-service/pkg/logger"

	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	proto "github.com/aidostt/protos/gen/go/reservista/mailer"
)

func Run(configPath, envPath string) {
	cfg, err := config.Init(configPath, envPath)
	if err != nil {
		logger.Error(err)

		return
	}

	shutdownTracing, err := tracing.Init(context.Background(), "notification-service")
	if err != nil {
		logger.Errorf("tracing init: %s", err.Error())
	} else {
		defer func() { _ = shutdownTracing(context.Background()) }()
	}

	dial := dialog.NewDialog(cfg.Authority,
		fmt.Sprintf("%v:%v", cfg.QRs.Host, cfg.QRs.Port),
		fmt.Sprintf("%v:%v", cfg.Reservations.Host, cfg.Reservations.Port),
		fmt.Sprintf("%v:%v", cfg.Users.Host, cfg.Users.Port),
	)

	services := service.NewService(service.Dependencies{
		Host:     cfg.SMTP.Host,
		Port:     cfg.SMTP.Port,
		Username: cfg.SMTP.Username,
		Password: cfg.SMTP.Password,
		Sender:   cfg.SMTP.Sender,
		Dialog:   dial,
	})
	handlers := delivery.NewHandler(services)

	// gRPC Server
	srv := server.NewServer()
	proto.RegisterMailerServer(srv.GrpcServer, handlers)

	l, err := net.Listen("tcp", fmt.Sprintf("%v:%v", cfg.GRPC.Host, cfg.GRPC.Port))
	if err != nil {
		logger.Errorf("error occurred while getting listener for the server: %s\n", err.Error())
		return
	}
	go func() {
		if err := srv.Run(l); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("error occurred while running grpc server: %s\n", err.Error())
		}
	}()

	metricsSrv := &http.Server{Addr: ":" + metricsPort(), Handler: server.MetricsHandler()}
	go func() {
		if err := metricsSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("metrics server error: %s\n", err.Error())
		}
	}()

	logger.Info("Server started at: " + cfg.GRPC.Host + ":" + cfg.GRPC.Port)

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
	srv.Stop()
	_ = metricsSrv.Close()
	logger.Info("Stopping server at: " + cfg.GRPC.Host + ":" + cfg.GRPC.Port)

}

// metricsPort is the port for the Prometheus metrics endpoint; it defaults to
// 9464 and can be overridden with METRICS_PORT.
func metricsPort() string {
	if p := os.Getenv("METRICS_PORT"); p != "" {
		return p
	}
	return "9464"
}
