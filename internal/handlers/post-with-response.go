package handlers

import (
	"calendar-api/internal/core"
	"calendar-api/internal/pkg/context"
	helpers2 "calendar-api/internal/pkg/helpers"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

func PostWithResponse[T any, V any](app App, service serviceWithResponseFunc[T, V], authorizedOnly bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := app.Logger().With(
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

		var user *core.User
		if authorizedOnly {
			user, err = helpers2.GetUser(r.Context())
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				l.Debug(err.Error())
				return
			}
		}

		ctx := &context.Context{
			Context:   r.Context(),
			User:      user,
			UserAgent: r.UserAgent(),
		}

		result, err := service(ctx, requestBody)
		if err != nil {
			helpers2.ProcessServiceError(w, r, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, core.Response[any]{
			Ok:   true,
			Data: result,
		})
	}
}
