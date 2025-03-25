package auth

import (
	"calendar-api/internal/pkg/context"
	"errors"
	"log/slog"
)

func (s *Service) Revoke(ctx *context.Context, requestBody RefreshRequest) error {
	l := s.l.With(
		slog.String("op", "internal.auth.Revoke"),
	)

	_, session, err := s.validateRefreshToken(ctx, requestBody.RefreshToken)
	if err != nil {
		return err
	}

	err = s.storage.DeleteSession(ctx, session.ID)
	if err != nil {
		l.Error("failed to delete session", slog.String("err", err.Error()))
		return errors.New("failed to delete session")
	}

	return nil
}
