package app

import (
	"fmt"
	"net"
	"net/http"
	"notification-service/internal/config"
	"notification-service/internal/delivery"
	"notification-service/internal/server"
	"notification-service/internal/service"
	"notification-service/pkg/dialog"
	"notification-service/pkg/logger"

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

	logger.Info("Server started at: " + cfg.GRPC.Host + ":" + cfg.GRPC.Port)

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
	srv.Stop()
	logger.Info("Stopping server at: " + cfg.GRPC.Host + ":" + cfg.GRPC.Port)

}
