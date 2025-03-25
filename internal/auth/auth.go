package auth

import (
	"calendar-api/internal/core"
	"calendar-api/internal/pkg/context"
	"calendar-api/internal/pkg/errors"
	"github.com/google/uuid"
	"log/slog"
)

func (s *Service) Auth(ctx *context.Context, claimsUser core.User) (Response, error) {
	l := s.l.With(
		slog.String("op", "internal.auth.Auth"),
	)

	l.Debug("getting user from db")
	user, err := s.storage.GetUser(ctx, claimsUser.Email)
	if err != nil {
		l.Debug("failed to get user from db, registering user")
		if err = s.storage.AddUser(ctx, &claimsUser); err != nil {
			l.Error("unable to add user to db", slog.String("err", err.Error()))
			return Response{}, errors.NewInternalError("internal registration error")
		}
		user = &claimsUser
	} else {
		l.Debug("updating user data")
		user.Name = claimsUser.Name
		user.Picture = claimsUser.Picture
		err = s.storage.UpdateUser(ctx, user)
		if err != nil {
			l.Error("unable to update user data", slog.String("err", err.Error()))
			return Response{}, errors.NewInternalError("internal registration error")
		}
	}

	l.Debug("creating session")
	guid := uuid.New()
	session := &core.Session{
		UserEmail:    user.Email,
		RefreshToken: guid.String(),
		DeviceData:   ctx.UserAgent,
	}
	err = s.storage.AddSession(ctx, session)
	if err != nil {
		l.Error("unable to add session to db", slog.String("err", err.Error()))
		return Response{}, errors.NewInternalError("internal auth error")
	}

	return s.generateTokens(user, session)
}

func (s *Service) generateTokens(user *core.User, session *core.Session) (Response, error) {
	l := s.l.With(
		slog.String("op", "internal.auth.generateTokens"),
	)

	l.Debug("generating refresh token")
	_, refreshTokenString, err := s.refreshTokenAuth.Encode(session.GetClaims())
	if err != nil {
		l.Error("unable to encode refresh token", slog.String("err", err.Error()))
		return Response{}, errors.NewInternalError("internal auth error")
	}

	l.Debug("generating access token")
	_, accessTokenString, err := s.accessTokenAuth.Encode(user.GetClaims())
	if err != nil {
		l.Error("unable to generate access token", slog.String("err", err.Error()))
		return Response{}, errors.NewInternalError("internal error")
	}

	return Response{AccessToken: accessTokenString, RefreshToken: refreshTokenString}, nil
}
