package main

import (
	"calendar-api/handlers/events"
	"calendar-api/internal/extensions"
	"calendar-api/lib/config"
	"calendar-api/lib/log"
	"calendar-api/middlewares/authmock"
	"calendar-api/storage/postgres"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/exp/slog"
	"io"
	"net/http"
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
	storage, err := postgres.NewStorage(cfg.ConnectionString)
	if err != nil {
		panic(err)
	}

	extensionMapper := extensions.NewExtensionMapper()

	// Router
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/event", func(r chi.Router) {
		r.Use(authmock.MockAuthMiddleware(logger, cfg, storage))

		r.Get("/byId", events.GetById(logger, storage, extensionMapper))
		r.Get("/byDate", events.GetByDate(logger, storage, extensionMapper))
		r.Post("/", events.Add(logger, storage))
		r.Put("/", events.Update(logger, storage))
		r.Delete("/", events.Delete(logger, storage))
	})

	http.ListenAndServe(":8080", r)
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
