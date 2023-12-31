package tags

import (
	"calendar-api/handlers/events"
	"calendar-api/internal/usercontext"
	"calendar-api/types"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type Adder interface {
	events.ByIDGetter
	AddTag(string, uint) error
}

func Add(logger *slog.Logger, tagAdder Adder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.With(
			slog.String("op", "handlers.tags.Add"),
			slog.String("requestId", middleware.GetReqID(r.Context())),
		)

		var requestBody types.Tag
		err := render.DecodeJSON(r.Body, &requestBody)
		if err != nil {
			render.Status(r, 400)
			l.Debug("err in decoding json: " + err.Error())
			return
		}

		user, err := usercontext.GetUser(r.Context())
		if err != nil {
			render.Status(r, 401)
			l.Debug("err to get user: " + err.Error())
			return
		}

		if requestBody.EventID == 0 || requestBody.ID != 0 || requestBody.TagText == "" || requestBody.Validate() != nil {
			render.Status(r, 403)
			l.Debug("validation error")
			return
		}

		initialEvent, err := tagAdder.GetEventByID(requestBody.EventID)
		if err != nil {
			render.Status(r, 404)
			l.Debug("err to get event from db: " + err.Error())
			return
		}
		if initialEvent.UserEmail != user.Email {
			render.Status(r, 401)
			l.Debug("user != user")
			return
		}

		l.Info("adding tag to db")
		err = tagAdder.AddTag(requestBody.TagText, requestBody.EventID)
		if err != nil {
			render.Status(r, 500)
			l.Error("err to add tag to db: " + err.Error())
			return
		}

		render.Status(r, 200)
	}
}
