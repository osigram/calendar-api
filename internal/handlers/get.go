package handlers

import (
	"calendar-api/internal/context"
	"calendar-api/internal/core"
	"calendar-api/internal/helpers"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type serviceWithResponseFunc[T any, V any] func(*context.Context, T) (V, error)

func Get[T any, V any](cfg *Configuration, service serviceWithResponseFunc[T, V]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := cfg.L.With(
			slog.String("op", "internal.handlers.Get"),
			slog.String("requestId", middleware.GetReqID(r.Context())),
		)

		var requestBody T
		err := cfg.decoder.Decode(&requestBody, r.URL.Query())
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

		ctx := &context.Context{
			Context: r.Context(),
			User:    user,
		}

		result, err := service(ctx, requestBody)
		if err != nil {
			statusCode, message := helpers.ProcessServiceError(w, r, err)
			l.Debug("service returned an error", slog.Int("statusCode", statusCode), slog.String("message", message))
			return
		}

		w.WriteHeader(http.StatusOK)
		l.Debug("service returned a response")
		render.JSON(w, r, core.Response[any]{
			Ok:   true,
			Data: result,
		})
	}
}
