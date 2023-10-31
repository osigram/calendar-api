package main

import (
	"calendar-api/handlers/events"
	"calendar-api/handlers/extensions"
	"calendar-api/handlers/tags"
	"calendar-api/internal/config"
	"calendar-api/internal/extensionsmapping"
	"calendar-api/internal/khnure"
	"calendar-api/internal/log"
	"calendar-api/middlewares/authmock"
	"calendar-api/storage/gormstorage"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"io"
	"log/slog"
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
	storage, err := gormstorage.NewStorage(cfg.ConnectionString)
	if err != nil {
		panic(err)
	}

	extensionMapper := extensionsmapping.NewExtensionMapper()
	extensionMapper.RegisterExtension(1, khnure.NewTimeTableExtension())

	// Router
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/event", func(r chi.Router) {
		r.Use(authmock.MockAuthMiddleware(logger, cfg, storage))

		r.Get("/byID", events.GetByID(logger, storage, extensionMapper))
		r.Get("/byDate", events.GetByDate(logger, storage, extensionMapper))
		r.Post("/", events.Add(logger, storage))
		r.Put("/", events.Update(logger, storage))
		r.Delete("/", events.Delete(logger, storage))
	})

	r.Route("/tag", func(r chi.Router) {
		r.Use(authmock.MockAuthMiddleware(logger, cfg, storage))

		r.Post("/", tags.Add(logger, storage))
		r.Delete("/", tags.Delete(logger, storage))
	})

	r.Route("/extension", func(r chi.Router) {
		r.Use(authmock.MockAuthMiddleware(logger, cfg, storage))

		r.Post("/", extensions.InstallOrUpdate(logger, storage, extensionMapper))
		r.Delete("/", extensions.Delete(logger, storage))
	})

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		logger.Error(fmt.Sprint(err))
	}
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
