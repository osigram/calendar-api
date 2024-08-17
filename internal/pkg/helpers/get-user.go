package helpers

import (
	"calendar-api/internal/core"
	"context"
	"github.com/go-chi/jwtauth/v5"
)

func GetUser(ctx context.Context) (*core.User, error) {
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	user, err := core.NewUserFromClaims(claims)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
