package handlers

import (
	"github.com/gorilla/schema"
	"log/slog"
)

type App interface {
	Logger() *slog.Logger
	Decoder() *schema.Decoder
}
