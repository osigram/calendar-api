package events

import (
	"calendar-api/internal/core"
	"calendar-api/internal/errors"
	"log/slog"
)

type GetEventByIDRequest struct {
	ID     uint `json:"id"`
	Source uint `json:"source,omitempty"`
}

type ByIDGetter interface {
	GetEventByID(uint) (*core.Event, error)
}

func (s *Service) GetByID(user *core.User, requestBody GetEventByIDRequest) (core.Event, error) {
	l := s.l.With(
		slog.String("op", "services.events.GetById"),
	)

	if requestBody.ID == 0 {
		l.Debug("validation error: id is zero")
		return core.Event{}, errors.NewValidationError("id is zero")
	}

	var eventGetter ByIDGetter
	switch requestBody.Source {
	case 0:
		eventGetter = s.storage
	default:
		if !hasUserExtension(user, requestBody.Source) {
			l.Debug("validation error: source is invalid")
			return core.Event{}, errors.NewValidationError("source is invalid")
		}

		extension, err := s.extensionsGetter.Get(requestBody.Source)
		if err != nil {
			l.Error("extension is not implemented", slog.Uint64("extensionID", uint64(requestBody.Source)))
			return core.Event{}, errors.NewInternalError("extension is not implemented")
		}

		eventGetter = extension
	}

	l.Info("getting event from db")
	event, err := eventGetter.GetEventByID(requestBody.ID)
	if err != nil {
		l.Debug("unable to get event from db", slog.String("err", err.Error()))
		return core.Event{}, errors.NewNotFoundError("event not found")
	}

	if user.Email != event.UserEmail && requestBody.Source == 0 {
		l.Debug("user is not an owner of this event")
		return core.Event{}, errors.NewAuthError("user is not an owner of this event")
	}

	return *event, nil
}

func hasUserExtension(user *core.User, id uint) bool {
	for _, ed := range user.ExtensionsData {
		if id == ed.ID {
			return true
		}
	}

	return false
}
