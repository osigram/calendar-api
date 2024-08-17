package auth

import (
	"calendar-api/internal/context"
	"calendar-api/internal/core"
	"calendar-api/internal/errors"
	"google.golang.org/api/idtoken"
	"log/slog"
)

type LoginData struct {
	GoogleJWT string `json:"googleJWT"`
}

func (s *Service) Google(ctx *context.Context, loginData LoginData) (Response, error) {
	l := s.l.With(
		slog.String("op", "internal.auth.Google"),
	)

	l.Debug("processing jwt from Google")
	payload, err := idtoken.Validate(ctx, loginData.GoogleJWT, s.googleClientID)
	if err != nil {
		l.Debug("unable to parse jwt from Google", slog.String("err", err.Error()))
		return Response{}, errors.NewValidationError("invalid AccessToken token")
	}

	l.Debug("processing claims")
	claimsUser, err := core.NewUserFromGoogleClaims(payload.Claims)
	if err != nil {
		l.Error("unable to parse claims from Google", slog.String("err", err.Error()))
		return Response{}, errors.NewInternalError("internal auth error")
	}

	return s.Auth(ctx, claimsUser)
}
