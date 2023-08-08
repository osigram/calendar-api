package tags

import (
	"calendar-api/handlers/events"
	"calendar-api/internal/userContext"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
	"net/http"
)

type RequestBody struct {
	TagText string `json:"tag_text"`
	EventId int64  `json:"event_id"`
}

type Adder interface {
	events.ByIdGetter
	AddTag(string, int64) error
}

func Add(logger *slog.Logger, tagAdder Adder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.With(
			slog.String("op", "handlers.tags.Add"),
			slog.String("requestId", middleware.GetReqID(r.Context())),
		)

		var requestBody RequestBody
		err := render.DecodeJSON(r.Body, &requestBody)
		if err != nil {
			render.Status(r, 400)
			l.Debug("err in decoding json: " + err.Error())
			return
		}

		user, err := userContext.GetUser(r.Context())
		if err != nil {
			render.Status(r, 401)
			l.Debug("err to get user: " + err.Error())
			return
		}

		// TODO: add validation
		if requestBody.EventId == 0 || requestBody.TagText == "" {
			render.Status(r, 403)
			l.Debug("validation error")
			return
		}

		initialEvent, err := tagAdder.GetEventById(requestBody.EventId)
		if err != nil {
			render.Status(r, 404)
			l.Debug("err to get event from db: " + err.Error())
			return
		}
		if initialEvent.User.Email != user.Email {
			render.Status(r, 401)
			l.Debug("user != user")
			return
		}

		l.Info("adding tag to db")
		err = tagAdder.AddTag(requestBody.TagText, requestBody.EventId)
		if err != nil {
			render.Status(r, 500)
			l.Error("err to add tag to db: " + err.Error())
			return
		}

		render.Status(r, 200)
	}
}
