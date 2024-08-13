package extensions

import (
	"calendar-api/internal/extensions"
	"log/slog"
)

type Storage interface {
	Installer
	Deleter
}

type Service struct {
	l                *slog.Logger
	storage          Storage
	extensionsGetter extensions.Getter
}

func New(logger *slog.Logger, storage Storage, extensionsGetter extensions.Getter) *Service {
	return &Service{
		l:                logger,
		storage:          storage,
		extensionsGetter: extensionsGetter,
	}
}
