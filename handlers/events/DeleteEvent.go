package events

import (
	"calendar-api/internal/userContext"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
	"net/http"
)

type DeleteEventRequest struct {
	ID uint `json:"id"`
}

type Deleter interface {
	ByIDGetter
	DeleteEvent(uint) error
}

func Delete(logger *slog.Logger, eventDeleter Deleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.With(
			slog.String("op", "handlers.events.Delete"),
			slog.String("requestId", middleware.GetReqID(r.Context())),
		)

		var requestBody DeleteEventRequest
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
		if requestBody.ID == 0 {
			render.Status(r, 403)
			l.Debug("id is null")
			return
		}
		initialEvent, err := eventDeleter.GetEventByID(requestBody.ID)
		if err != nil {
			render.Status(r, 404)
			l.Error("err to get event from db")
			return
		}
		if initialEvent.User.Email != user.Email {
			render.Status(r, 401)
			l.Debug("user is not equal")
			return
		}

		l.Info("deleting event from db")
		err = eventDeleter.DeleteEvent(requestBody.ID)
		if err != nil {
			render.Status(r, 404)
			l.Debug("err to delete event from db")
			return
		}

		render.Status(r, 200)
	}
}
