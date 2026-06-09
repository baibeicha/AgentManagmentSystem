package repository

import (
	"AgentManagmentSystem/internal/auth/domain"
	"AgentManagmentSystem/pkg/storage"
	"context"
)

type TenantRepository struct {
	db *storage.DB
}

func NewTenantRepository(db *storage.DB) *TenantRepository {
	return &TenantRepository{db: db}
}

func (r *TenantRepository) Create(ctx context.Context, tenant *domain.Tenant) error {
	query := `INSERT INTO tenants (name) VALUES ($1) RETURNING id, created_at, updated_at`
	return r.db.QueryRow(ctx, query, tenant.Name).Scan(&tenant.ID, &tenant.CreatedAt, &tenant.UpdatedAt)
}

func (r *TenantRepository) GetByID(ctx context.Context, id int64) (*domain.Tenant, error) {
	query := `SELECT id, name, created_at, updated_at FROM tenants WHERE id = $1`
	tenant := &domain.Tenant{}
	err := r.db.QueryRow(ctx, query, id).Scan(&tenant.ID, &tenant.Name, &tenant.CreatedAt, &tenant.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return tenant, nil
}

func (r *TenantRepository) Update(ctx context.Context, tenant *domain.Tenant) error {
	query := `UPDATE tenants SET name = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2 RETURNING updated_at`
	return r.db.QueryRow(ctx, query, tenant.Name, tenant.ID).Scan(&tenant.UpdatedAt)
}

func (r *TenantRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM tenants WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
