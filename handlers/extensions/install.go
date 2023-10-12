package extensions

import (
	"calendar-api/internal/extensionsmapping"
	"calendar-api/internal/usercontext"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type InstallRequestBody struct {
	ExtensionID    uint   `json:"extensionID"`
	AdditionalData string `json:"additionalData"`
}

type Installer interface {
	InstallOrUpdateExtension(email string, extensionID uint, additionalData string) error
}

type ExtensionGetter interface {
	Get(id uint) (extensionsmapping.Extension, error)
}

func InstallOrUpdate(logger *slog.Logger, extensionInstaller Installer, extensionMapper ExtensionGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.With(
			slog.String("op", "handlers.extensions.InstallOrUpdate"),
			slog.String("requestId", middleware.GetReqID(r.Context())),
		)

		var requestBody InstallRequestBody
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

		if requestBody.ExtensionID == 0 || requestBody.AdditionalData == "" {
			render.Status(r, 403)
			l.Debug("validation error")
			return
		}

		extension, err := extensionMapper.Get(requestBody.ExtensionID)
		if err != nil {
			render.Status(r, 404)
			l.Debug("err to get extension from extensionMapper: " + err.Error())
			return
		}
		if !extension.ValidateAdditionalData(requestBody.AdditionalData) {
			render.Status(r, 403)
			l.Debug("additional data validation error")
			return
		}

		l.Info("adding ExtensionData to db")
		err = extensionInstaller.InstallOrUpdateExtension(user.Email, requestBody.ExtensionID, requestBody.AdditionalData)
		if err != nil {
			render.Status(r, 500)
			l.Error("err to add ExtensionData to db: " + err.Error())
			return
		}

		render.Status(r, 200)
	}
}
