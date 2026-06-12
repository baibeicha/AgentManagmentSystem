package repository

import (
	"AgentManagmentSystem/internal/auth/domain"
	"AgentManagmentSystem/pkg/storage"
	"context"
)

type SessionRepository struct {
	db *storage.DB
}

func NewSessionRepository(db *storage.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(ctx context.Context, session *domain.Session) error {
	query := `
		INSERT INTO sessions (user_id, refresh_token, device_id, client_ip, user_agent, is_revoked, expires_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) 
		RETURNING id, created_at`

	return r.db.QueryRow(
		ctx, query,
		session.UserID, session.RefreshToken, session.DeviceID,
		session.ClientIP, session.UserAgent, session.IsRevoked, session.ExpiresAt,
	).Scan(&session.ID, &session.CreatedAt)
}

func (r *SessionRepository) GetByRefreshToken(ctx context.Context, token string) (*domain.Session, error) {
	query := `
		SELECT id, user_id, refresh_token, device_id, client_ip, user_agent, is_revoked, expires_at, created_at 
		FROM sessions WHERE refresh_token = $1`

	session := &domain.Session{}
	err := r.db.QueryRow(ctx, query, token).Scan(
		&session.ID, &session.UserID, &session.RefreshToken, &session.DeviceID,
		&session.ClientIP, &session.UserAgent, &session.IsRevoked,
		&session.ExpiresAt, &session.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (r *SessionRepository) RevokeSession(ctx context.Context, id string) error {
	query := `UPDATE sessions SET is_revoked = true WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *SessionRepository) RevokeAllUserSessions(ctx context.Context, userID string) error {
	query := `UPDATE sessions SET is_revoked = true WHERE user_id = $1 AND is_revoked = false`
	_, err := r.db.Exec(ctx, query, userID)
	return err
}
