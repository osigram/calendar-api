package auth

import (
	"calendar-api/internal/core"
	"github.com/go-chi/jwtauth/v5"
	"log/slog"
)

type Response struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type Storage interface {
	GetUser(email string) (*core.User, error)
	AddUser(user *core.User) error
	AddSession(session *core.Session) error
	UpdateUser(user *core.User) error
	UpdateSession(session *core.Session) error
	DeleteSession(sessionID uint) error
}

type Service struct {
	l                *slog.Logger
	storage          Storage
	accessTokenAuth  *jwtauth.JWTAuth
	refreshTokenAuth *jwtauth.JWTAuth
	googleClientID   string
}

func New(logger *slog.Logger, storage Storage, auth *jwtauth.JWTAuth, googleClientID string) *Service {
	return &Service{
		l:                logger,
		storage:          storage,
		accessTokenAuth:  auth,
		refreshTokenAuth: auth,
		googleClientID:   googleClientID,
	}
}
