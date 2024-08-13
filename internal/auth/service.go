package auth

import (
	"calendar-api/internal/core"
	"github.com/go-chi/jwtauth/v5"
	"log/slog"
)

type Storage interface {
	GetUser(email string) (*core.User, error)
	AddUser(user *core.User) error
}

type Service struct {
	l       *slog.Logger
	storage Storage
	jwtAuth *jwtauth.JWTAuth
}

func New(logger *slog.Logger, storage Storage, auth *jwtauth.JWTAuth) *Service {
	return &Service{
		l:       logger,
		storage: storage,
		jwtAuth: auth,
	}
}
