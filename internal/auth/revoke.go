package auth

import (
	"calendar-api/internal/context"
	"errors"
	"log/slog"
)

func (s *Service) Revoke(_ *context.Context, requestBody RefreshRequest) error {
	l := s.l.With(
		slog.String("op", "internal.auth.Revoke"),
	)

	_, session, err := s.validateRefreshToken(requestBody.RefreshToken)
	if err != nil {
		return err
	}

	err = s.storage.DeleteSession(session.ID)
	if err != nil {
		l.Error("failed to delete session", slog.String("err", err.Error()))
		return errors.New("failed to delete session")
	}

	return nil
}
