package repository

import (
	"AgentManagmentSystem/internal/auth/domain"
	"AgentManagmentSystem/pkg/storage"
	"context"
)

type UserRepository struct {
	db *storage.DB
}

func NewUserRepository(db *storage.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (tenant_id, login, password_hash, global_role) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	return r.db.QueryRow(ctx, query, user.TenantID, user.Login, user.PasswordHash, user.GlobalRole).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	query := `SELECT id, tenant_id, login, password_hash, global_role, created_at, updated_at FROM users WHERE id = $1`
	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, id).Scan(&user.ID, &user.TenantID, &user.Login, &user.PasswordHash, &user.GlobalRole, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByLogin(ctx context.Context, login string) (*domain.User, error) {
	query := `SELECT id, tenant_id, login, password_hash, global_role, created_at, updated_at FROM users WHERE login = $1`
	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, login).Scan(&user.ID, &user.TenantID, &user.Login, &user.PasswordHash, &user.GlobalRole, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByTenantID(ctx context.Context, tenantID int64) ([]*domain.User, error) {
	query := `SELECT id, tenant_id, login, password_hash, global_role, created_at, updated_at FROM users WHERE tenant_id = $1`
	rows, err := r.db.Query(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.TenantID, &user.Login, &user.PasswordHash, &user.GlobalRole, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, rows.Err()
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `UPDATE users SET login = $1, password_hash = $2, global_role = $3, updated_at = CURRENT_TIMESTAMP WHERE id = $4 RETURNING updated_at`
	return r.db.QueryRow(ctx, query, user.Login, user.PasswordHash, user.GlobalRole, user.ID).Scan(&user.UpdatedAt)
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
