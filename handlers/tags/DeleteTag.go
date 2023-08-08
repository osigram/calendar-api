package tags

import (
	"calendar-api/handlers/events"
	"calendar-api/internal/userContext"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
	"net/http"
)

type DeleteRequestBody struct {
	Id      int64 `json:"id"`
	EventId int64 `json:"event_id"`
}

type Deleter interface {
	events.ByIdGetter
	DeleteTag(int64) error
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

		user, err := userContext.GetUser(r.Context())
		if err != nil {
			render.Status(r, 401)
			l.Debug("err in decoding json: " + err.Error())
			return
		}

		// TODO: add validation
		if requestBody.Id == 0 || requestBody.EventId == 0 {
			render.Status(r, 403)
			l.Debug("id is null")
			return
		}
		initialEvent, err := eventDeleter.GetEventById(requestBody.EventId)
		if err != nil {
			render.Status(r, 404)
			l.Error("err to get event from db: " + err.Error())
			return
		}
		if initialEvent.User.Email != user.Email {
			render.Status(r, 401)
			l.Debug("user is not equal")
			return
		}
		hasTag := false
		for _, tag := range initialEvent.Tags {
			if tag.Id == requestBody.Id {
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
		err = eventDeleter.DeleteTag(requestBody.Id)
		if err != nil {
			render.Status(r, 404)
			l.Debug("err to delete tag from db: " + err.Error())
			return
		}

		render.Status(r, 200)
	}
}
