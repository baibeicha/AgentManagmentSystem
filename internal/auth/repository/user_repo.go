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
	query := `
		INSERT INTO users (tenant_id, email, password_hash, role, status, is_2fa_enabled, totp_secret) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) 
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(
		ctx, query,
		user.TenantID, user.Email, user.PasswordHash, user.Role, user.Status, user.Is2FAEnabled, user.TOTPSecret,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	query := `
		SELECT id, tenant_id, email, password_hash, role, status, is_2fa_enabled, COALESCE(totp_secret, ''), created_at, updated_at 
		FROM users WHERE id = $1`
	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.TenantID, &user.Email, &user.PasswordHash,
		&user.Role, &user.Status, &user.Is2FAEnabled, &user.TOTPSecret,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, tenant_id, email, password_hash, role, status, is_2fa_enabled, COALESCE(totp_secret, ''), created_at, updated_at 
		FROM users WHERE email = $1`
	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.TenantID, &user.Email, &user.PasswordHash,
		&user.Role, &user.Status, &user.Is2FAEnabled, &user.TOTPSecret,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByTenantID(ctx context.Context, tenantID string) ([]*domain.User, error) {
	query := `
		SELECT id, tenant_id, email, password_hash, role, status, is_2fa_enabled, COALESCE(totp_secret, ''), created_at, updated_at 
		FROM users WHERE tenant_id = $1`
	rows, err := r.db.Query(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(
			&user.ID, &user.TenantID, &user.Email, &user.PasswordHash,
			&user.Role, &user.Status, &user.Is2FAEnabled, &user.TOTPSecret,
			&user.CreatedAt, &user.UpdatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, rows.Err()
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users 
		SET email = $1, password_hash = $2, role = $3, status = $4, is_2fa_enabled = $5, totp_secret = $6, updated_at = CURRENT_TIMESTAMP 
		WHERE id = $7 RETURNING updated_at`
	return r.db.QueryRow(
		ctx, query,
		user.Email, user.PasswordHash, user.Role, user.Status, user.Is2FAEnabled, user.TOTPSecret, user.ID,
	).Scan(&user.UpdatedAt)
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
