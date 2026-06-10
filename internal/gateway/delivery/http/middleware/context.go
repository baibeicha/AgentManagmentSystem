package middleware

import "context"

type ctxKey string

const (
	ContextKeyUserID   ctxKey = "user_id"
	ContextKeyTenantID ctxKey = "tenant_id"
	ContextKeyRole     ctxKey = "role"
)

func WithUserContext(ctx context.Context, userID, tenantID, role string) context.Context {
	ctx = context.WithValue(ctx, ContextKeyUserID, userID)
	ctx = context.WithValue(ctx, ContextKeyTenantID, tenantID)
	ctx = context.WithValue(ctx, ContextKeyRole, role)
	return ctx
}

func getValue(ctx context.Context, key ctxKey) (string, bool) {
	val := ctx.Value(key)
	str, ok := val.(string)
	return str, ok
}

func GetUserID(ctx context.Context) (string, bool) {
	return getValue(ctx, ContextKeyUserID)
}

func GetTenantID(ctx context.Context) (string, bool) {
	return getValue(ctx, ContextKeyTenantID)
}

func GetRole(ctx context.Context) (string, bool) {
	return getValue(ctx, ContextKeyRole)
}
