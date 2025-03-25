package auth

import (
	"calendar-api/internal/core"
	"calendar-api/internal/pkg/context"
	"calendar-api/internal/pkg/errors"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"log/slog"
)

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

func (s *Service) Refresh(ctx *context.Context, requestBody RefreshRequest) (Response, error) {
	l := s.l.With(
		slog.String("op", "internal.auth.Refresh"),
	)

	user, session, err := s.validateRefreshToken(ctx, requestBody.RefreshToken)
	if err != nil {
		return Response{}, err
	}

	guid := uuid.New()
	session.RefreshToken = guid.String()
	session.DeviceData = ctx.UserAgent
	err = s.storage.UpdateSession(ctx, session)
	if err != nil {
		l.Error("failed to update session", slog.String("err", err.Error()))
		return Response{}, errors.NewInternalError("failed to update session")
	}

	return s.generateTokens(user, session)
}

func (s *Service) validateRefreshToken(ctx *context.Context, refreshToken string) (*core.User, *core.Session, error) {
	l := s.l.With(
		slog.String("op", "internal.auth.validateRefreshToken"),
	)

	token, err := jwtauth.VerifyToken(s.refreshTokenAuth, refreshToken)
	if err != nil {
		l.Debug("failed to decode refresh token", slog.String("err", err.Error()))
		return nil, nil, errors.NewValidationError("failed to decode refresh token")
	}

	claimsSession, err := core.NewSessionFromClaims(token.PrivateClaims())
	if err != nil {
		l.Debug("failed to decode session", slog.String("err", err.Error()))
		return nil, nil, errors.NewValidationError("failed to decode session")
	}

	user, err := s.storage.GetUser(ctx, claimsSession.UserEmail)
	if err != nil {
		l.Debug("failed to get user", slog.String("err", err.Error()))
		return nil, nil, errors.NewValidationError("failed to get user")
	}

	var session *core.Session
	for _, userSession := range user.Sessions {
		if userSession.ID == claimsSession.ID {
			session = &userSession
		}
	}
	if session == nil {
		return nil, nil, errors.NewValidationError("session not found")
	}
	if session.RefreshToken != claimsSession.RefreshToken {
		return nil, nil, errors.NewValidationError("invalid refresh token")
	}

	return user, session, nil
}
