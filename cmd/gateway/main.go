package main

import (
	"AgentManagmentSystem/internal/gateway/delivery/http/router"
	"AgentManagmentSystem/pkg/config"
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
	cfg := config.MustLoad("gateway-config.yaml")

	log := logger.SetupLogger(
		cfg.GetString("log.type"),
		cfg.GetString("log.level"),
	)

	r := router.SetupRouter(cfg, log)

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
