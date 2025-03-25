package events

import (
	pkgcontext "calendar-api/internal/pkg/context"
	"calendar-api/internal/pkg/errors"
	"context"
	"log/slog"
)

type DeleteEventRequest struct {
	ID uint `json:"id"`
}

type Deleter interface {
	ByIDGetter
	DeleteEvent(context.Context, uint) error
}

func (s *Service) Delete(ctx *pkgcontext.Context, requestBody DeleteEventRequest) error {
	l := s.l.With(
		slog.String("op", "services.events.Delete"),
	)

	if requestBody.ID == 0 {
		l.Debug("validation error: id is zero")
		return errors.NewValidationError("id is zero")
	}

	initialEvent, err := s.storage.GetEventByID(ctx, requestBody.ID)
	if err != nil {
		l.Debug("unable to get event from db")
		return errors.NewNotFoundError("unable to get event from db")
	}
	if initialEvent.UserEmail != ctx.User.Email {
		l.Debug("user email not match")
		return errors.NewAuthError("wrong user")
	}

	l.Info("deleting event from db")
	err = s.storage.DeleteEvent(ctx, requestBody.ID)
	if err != nil {
		l.Debug("unable to delete event from db")
		return errors.NewNotFoundError("unable to delete event from db")
	}

	return nil
}
