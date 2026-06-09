package repository

import (
	"AgentManagmentSystem/internal/auth/domain"
	"AgentManagmentSystem/pkg/storage"
	"context"
)

type HostRepository struct {
	db *storage.DB
}

func NewHostRepository(db *storage.DB) *HostRepository {
	return &HostRepository{db: db}
}

func (r *HostRepository) Create(ctx context.Context, host *domain.Host) error {
	query := `INSERT INTO hosts (tenant_id, group_id, hostname) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
	return r.db.QueryRow(ctx, query, host.TenantID, host.GroupID, host.Hostname).Scan(&host.ID, &host.CreatedAt, &host.UpdatedAt)
}

func (r *HostRepository) GetByID(ctx context.Context, id int64) (*domain.Host, error) {
	query := `SELECT id, tenant_id, group_id, hostname, created_at, updated_at FROM hosts WHERE id = $1`
	host := &domain.Host{}
	err := r.db.QueryRow(ctx, query, id).Scan(&host.ID, &host.TenantID, &host.GroupID, &host.Hostname, &host.CreatedAt, &host.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return host, nil
}

func (r *HostRepository) GetByTenantID(ctx context.Context, tenantID int64) ([]*domain.Host, error) {
	query := `SELECT id, tenant_id, group_id, hostname, created_at, updated_at FROM hosts WHERE tenant_id = $1`
	rows, err := r.db.Query(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hosts []*domain.Host
	for rows.Next() {
		var host domain.Host
		if err := rows.Scan(&host.ID, &host.TenantID, &host.GroupID, &host.Hostname, &host.CreatedAt, &host.UpdatedAt); err != nil {
			return nil, err
		}
		hosts = append(hosts, &host)
	}
	return hosts, rows.Err()
}

func (r *HostRepository) GetByGroupID(ctx context.Context, groupID int64) ([]*domain.Host, error) {
	query := `SELECT id, tenant_id, group_id, hostname, created_at, updated_at FROM hosts WHERE group_id = $1`
	rows, err := r.db.Query(ctx, query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hosts []*domain.Host
	for rows.Next() {
		var host domain.Host
		if err := rows.Scan(&host.ID, &host.TenantID, &host.GroupID, &host.Hostname, &host.CreatedAt, &host.UpdatedAt); err != nil {
			return nil, err
		}
		hosts = append(hosts, &host)
	}
	return hosts, rows.Err()
}

func (r *HostRepository) Update(ctx context.Context, host *domain.Host) error {
	query := `UPDATE hosts SET group_id = $1, hostname = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3 RETURNING updated_at`
	return r.db.QueryRow(ctx, query, host.GroupID, host.Hostname, host.ID).Scan(&host.UpdatedAt)
}

func (r *HostRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM hosts WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
