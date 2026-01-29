package ctxutil

import "context"

type ContextKey string

const (
	UserIDCtxKey    ContextKey = "user-id"
	UserRolenCtxKey ContextKey = "user-role"
	UserTokenCtxKey ContextKey = "user-token"
)

func ContextWithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDCtxKey, userID)
}

func UserIDFromContext(ctx context.Context) string {
	userID, ok := ctx.Value(UserIDCtxKey).(string)
	if !ok {
		return ""
	}
	return userID
}

func ContextWithUserRole(ctx context.Context, userRole string) context.Context {
	return context.WithValue(ctx, UserRolenCtxKey, userRole)
}

func UserRoleFromContext(ctx context.Context) string {
	userRole, ok := ctx.Value(UserRolenCtxKey).(string)
	if !ok {
		return ""
	}
	return userRole
}

func ContextWithUserToken(ctx context.Context, userToken string) context.Context {
	return context.WithValue(ctx, UserTokenCtxKey, userToken)
}

func UserTokenFromContext(ctx context.Context) string {
	userToken, ok := ctx.Value(UserTokenCtxKey).(string)
	if !ok {
		return ""
	}
	return userToken
}
