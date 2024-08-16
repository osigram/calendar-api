package server

import (
	"calendar-api/internal/config"
	"calendar-api/internal/log"
	"io"
	"log/slog"
	"os"
)

func mustNewLogger(cfg *config.Config, w io.Writer) *slog.Logger {
	var logger *slog.Logger
	switch cfg.BuildMode {
	case config.Prod:
		logger = log.MustNewProdLogger(w)
	case config.Dev:
		logger = log.MustNewDevLogger(w)
	default:
		panic("Error to initialize logger in main")
	}

	return logger
}

func mustNewLogWriter(enableConsoleLogging bool, filepath string) (logWriter io.WriteCloser) {
	if enableConsoleLogging {
		logWriter = os.Stdout
	} else {
		var err error
		logWriter, err = os.Create(filepath)
		if err != nil {
			panic("Error in opening log file")
		}
	}

	return logWriter
}
