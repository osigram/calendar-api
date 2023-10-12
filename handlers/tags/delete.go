package tags

import (
	"calendar-api/handlers/events"
	"calendar-api/internal/usercontext"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type DeleteRequestBody struct {
	ID      uint `json:"ID"`
	EventID uint `json:"eventID"`
}

type Deleter interface {
	events.ByIDGetter
	DeleteTag(uint) error
}

func Delete(logger *slog.Logger, eventDeleter Deleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.With(
			slog.String("op", "handlers.tags.Delete"),
			slog.String("requestId", middleware.GetReqID(r.Context())),
		)
		var requestBody DeleteRequestBody
		err := render.DecodeJSON(r.Body, &requestBody)
		if err != nil {
			render.Status(r, 400)
			l.Debug("err in decoding json: " + err.Error())
			return
		}

		user, err := usercontext.GetUser(r.Context())
		if err != nil {
			render.Status(r, 401)
			l.Debug("err in decoding json: " + err.Error())
			return
		}

		// TODO: add validation
		if requestBody.ID == 0 || requestBody.EventID == 0 {
			render.Status(r, 403)
			l.Debug("id is null")
			return
		}
		initialEvent, err := eventDeleter.GetEventByID(requestBody.EventID)
		if err != nil {
			render.Status(r, 404)
			l.Error("err to get event from db: " + err.Error())
			return
		}
		if initialEvent.UserEmail != user.Email {
			render.Status(r, 401)
			l.Debug("user is not equal")
			return
		}
		hasTag := false
		for _, tag := range initialEvent.Tags {
			if tag.ID == requestBody.ID {
				hasTag = true
				break
			}
		}
		if !hasTag {
			render.Status(r, 403)
			l.Debug("event does not have such tag")
			return
		}

		l.Info("deleting tag from db")
		err = eventDeleter.DeleteTag(requestBody.ID)
		if err != nil {
			render.Status(r, 404)
			l.Debug("err to delete tag from db: " + err.Error())
			return
		}

		render.Status(r, 200)
	}
}
