package events

import (
	"calendar-api/internal/extensions"
	"calendar-api/internal/userContext"
	"calendar-api/types"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
	"net/http"
	"strconv"
	"time"
)

var DefaultTime time.Time

type GetEventByDateRequest struct {
	TimeOfStart  time.Time `json:"timeOfStart"`
	TimeOfFinish time.Time `json:"timeOfFinish,omitempty"`
}

type ByDateGetter interface {
	GetEventsByDate(*types.User, time.Time, time.Time) ([]types.Event, error)
}

type ExtensionGetter interface {
	Get(id int64) (extensions.Extension, error)
}

func GetByDate(logger *slog.Logger, eventGetter ByDateGetter, extensionMapper ExtensionGetter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := logger.With(
			slog.String("op", "handlers.events.GetByDate"),
			slog.String("requestId", middleware.GetReqID(r.Context())),
		)

		var requestBody GetEventByDateRequest
		err := render.DecodeJSON(r.Body, &requestBody)
		if err != nil {
			render.Status(r, 400)
			l.Debug("err in decoding json")
			return
		}

		user, err := userContext.GetUser(r.Context())
		if err != nil {
			render.Status(r, 401)
			l.Debug(err.Error())
		}

		l.Debug("getting events from db")

		if requestBody.TimeOfStart == DefaultTime {
			requestBody.TimeOfStart = time.UnixMicro(0)
		}
		if requestBody.TimeOfFinish == DefaultTime {
			requestBody.TimeOfFinish = time.Now().Add(52 * 365 * 24 * time.Hour)
		}
		events, err := eventGetter.GetEventsByDate(user, requestBody.TimeOfStart, requestBody.TimeOfFinish)
		if err != nil {
			render.Status(r, 404)
			l.Error("err to get events from db")
			return
		}

		l.Debug("getting events from extensions")

		for _, extensionData := range user.ExtensionsUsed {
			extension, err := extensionMapper.Get(extensionData.Id)
			if err != nil {
				l.Error("ExtensionData " + strconv.Itoa(int(extensionData.Id)) + "is not implemented")
				continue
			}

			extensionEvents, err := extension.GetEventsByDate(user, requestBody.TimeOfStart, requestBody.TimeOfFinish)
			if err != nil {
				l.Error("err to get events from extension")
				continue
			}

			events = append(events, extensionEvents...)
		}

		render.Status(r, 200)
		render.JSON(w, r, events)
	})
}
