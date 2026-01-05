package context

import (
	"context"

	"github.com/raminderis/lenslocked/models"
)

type key string

const (
	UserKey key = "user"
)

func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, UserKey, user)
}

func User(ctx context.Context) *models.User {
	val := ctx.Value(UserKey)
	user, ok := val.(*models.User)
	if !ok {
		return nil
	}
	return user
}
