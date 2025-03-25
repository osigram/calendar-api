package log

import (
	"calendar-api/internal/config"
	"io"
	"log/slog"
	"os"
)

func MustNewLogger(cfg *config.Config, w io.Writer) *slog.Logger {
	var logger *slog.Logger
	switch cfg.BuildMode {
	case config.Prod:
		logger = MustNewProdLogger(w)
	case config.Dev:
		logger = MustNewDevLogger(w)
	default:
		panic("Error to initialize logger in main")
	}

	return logger
}

func MustNewLogWriter(enableConsoleLogging bool, filepath string) (logWriter io.WriteCloser) {
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
