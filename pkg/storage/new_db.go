package storage

import "AgentManagmentSystem/pkg/config"

func NewPostgresDB(cfg *config.Config) (*DB, error) {
	return ConnectSource(config.NewDatasourceFromConfig(cfg))
}
