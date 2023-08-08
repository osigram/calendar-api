package authmock

import (
	"calendar-api/lib/config"
	"calendar-api/middlewares"
	"calendar-api/types"
	"context"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/exp/slog"
	"net/http"
)

func MockAuthMiddleware(logger *slog.Logger, cfg config.Config, userGetSetter middlewares.UserGetSetter) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l := logger.With(
				slog.String("op", "middlewares.MockAuthMiddleware"),
				slog.String("requestId", middleware.GetReqID(r.Context())),
			)

			l.Debug("Using auth mock")

			var user *types.User
			var err error
			if user, err = userGetSetter.GetUser("user@example.com"); err != nil {
				user = &types.User{
					Email:          "user@example.com",
					Name:           "Example",
					PicturePath:    "https://instagram.com/favicon.ico",
					ExtensionsUsed: nil,
				}
				err = userGetSetter.AddUser(user)
				if err != nil {
					l.Error("could not add user")
					panic(err)
				}
			}

			ctx := context.WithValue(r.Context(), "user", user)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
