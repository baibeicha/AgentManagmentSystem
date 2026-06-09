package repository

import (
	"AgentManagmentSystem/internal/auth/domain"
	"AgentManagmentSystem/pkg/storage"
	"context"
)

type GroupPolicyRepository struct {
	db *storage.DB
}

func NewGroupPolicyRepository(db *storage.DB) *GroupPolicyRepository {
	return &GroupPolicyRepository{db: db}
}

func (r *GroupPolicyRepository) AddPolicy(ctx context.Context, policy *domain.GroupPolicy) error {
	query := `INSERT INTO group_policies (user_id, group_id, permission) VALUES ($1, $2, $3) RETURNING created_at`
	return r.db.QueryRow(ctx, query, policy.UserID, policy.GroupID, policy.Permission).Scan(&policy.CreatedAt)
}

func (r *GroupPolicyRepository) RemovePolicy(ctx context.Context, userID int64, groupID int64, permission string) error {
	query := `DELETE FROM group_policies WHERE user_id = $1 AND group_id = $2 AND permission = $3`
	_, err := r.db.Exec(ctx, query, userID, groupID, permission)
	return err
}

func (r *GroupPolicyRepository) GetPoliciesByUserID(ctx context.Context, userID int64) ([]*domain.GroupPolicy, error) {
	query := `SELECT user_id, group_id, permission, created_at FROM group_policies WHERE user_id = $1`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var policies []*domain.GroupPolicy
	for rows.Next() {
		var policy domain.GroupPolicy
		if err := rows.Scan(&policy.UserID, &policy.GroupID, &policy.Permission, &policy.CreatedAt); err != nil {
			return nil, err
		}
		policies = append(policies, &policy)
	}
	return policies, rows.Err()
}

func (r *GroupPolicyRepository) CheckPermission(ctx context.Context, userID int64, groupID int64, permission string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM group_policies WHERE user_id = $1 AND group_id = $2 AND permission = $3)`
	var exists bool
	err := r.db.QueryRow(ctx, query, userID, groupID, permission).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
