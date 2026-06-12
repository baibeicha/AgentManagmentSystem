package repository

import (
	"AgentManagmentSystem/internal/auth/domain"
	"AgentManagmentSystem/pkg/storage"
	"context"
)

type ResourcePolicyRepository struct {
	db *storage.DB
}

func NewResourcePolicyRepository(db *storage.DB) *ResourcePolicyRepository {
	return &ResourcePolicyRepository{db: db}
}

func (r *ResourcePolicyRepository) AddPolicy(ctx context.Context, policy *domain.ResourcePolicy) error {
	query := `
		INSERT INTO resource_policies (user_id, resource_id, action) 
		VALUES ($1, $2, $3) RETURNING id, created_at`
	return r.db.QueryRow(ctx, query, policy.UserID, policy.ResourceID, policy.Action).Scan(&policy.ID, &policy.CreatedAt)
}

func (r *ResourcePolicyRepository) RemovePolicy(ctx context.Context, userID string, resourceID string, action string) error {
	query := `DELETE FROM resource_policies WHERE user_id = $1 AND resource_id = $2 AND action = $3`
	_, err := r.db.Exec(ctx, query, userID, resourceID, action)
	return err
}

func (r *ResourcePolicyRepository) GetPoliciesByUserID(ctx context.Context, userID string) ([]*domain.ResourcePolicy, error) {
	query := `SELECT id, user_id, resource_id, action, created_at FROM resource_policies WHERE user_id = $1`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var policies []*domain.ResourcePolicy
	for rows.Next() {
		var policy domain.ResourcePolicy
		if err := rows.Scan(&policy.ID, &policy.UserID, &policy.ResourceID, &policy.Action, &policy.CreatedAt); err != nil {
			return nil, err
		}
		policies = append(policies, &policy)
	}
	return policies, rows.Err()
}

func (r *ResourcePolicyRepository) CheckPermission(ctx context.Context, userID string, resourceID string, action string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM resource_policies WHERE user_id = $1 AND resource_id = $2 AND action = $3)`
	var exists bool
	err := r.db.QueryRow(ctx, query, userID, resourceID, action).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
