package extensions

import (
	"calendar-api/internal/usercontext"
	"calendar-api/types"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type DeleteRequestBody struct {
	ExtensionID uint `json:"extensionID"`
}

type Deleter interface {
	DeleteExtension(email string, extensionID uint) error
}

func Delete(logger *slog.Logger, extensionDeleter Deleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.With(
			slog.String("op", "handlers.extensions.Delete"),
			slog.String("requestId", middleware.GetReqID(r.Context())),
		)

		var requestBody DeleteRequestBody
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

		if !hasUserExtension(user, requestBody.ExtensionID) {
			render.Status(r, 403)
			l.Debug("validation error")
			return
		}

		l.Info("deleting ExtensionData from db")
		err = extensionDeleter.DeleteExtension(user.Email, requestBody.ExtensionID)
		if err != nil {
			render.Status(r, 500)
			l.Error("err to delete ExtensionData from db: " + err.Error())
			return
		}

		render.Status(r, 200)
	}
}

func hasUserExtension(user *types.User, extensionID uint) bool {
	for _, extensionData := range user.ExtensionsData {
		if extensionData.Extension == extensionID {
			return true
		}
	}

	return false
}
