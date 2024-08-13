package tags

import (
	"calendar-api/internal/extensions"
	"log/slog"
)

type Storage interface {
	Adder
	Deleter
}

type Service struct {
	l                *slog.Logger
	storage          Storage
	extensionsGetter extensions.Getter
}

func New(logger *slog.Logger, storage Storage) *Service {
	return &Service{
		l:       logger,
		storage: storage,
	}
}
