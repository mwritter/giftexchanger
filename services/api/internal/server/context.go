package server

import (
	"context"

	"github.com/mwritter/giftexchanger/services/api/internal/auth"
)

type userContextKey struct{}

func WithUser(ctx context.Context, user *auth.User) context.Context {
	return context.WithValue(ctx, userContextKey{}, user)
}

func UserFromContext(ctx context.Context) *auth.User {
	user, _ := ctx.Value(userContextKey{}).(*auth.User)
	return user
}
