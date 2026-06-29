package client

import (
	"AgentManagmentSystem/pkg/config"
	"fmt"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type GRPCClientConfig struct {
	Addr                string
	Port                string
	PubliclyTrustedCert bool
	TLSCertPath         string
}

func GetGrpcConfig(cfg *config.Config, serviceName string) *GRPCClientConfig {
	return &GRPCClientConfig{
		Addr:                cfg.GetString("service." + serviceName + ".addr"),
		Port:                cfg.GetString("service." + serviceName + ".port"),
		PubliclyTrustedCert: cfg.GetBool("service." + serviceName + ".trusted_cert"),
		TLSCertPath:         cfg.GetString("service." + serviceName + ".cert_path"),
	}
}

func NewGrpcClient(cfg *config.Config, clientName string) (*grpc.ClientConn, error) {
	log := slog.Default()
	client := GetGrpcConfig(cfg, clientName)

	if client == nil {
		return nil, fmt.Errorf("client config is empty")
	}

	var creds credentials.TransportCredentials
	var err error
	if client.PubliclyTrustedCert {
		creds = credentials.NewClientTLSFromCert(nil, "")
	} else {
		creds, err = credentials.NewClientTLSFromFile(client.TLSCertPath, "")
	}

	if err != nil {
		log.Error(fmt.Sprintf("Can not load TLS cert: %v", err))
	}

	conn, err := grpc.NewClient(client.Addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Error(fmt.Sprintf("Can not connect to the server: %v", err))
	}

	return conn, nil
}
