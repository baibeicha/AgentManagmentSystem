package main

import (
	clients "AgentManagmentSystem/internal/gateway/client"
	"AgentManagmentSystem/internal/gateway/client/grpc"
	"AgentManagmentSystem/internal/gateway/client/mock"
	"AgentManagmentSystem/internal/gateway/delivery/http/router"
	"AgentManagmentSystem/internal/gateway/delivery/http/router/handler"
	"AgentManagmentSystem/pkg/config"
	"AgentManagmentSystem/pkg/grpc/client"
	"AgentManagmentSystem/pkg/logger"
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.MustLoad("gateway-config")

	log, err := logger.New(
		logger.WithLevel(cfg.Log.Level),
		logger.WithEnv(cfg.Log.Type),
		logger.WithFileOutput(cfg.Log.File.ToArgs()),
	)

	if err != nil {
		log.Error("can not setup logger", "err", err)
		return
	}

	authGrpcClient, err := client.NewGrpcClient(cfg, clients.AUTH)
	if err != nil {
		log.Error("can not setup auth grpc client", "err", err)
		return
	}
	defer authGrpcClient.Close()

	deviceMock := mock.NewDeviceMock()
	metricsMock := mock.NewMetricsMock()
	commandMock := mock.NewCommandMock()
	automationMock := mock.NewAutomationMock()
	notificationMock := mock.NewNotificationMock()
	incidentMock := mock.NewIncidentMock()
	teamMock := mock.NewTeamMock()
	auditMock := mock.NewAuditMock()

	h := handler.NewGatewayHandlers(
		grpc.NewAuthServiceClient(authGrpcClient),
		deviceMock,
		deviceMock,
		metricsMock,
		commandMock,
		automationMock,
		notificationMock,
		incidentMock,
		teamMock,
		auditMock,
	)

	r := router.SetupRouter(h)

	serverPort := cfg.GetString("server.port")
	if serverPort == "" {
		serverPort = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + serverPort,
		Handler: r,
	}

	go func() {
		log.Info("Started server on port " + serverPort)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("Error starting server", "err", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Info("Gracefully shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Stopping server because of error", "err", err)
	}

	log.Info("Server successfully stopped")
}
