package tags

import (
	"log/slog"
)

type Storage interface {
	Adder
	Deleter
}

type Service struct {
	l       *slog.Logger
	storage Storage
}

func New(logger *slog.Logger, storage Storage) *Service {
	return &Service{
		l:       logger,
		storage: storage,
	}
}
