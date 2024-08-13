package events

import (
	"calendar-api/internal/context"
	"calendar-api/internal/errors"
	"log/slog"
)

type DeleteEventRequest struct {
	ID uint `json:"id"`
}

type Deleter interface {
	ByIDGetter
	DeleteEvent(uint) error
}

func (s *Service) Delete(ctx *context.Context, requestBody DeleteEventRequest) error {
	l := s.l.With(
		slog.String("op", "services.events.Delete"),
	)

	if requestBody.ID == 0 {
		l.Debug("validation error: id is zero")
		return errors.NewValidationError("id is zero")
	}

	initialEvent, err := s.storage.GetEventByID(requestBody.ID)
	if err != nil {
		l.Debug("unable to get event from db")
		return errors.NewNotFoundError("unable to get event from db")
	}
	if initialEvent.UserEmail != ctx.User.Email {
		l.Debug("user email not match")
		return errors.NewAuthError("wrong user")
	}

	l.Info("deleting event from db")
	err = s.storage.DeleteEvent(requestBody.ID)
	if err != nil {
		l.Debug("unable to delete event from db")
		return errors.NewNotFoundError("unable to delete event from db")
	}

	return nil
}
