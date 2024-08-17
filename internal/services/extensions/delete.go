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

type Getter interface {
	GetExtensionData(userEmail string, extensionID uint) (*core.ExtensionData, error)
}

type Deleter interface {
	DeleteExtension(email string, extensionID uint) error
}

func (s *Service) Delete(ctx *context.Context, requestBody DeleteRequestBody) error {
	l := s.l.With(
		slog.String("op", "services.extensions.Delete"),
	)

	if requestBody.ExtensionID == 0 {
		l.Debug("validation error: extensionID is zero", slog.Uint64("extensionID", uint64(requestBody.ExtensionID)))
		return errors.NewValidationError("extensionID is zero")
	}

	if _, err := s.storage.GetExtensionData(ctx.User.Email, requestBody.ExtensionID); err != nil {
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
