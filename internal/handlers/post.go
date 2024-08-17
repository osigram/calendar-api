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

type serviceFunc[T any] func(*context.Context, T) error

func Post[T any](app App, service serviceFunc[T]) http.HandlerFunc {
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

		user, err := helpers2.GetUser(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			l.Debug(err.Error())
			return
		}

		ctx := &context.Context{
			Context:   r.Context(),
			User:      user,
			UserAgent: r.UserAgent(),
		}

		err = service(ctx, requestBody)
		if err != nil {
			helpers2.ProcessServiceError(w, r, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, core.Response[any]{
			Ok: true,
		})
	}
}
