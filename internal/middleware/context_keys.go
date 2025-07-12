package middleware

import "context"

type contextKey string

const (
	userIDKey   contextKey = "user_id"
	userRoleKey contextKey = "user_role"
)

func GetUserIdFromContext(ctx context.Context) (int64, bool) {
	v := ctx.Value(userIDKey)
	id, ok := v.(int64)
	return id, ok
}

func GetUserRoleFromContext(ctx context.Context) (string, bool) {
	v := ctx.Value(userRoleKey)
	role, ok := v.(string)
	return role, ok
}
