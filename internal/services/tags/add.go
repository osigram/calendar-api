package tags

import (
	"calendar-api/internal/core"
	pkgcontext "calendar-api/internal/pkg/context"
	"calendar-api/internal/pkg/errors"
	"calendar-api/internal/services/events"
	"context"
	"log/slog"
)

type Adder interface {
	events.ByIDGetter
	AddTag(context.Context, string, uint) error
}

func (s *Service) Add(ctx *pkgcontext.Context, requestBody core.Tag) error {
	l := s.l.With(
		slog.String("op", "services.tags.Add"),
	)

	if requestBody.EventID == 0 || requestBody.ID != 0 {
		l.Debug("validation error: eventID is zero or tagID is not zero")
		return errors.NewValidationError("eventID is zero or tagID is not zero")
	}

	if err := requestBody.Validate(); err != nil {
		l.Debug("validation error: request body is invalid", slog.String("err", err.Error()))
		return errors.NewValidationError("request body is invalid: " + err.Error())
	}

	initialEvent, err := s.storage.GetEventByID(ctx, requestBody.EventID)
	if err != nil {
		l.Debug("unable to get event from db", slog.String("err", err.Error()))
		return errors.NewNotFoundError("event not found")
	}
	if initialEvent.UserEmail != ctx.User.Email {
		l.Debug("user emails not match", slog.String("userEmail", ctx.User.Email), slog.String("eventEmail", initialEvent.UserEmail))
		return errors.NewAuthError("wrong user")
	}

	l.Info("adding tag to db")
	err = s.storage.AddTag(ctx, requestBody.TagText, requestBody.EventID)
	if err != nil {
		l.Error("unable to add tag to db", slog.String("err", err.Error()))
		return errors.NewInternalError("unable to add tag to db")
	}

	return nil
}
