package log

import (
	"io"
	"log/slog"
)

func MustNewDevLogger(w io.Writer) *slog.Logger {
	if w == nil {
		panic("error to initialize production logger")
	}

	return slog.New(
		slog.NewTextHandler(w, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
		),
	)
}

func MustNewProdLogger(w io.Writer) *slog.Logger {
	if w == nil {
		panic("error to initialize production logger")
	}

	return slog.New(
		slog.NewJSONHandler(w, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		},
		),
	)
}
