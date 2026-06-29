package main

import (
	"AgentManagmentSystem/internal/auth/api"
	"AgentManagmentSystem/pkg/grpc/interceptor"
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"AgentManagmentSystem/internal/auth/repository"
	"AgentManagmentSystem/internal/auth/service"
	"AgentManagmentSystem/pkg/config"
	"AgentManagmentSystem/pkg/jwt"
	"AgentManagmentSystem/pkg/logger"
	"AgentManagmentSystem/pkg/storage"

	authv1 "AgentManagmentSystem/pkg/api/grpc/auth/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.MustLoad("auth-config")

	_ = cfg.Datasource

	log, err := logger.New(
		logger.WithLevel(cfg.Log.Level),
		logger.WithEnv(cfg.Log.Type),
		logger.WithFileOutput(cfg.Log.File.ToArgs()),
	)

	if err != nil {
		log.Error("can not setup logger", "err", err)
		return
	}

	log.Info("Starting Auth Service...")

	db, err := storage.NewPostgresDB(cfg)
	if err != nil {
		log.Error("failed to connect to PostgreSQL", "err", err)
		return
	}
	defer db.Close()
	log.Info("Connected to PostgreSQL")

	redisClient, err := storage.NewRedisClient(context.Background(), cfg)
	if err != nil {
		log.Error("failed to connect to Redis", "err", err)
		return
	}
	defer redisClient.Close()
	log.Info("Connected to Redis")

	tokenRepo := repository.NewRedisTokenRepo(cfg, redisClient)

	jwtCfg := cfg.JWT
	tokenProvider, err := jwt.NewTokenProvider(cfg, tokenRepo, jwtCfg.TTL.GetAccessTTL(), jwtCfg.TTL.GetRefreshTTL())

	if err != nil {
		log.Error("failed to create token provider", "err", err)
		return
	}

	userRepo := repository.NewUserRepository(db)
	policyRepo := repository.NewResourcePolicyRepository(db)
	sessionRepo := repository.NewSessionRepository(db)

	authService := service.NewAuthService(
		cfg,
		tokenProvider,
		userRepo,
		policyRepo,
		sessionRepo,
	)

	authApi := api.NewAuthServiceApi(authService)

	grpcPort := cfg.GRPC.Port
	if grpcPort == "" {
		grpcPort = "50051"
	}

	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Error("failed to listen tcp", "err", err)
		return
	}

	var opts []grpc.ServerOption

	if cfg.GRPC.TLS.Enabled {
		creds, err := credentials.NewServerTLSFromFile(cfg.GRPC.TLS.CertPath, cfg.GRPC.TLS.KeyPath)
		if err != nil {
			log.Error("Failed to setup TLS", "err", err)
			return
		}
		opts = append(opts, grpc.Creds(creds))
		log.Info("gRPC server is starting in SECURE mode (TLS enabled)")
	} else {
		log.Info("gRPC server is starting in INSECURE mode (TLS disabled)")
	}

	opts = append(opts, grpc.ChainUnaryInterceptor(
		interceptor.LoggerInterceptor(),
	))

	grpcServer := grpc.NewServer(opts...)

	authv1.RegisterAuthServiceServer(grpcServer, authApi)

	reflection.Register(grpcServer)

	go func() {
		log.Info("Started gRPC server", "port", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Error("gRPC server crashed", "err", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Info("Gracefully shutting down Auth Service...")

	stopped := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(stopped)
	}()

	timer := time.NewTimer(5 * time.Second)
	select {
	case <-timer.C:
		log.Warn("Forcing gRPC server shutdown due to timeout")
		grpcServer.Stop()
	case <-stopped:
		timer.Stop()
	}

	log.Info("Auth Service successfully stopped")
}
