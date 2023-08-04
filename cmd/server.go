package main

import (
	"calendar-api/lib/config"
	"calendar-api/lib/log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/exp/slog"
	"io"
	"os"
	"time"
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

	// Source initialisation
	// TODO: source

	// Router
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(middleware.Timeout(60 * time.Second))

	// r.Use(google_auth.GoogleAuthMiddleware(logger, cfg, storage))
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
