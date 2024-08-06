package helpers

import (
	"calendar-api/internal/core"
	"context"
	"errors"
)

func GetUser(ctx context.Context) (*core.User, error) {
	userAny := ctx.Value("user")
	if user, ok := userAny.(*core.User); ok {
		return user, nil
	}

	return nil, errors.New("cannot get user from context")
}
