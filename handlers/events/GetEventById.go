package events

import (
	"calendar-api/internal/userContext"
	"calendar-api/types"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
	"net/http"
	"strconv"
)

type GetEventByIdRequest struct {
	Id     int64 `json:"id"`
	Source int64 `json:"source,omitempty"`
}

type ByIdGetter interface {
	GetEventById(int64) (*types.Event, error)
}

func GetById(logger *slog.Logger, eventGetter ByIdGetter, extensionMapper ExtensionGetter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := logger.With(
			slog.String("op", "handlers.events.GetById"),
			slog.String("requestId", middleware.GetReqID(r.Context())),
		)

		var requestBody GetEventByIdRequest
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

		var eg ByIdGetter
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
		event, err := eg.GetEventById(requestBody.Id)
		if err != nil {
			render.Status(r, 404)
			l.Debug("err to get event from db")
			return
		}

		render.Status(r, 200)
		render.JSON(w, r, event)
	})
}

func hasUserExtension(user *types.User, id int64) bool {
	for _, ed := range user.ExtensionsUsed {
		if id == ed.Id {
			return true
		}
	}

	return false
}
