package extensions

import (
	"calendar-api/pkg/extensions"
	"log/slog"
)

type Storage interface {
	Installer
	Deleter
	Getter
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
