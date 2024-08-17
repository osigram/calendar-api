package events

import (
	"calendar-api/internal/core"
	"calendar-api/internal/extensions"
	"calendar-api/internal/pkg/context"
	"calendar-api/internal/pkg/errors"
	"log/slog"
	"time"
)

type GetEventByDateRequest struct {
	TimeOfStart  time.Time `schema:"timeOfStart"`
	TimeOfFinish time.Time `schema:"timeOfFinish,omitempty"`
}

type ByDateGetter interface {
	GetEventsByDate(*core.User, time.Time, time.Time) ([]core.Event, error)
}

type UserGetter interface {
	GetUser(string) (*core.User, error)
}

type ExtensionGetter interface {
	Get(id uint) (extensions.Extension, error)
}

func (s *Service) GetByDate(ctx *context.Context, requestBody GetEventByDateRequest) ([]core.Event, error) {
	l := s.l.With(
		slog.String("op", "services.events.GetByDate"),
	)

	l.Debug("getting events from db")
	if requestBody.TimeOfStart.IsZero() {
		requestBody.TimeOfStart = time.UnixMicro(0)
	}
	if requestBody.TimeOfFinish.IsZero() {
		requestBody.TimeOfFinish = time.Now().Add(52 * 365 * 24 * time.Hour)
	}
	events, err := s.storage.GetEventsByDate(ctx.User, requestBody.TimeOfStart, requestBody.TimeOfFinish)
	if err != nil {
		l.Error("failed to get events by date", slog.String("err", err.Error()))
		return nil, errors.NewInternalError("unable to get events from db")
	}

	user, err := s.storage.GetUser(ctx.User.Email)
	if err != nil {
		l.Error("failed to get user from db", slog.String("err", err.Error()))
		return nil, errors.NewInternalError("unable to get user from db")
	}
	l.Debug("getting events from extensions")
	for _, extensionData := range user.ExtensionsData {
		extension, err := s.extensionsGetter.Get(extensionData.Extension)
		if err != nil {
			l.Error("Extension is not implemented", slog.Uint64("extensionID", uint64(extensionData.ID)), slog.String("err", err.Error()))
			continue
		}

		extensionEvents, err := extension.GetEventsByDate(extensionData.AdditionalData, requestBody.TimeOfStart, requestBody.TimeOfFinish)
		if err != nil {
			l.Error("unable to get events from extension", slog.Uint64("extensionID", uint64(extensionData.ID)), slog.String("err", err.Error()))
			continue
		}

		events = append(events, extensionEvents...)
	}

	return events, nil
}
