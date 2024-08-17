package extensions

import (
	"calendar-api/internal/core"
	"calendar-api/internal/pkg/context"
	"calendar-api/internal/pkg/errors"
	"log/slog"
)

type DeleteRequestBody struct {
	ExtensionID uint `json:"extensionID"`
}

type Deleter interface {
	DeleteExtension(email string, extensionID uint) error
}

func (s *Service) Delete(ctx *context.Context, requestBody DeleteRequestBody) error {
	l := s.l.With(
		slog.String("op", "services.extensions.Delete"),
	)

	if !hasUserExtension(ctx.User, requestBody.ExtensionID) {
		l.Debug("validation error: user has no such extension")
		return errors.NewValidationError("user has no such extension")
	}

	l.Info("deleting ExtensionData from db")
	err := s.storage.DeleteExtension(ctx.User.Email, requestBody.ExtensionID)
	if err != nil {
		l.Error("unable to delete ExtensionData from db", slog.String("err", err.Error()))
		return errors.NewInternalError("unable to delete extension data from db")
	}

	return nil
}

// TODO: this is a wrong way to get extensionsData, because of the auth logic. Later, should be rewritten and moved to helpers
func hasUserExtension(user *core.User, extensionID uint) bool {
	for _, extensionData := range user.ExtensionsData {
		if extensionData.Extension == extensionID {
			return true
		}
	}

	return false
}
