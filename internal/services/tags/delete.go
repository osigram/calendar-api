package tags

import (
	"calendar-api/internal/core"
	"calendar-api/internal/errors"
	"calendar-api/internal/services/events"
	"log/slog"
)

type DeleteRequestBody struct {
	ID      uint `json:"ID"`
	EventID uint `json:"eventID"`
}

type Deleter interface {
	events.ByIDGetter
	DeleteTag(uint) error
}

func (s *Service) Delete(user *core.User, requestBody DeleteRequestBody) error {
	l := s.l.With(
		slog.String("op", "services.tags.Delete"),
	)

	// TODO: add validation
	if requestBody.ID == 0 || requestBody.EventID == 0 {
		l.Debug("validation error: requestBodyID is zero or eventID is zero")
		return errors.NewValidationError("requestBodyID is zero or eventID is zero")
	}
	initialEvent, err := s.storage.GetEventByID(requestBody.EventID)
	if err != nil {
		l.Error("unable to get event from db", slog.String("err", err.Error()))
		return errors.NewNotFoundError("event not found")
	}
	if initialEvent.UserEmail != user.Email {
		l.Debug("user emails not match", slog.String("userEmail", user.Email), slog.String("eventEmail", initialEvent.UserEmail))
		return errors.NewAuthError("wrong user")
	}

	l.Debug("checking if the event has such tag")
	var hasTag bool
	for _, tag := range initialEvent.Tags {
		if tag.ID == requestBody.ID {
			hasTag = true
			break
		}
	}
	if !hasTag {
		l.Debug("validation error: the event does not have such tag")
		return errors.NewValidationError("the event does not have such tag")
	}

	l.Info("deleting tag from db")
	err = s.storage.DeleteTag(requestBody.ID)
	if err != nil {
		l.Error("unable to delete tag from db", slog.String("err", err.Error()))
		return errors.NewInternalError("tag not found")
	}

	return nil
}
