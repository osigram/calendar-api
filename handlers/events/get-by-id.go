package events

import (
	"calendar-api/internal/usercontext"
	"calendar-api/types"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
)

type GetEventByIDRequest struct {
	ID     uint `json:"id"`
	Source uint `json:"source,omitempty"`
}

type ByIDGetter interface {
	GetEventByID(uint) (*types.Event, error)
}

func GetByID(logger *slog.Logger, eventGetter ByIDGetter, extensionMapper ExtensionGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.With(
			slog.String("op", "handlers.events.GetById"),
			slog.String("requestId", middleware.GetReqID(r.Context())),
		)

		var requestBody GetEventByIDRequest
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
		// TODO: if id == 0 return error

		var eg ByIDGetter
		if requestBody.Source == 0 {
			eg = eventGetter
		} else {
			if !hasUserExtension(user, requestBody.Source) {
				render.Status(r, 403)
				return
			}

			extension, err := extensionMapper.Get(requestBody.Source)
			if err != nil {
				render.Status(r, 501)
				l.Error("ExtensionData " + strconv.Itoa(int(requestBody.Source)) + "is not implemented")
				return
			}

			eg = extension
		}

		l.Info("getting event from db")
		event, err := eg.GetEventByID(requestBody.ID)
		if err != nil {
			render.Status(r, 404)
			l.Debug("err to get event from db")
			return
		}

		if user.Email != event.UserEmail {
			render.Status(r, 401)
			l.Debug("user is not the owner of this event")
			return
		}

		render.Status(r, 200)
		render.JSON(w, r, event)
	}
}

func hasUserExtension(user *types.User, id uint) bool {
	for _, ed := range user.ExtensionsData {
		if id == ed.ID {
			return true
		}
	}

	return false
}
