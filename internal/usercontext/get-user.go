package usercontext

import (
	"calendar-api/types"
	"context"
	"errors"
)

func GetUser(ctx context.Context) (*types.User, error) {
	userAny := ctx.Value("user")
	if user, ok := userAny.(*types.User); ok {
		return user, nil
	}

	return nil, errors.New("cannot get user from context")
}
