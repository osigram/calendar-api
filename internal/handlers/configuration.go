package handlers

import (
	"calendar-api/internal/extensions"
	"calendar-api/internal/storage"
	"github.com/gorilla/schema"
	"log/slog"
)

type Configuration struct {
	L                *slog.Logger
	Storage          storage.Storage
	ExtensionsMapper extensions.Getter
	decoder          *schema.Decoder
}

func NewConfiguration(logger *slog.Logger, storage storage.Storage, extensionsMapper extensions.Getter) *Configuration {
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	decoder.ZeroEmpty(true)

	return &Configuration{
		L:                logger,
		Storage:          storage,
		ExtensionsMapper: extensionsMapper,
		decoder:          decoder,
	}
}
