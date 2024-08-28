package auth

import (
	"calendar-api/internal/core"
	"context"
	"github.com/go-chi/jwtauth/v5"
	"log/slog"
)

type Response struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type Storage interface {
	GetUser(ctx context.Context, email string) (*core.User, error)
	AddUser(ctx context.Context, user *core.User) error
	AddSession(ctx context.Context, session *core.Session) error
	UpdateUser(ctx context.Context, user *core.User) error
	UpdateSession(ctx context.Context, session *core.Session) error
	DeleteSession(ctx context.Context, sessionID uint) error
}

type Service struct {
	l                *slog.Logger
	storage          Storage
	accessTokenAuth  *jwtauth.JWTAuth
	refreshTokenAuth *jwtauth.JWTAuth
	googleClientID   string
}

func New(logger *slog.Logger, storage Storage, accessTokenAuth, refreshTokenAuth *jwtauth.JWTAuth, googleClientID string) *Service {
	return &Service{
		l:                logger,
		storage:          storage,
		accessTokenAuth:  accessTokenAuth,
		refreshTokenAuth: refreshTokenAuth,
		googleClientID:   googleClientID,
	}
}
