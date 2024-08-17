package events

import (
	"calendar-api/internal/core"
	"calendar-api/internal/pkg/context"
	"calendar-api/internal/pkg/errors"
	"log/slog"
)

type Adder interface {
	AddEvent(event *core.Event) error
}

func (s *Service) Add(ctx *context.Context, requestBody core.Event) error {
	l := s.l.With(
		slog.String("op", "services.events.Add"),
	)

	if requestBody.ID != 0 || requestBody.SourceID != 0 {
		l.Debug("validation error: id is zero")
		return errors.NewValidationError("id must be zero")
	}

	if err := requestBody.Validate(); err != nil {
		l.Debug("validation error: request body is invalid", slog.String("err", err.Error()))
		return errors.NewValidationError(err.Error())
	}

	requestBody.User = *ctx.User

	l.Info("adding event to db")
	err := s.storage.AddEvent(&requestBody)
	if err != nil {
		l.Error("unable to add event to db")
		return errors.NewInternalError("unable to add event to db")
	}

	return nil
}
