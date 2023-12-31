package events

import (
	"calendar-api/internal/usercontext"
	"calendar-api/types"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type Adder interface {
	AddEvent(event *types.Event) error
}

func Add(logger *slog.Logger, eventAdder Adder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.With(
			slog.String("op", "handlers.events.Add"),
			slog.String("requestId", middleware.GetReqID(r.Context())),
		)

		var requestBody types.Event
		err := render.DecodeJSON(r.Body, &requestBody)
		if err != nil {
			render.Status(r, 400)
			l.Debug("err in decoding json")
			return
		}

		user, err := usercontext.GetUser(r.Context())
		if err != nil {
			render.Status(r, 401)
			l.Debug(err.Error())
			return
		}

		if requestBody.ID != 0 || requestBody.SourceID != 0 || requestBody.Validate() != nil {
			render.Status(r, 403)
			return
		}
		requestBody.User = *user

		l.Info("adding event to db")
		err = eventAdder.AddEvent(&requestBody)
		if err != nil {
			render.Status(r, 500)
			l.Error("err to add event to db")
			return
		}

		render.Status(r, 200)
	}
}
