package extensions

import (
	pkgcontext "calendar-api/internal/pkg/context"
	"calendar-api/internal/pkg/errors"
	"context"
	"log/slog"
)

type InstallRequestBody struct {
	ExtensionID    uint   `json:"extensionID"`
	AdditionalData string `json:"additionalData"`
}

type Installer interface {
	InstallOrUpdateExtension(ctx context.Context, email string, extensionID uint, additionalData string) error
}

func (s *Service) InstallOrUpdate(ctx *pkgcontext.Context, requestBody InstallRequestBody) error {
	l := s.l.With(
		slog.String("op", "services.extensions.InstallOrUpdate"),
	)

	if requestBody.ExtensionID == 0 || requestBody.AdditionalData == "" {
		l.Debug(
			"validation error: extensionID is zero or additionalData is empty",
			slog.Uint64("extensionID", uint64(requestBody.ExtensionID)),
			slog.String("additionalData", requestBody.AdditionalData),
		)
		return errors.NewValidationError("extensionID is zero or additionalData is empty")
	}

	extension, err := s.extensionsGetter.Get(ctx, requestBody.ExtensionID)
	if err != nil {
		l.Debug("unable to get extension from extensionMapper", slog.String("err", err.Error()))
		return errors.NewNotFoundError("extension not found")
	}

	if !extension.ValidateAdditionalData(ctx, requestBody.AdditionalData) {
		l.Debug("validation error: unsuitable additional data", slog.String("additionalData", requestBody.AdditionalData))
		return errors.NewValidationError("additional data is invalid")
	}

	l.Info("adding ExtensionData to db")
	err = s.storage.InstallOrUpdateExtension(ctx, ctx.User.Email, requestBody.ExtensionID, requestBody.AdditionalData)
	if err != nil {
		l.Error("unable to add ExtensionData to db", slog.String("err", err.Error()))
		return errors.NewInternalError("unable to install extension")
	}

	return nil
}
