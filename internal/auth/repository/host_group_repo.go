package repository

import (
	"AgentManagmentSystem/internal/auth/domain"
	"AgentManagmentSystem/pkg/storage"
	"context"
)

type HostGroupRepository struct {
	db *storage.DB
}

func NewHostGroupRepository(db *storage.DB) *HostGroupRepository {
	return &HostGroupRepository{db: db}
}

func (r *HostGroupRepository) Create(ctx context.Context, group *domain.HostGroup) error {
	query := `INSERT INTO host_groups (tenant_id, name) VALUES ($1, $2) RETURNING id, created_at`
	return r.db.QueryRow(ctx, query, group.TenantID, group.Name).Scan(&group.ID, &group.CreatedAt)
}

func (r *HostGroupRepository) GetByID(ctx context.Context, id int64) (*domain.HostGroup, error) {
	query := `SELECT id, tenant_id, name, created_at FROM host_groups WHERE id = $1`
	group := &domain.HostGroup{}
	err := r.db.QueryRow(ctx, query, id).Scan(&group.ID, &group.TenantID, &group.Name, &group.CreatedAt)
	if err != nil {
		return nil, err
	}
	return group, nil
}

func (r *HostGroupRepository) GetByTenantID(ctx context.Context, tenantID int64) ([]*domain.HostGroup, error) {
	query := `SELECT id, tenant_id, name, created_at FROM host_groups WHERE tenant_id = $1`
	rows, err := r.db.Query(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*domain.HostGroup
	for rows.Next() {
		var group domain.HostGroup
		if err := rows.Scan(&group.ID, &group.TenantID, &group.Name, &group.CreatedAt); err != nil {
			return nil, err
		}
		groups = append(groups, &group)
	}
	return groups, rows.Err()
}

func (r *HostGroupRepository) Update(ctx context.Context, group *domain.HostGroup) error {
	query := `UPDATE host_groups SET name = $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, query, group.Name, group.ID)
	return err
}

func (r *HostGroupRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM host_groups WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
