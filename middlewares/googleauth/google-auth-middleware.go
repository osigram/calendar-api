package googleauth

import (
	"calendar-api/internal/config"
	"calendar-api/middlewares"
	"calendar-api/types"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"google.golang.org/api/idtoken"
	"log/slog"
	"net/http"
	"strings"
)

func GoogleAuthMiddleware(logger *slog.Logger, cfg config.Config, userGetSetter middlewares.UserGetSetter) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l := logger.With(
				slog.String("op", "middlewares.GoogleAuthMiddleware"),
				slog.String("requestId", middleware.GetReqID(r.Context())),
			)

			auth := r.Header.Get("Authorization")
			authList := strings.Split(auth, " ") // ["Bearer", "{token}"]
			if len(authList) != 2 {
				l.Debug("unable to parse Authorization header", slog.String("Authorization", auth))

				render.Status(r, 401)
				return
			}
			token := authList[1]

			payload, err := idtoken.Validate(r.Context(), token, cfg.GoogleClientId)
			if err != nil || payload == nil || payload.Claims == nil {
				l.Debug(
					"unable to validate google token",
					slog.String("token", token),
					slog.String("err", err.Error()),
				)

				render.Status(r, 401)
				return
			}

			// If user already registered
			var email string
			if emailAny, ok := payload.Claims["email"]; !ok {
				claimsJson, _ := json.Marshal(payload.Claims)
				l.Debug(
					"unable to get email from payload",
					slog.String("claims", string(claimsJson)),
				)

				render.Status(r, 401)
				return
			} else if email, ok = emailAny.(string); !ok {
				claimsJson, _ := json.Marshal(payload.Claims)
				l.Debug(
					"unable to get email from payload",
					slog.String("claims", string(claimsJson)),
				)

				render.Status(r, 401)
				return
			}

			var user *types.User
			if user, err = userGetSetter.GetUser(email); err != nil {
				if user, err = RegisterUser(userGetSetter, payload.Claims); err != nil {
					claimsJson, _ := json.Marshal(payload.Claims)
					l.Debug(
						"unable to register user",
						slog.String("claims", string(claimsJson)),
						slog.String("err", err.Error()),
					)

					render.Status(r, 401)
					return
				}
			}

			ctx := context.WithValue(r.Context(), "user", &user)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RegisterUser(userGetSetter middlewares.UserGetSetter, claims map[string]any) (*types.User, error) {
	var email string
	var err error
	if email, err = getClaim(claims, "email"); err != nil {
		return nil, err
	}

	var name string
	if name, err = getClaim(claims, "name"); err != nil {
		return nil, err
	}

	var picturePath string
	if picturePath, err = getClaim(claims, "picture"); err != nil {
		return nil, err
	}

	user := types.User{
		Email:          email,
		Name:           name,
		PicturePath:    picturePath,
		ExtensionsData: nil,
	}

	if err = userGetSetter.AddUser(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func getClaim(claims map[string]any, key string) (string, error) {
	var result string
	if resultAny, ok := claims[key]; !ok {
		err := errors.New("unable to get " + key + " from payload")

		return "", err
	} else if result, ok = resultAny.(string); !ok {
		err := errors.New("unable to get " + key + " from payload")

		return "", err
	}

	return result, nil
}
