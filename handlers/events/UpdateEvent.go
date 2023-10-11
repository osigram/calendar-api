package events

import (
	"calendar-api/internal/userContext"
	"calendar-api/types"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
	"net/http"
)

type Updater interface {
	ByIDGetter
	UpdateEvent(event *types.Event) error
}

func Update(logger *slog.Logger, eventUpdater Updater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.With(
			slog.String("op", "handlers.events.Update"),
			slog.String("requestId", middleware.GetReqID(r.Context())),
		)

		var requestBody types.Event
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
			return
		}

		// TODO: add validation
		if requestBody.ID == 0 || requestBody.SourceID != 0 {
			render.Status(r, 403)
			return
		}
		initialEvent, err := eventUpdater.GetEventByID(requestBody.ID)
		if err != nil {
			render.Status(r, 404)
			return
		}
		if initialEvent.User.Email != user.Email {
			render.Status(r, 401)
			return
		}

		l.Info("updating event in db")
		err = eventUpdater.UpdateEvent(&requestBody)
		if err != nil {
			render.Status(r, 404)
			l.Debug("err to update event in db")
			return
		}

		render.Status(r, 200)
	}
}
