package events

import (
	"calendar-api/internal/core"
	"calendar-api/internal/pkg/context"
	"calendar-api/internal/pkg/errors"
	"log/slog"
)

type Updater interface {
	ByIDGetter
	UpdateEvent(event *core.Event) error
}

func (s *Service) Update(ctx *context.Context, requestBody core.Event) error {
	l := s.l.With(
		slog.String("op", "services.events.Update"),
	)

	if requestBody.ID != 0 || requestBody.SourceID != 0 {
		l.Debug("validation error: id is zero")
		return errors.NewValidationError("id must be zero")
	}

	if err := requestBody.Validate(); err != nil {
		l.Debug("validation error: request body is invalid", slog.String("err", err.Error()))
		return errors.NewValidationError(err.Error())
	}

	initialEvent, err := s.storage.GetEventByID(requestBody.ID)
	if err != nil {
		l.Debug("event does not exist")
		return errors.NewNotFoundError("event not found")
	}
	if initialEvent.UserEmail != ctx.User.Email {
		l.Debug("user email does not match")
		return errors.NewAuthError("wrong user")
	}

	l.Info("updating event in db")
	err = s.storage.UpdateEvent(&requestBody)
	if err != nil {
		l.Error("unable to update event in db")
		return errors.NewInternalError("unable to update event in db")
	}

	return nil
}
