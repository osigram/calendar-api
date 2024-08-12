package handlers

import (
	"calendar-api/internal/core"
	"calendar-api/internal/helpers"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

func PostWithResponse[T any, V any](cfg *Configuration, service serviceWithResponseFunc[T, V]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := cfg.L.With(
			slog.String("op", "internal.handlers.PostWithResponse"),
			slog.String("requestId", middleware.GetReqID(r.Context())),
		)

		var requestBody T
		err := render.DecodeJSON(r.Body, &requestBody)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			l.Debug("unable to decode request body", slog.String("err", err.Error()))
			return
		}

		user, err := helpers.GetUser(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			l.Debug(err.Error())
			return
		}

		result, err := service(user, requestBody)
		if err != nil {
			helpers.ProcessServiceError(w, r, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, core.Response[any]{
			Ok:   true,
			Data: result,
		})
	}
}
