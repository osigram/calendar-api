package main

import (
	"calendar-api/lib/config"
	"calendar-api/lib/log"
	"golang.org/x/exp/slog"
	"io"
	"os"
)

func main() {
	// Config
	cfg := config.NewConfig()

	// Logger
	var logWriter io.WriteCloser
	if cfg.LogToConsole {
		logWriter = os.Stdout
	} else {
		var err error
		logWriter, err = os.Create(cfg.LogFilePath)
		if err != nil {
			panic("Error in opening log file")
		}
		defer logWriter.Close()
	}
	logger := NewLogger(cfg, io.Writer(logWriter))

	logger.Info("Starting application...")

}

func NewLogger(cfg config.Config, w io.Writer) *slog.Logger {
	var logger *slog.Logger
	if cfg.BuildMode == config.Prod {
		logger = log.NewProdLogger(w)
	} else if cfg.BuildMode == config.Dev {
		logger = log.NewDevLogger(w)
	} else {
		panic("Error to initialize logger in main")
	}

	return logger
}
