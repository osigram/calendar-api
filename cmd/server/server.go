package main

import (
	"calendar-api/internal/config"
	"calendar-api/internal/extensions/extensionsmapping"
	"calendar-api/internal/extensions/khnure"
	"calendar-api/internal/handlers/events"
	"calendar-api/internal/handlers/extensions"
	"calendar-api/internal/handlers/tags"
	"calendar-api/internal/log"
	"calendar-api/internal/middlewares/authmock"
	"calendar-api/internal/storage"
	"calendar-api/internal/storage/gormstorage"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type Middleware = func(http.Handler) http.Handler

func main() {
	// Config
	cfg := config.NewConfig()

	// Logger
	logWriter := mustNewLogWriter(cfg.EnableConsoleLogging, cfg.LogFilePath)
	defer logWriter.Close()

	logger := mustNewLogger(cfg, io.Writer(logWriter))
	logger.Info("Starting application...")

	// Source initialisation
	storage, err := gormstorage.NewStorage(cfg.ConnectionString)
	if err != nil {
		panic(err)
	}

	extensionMapper := extensionsmapping.NewExtensionMapper()
	extensionMapper.RegisterExtension(1, khnure.NewTimeTableExtension())

	// auth
	authMiddleware := authmock.MockAuthMiddleware(logger, cfg, storage)

	// Router
	router := NewRouter(logger, storage, extensionMapper, authMiddleware)

	logger.Info("Starting server...", slog.String("url", cfg.URL))
	err = http.ListenAndServe(cfg.URL, router)
	if err != nil {
		logger.Error(fmt.Sprint(err))
	}
}

func NewRouter(logger *slog.Logger,
	storage storage.Storage,
	extensionMapper *extensionsmapping.ExtensionMapper,
	authMiddleware Middleware,
) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/event", func(r chi.Router) {
		r.Use(authMiddleware)

		r.Get("/byID", events.GetByID(logger, storage, extensionMapper))
		r.Get("/byDate", events.GetByDate(logger, storage, extensionMapper))
		r.Post("/", events.Add(logger, storage))
		r.Put("/", events.Update(logger, storage))
		r.Delete("/", events.Delete(logger, storage))
	})

	r.Route("/tag", func(r chi.Router) {
		r.Use(authMiddleware)

		r.Post("/", tags.Add(logger, storage))
		r.Delete("/", tags.Delete(logger, storage))
	})

	r.Route("/extension", func(r chi.Router) {
		r.Use(authMiddleware)

		r.Post("/", extensions.InstallOrUpdate(logger, storage, extensionMapper))
		r.Delete("/", extensions.Delete(logger, storage))
	})

	return r
}

func mustNewLogger(cfg config.Config, w io.Writer) *slog.Logger {
	var logger *slog.Logger
	switch cfg.BuildMode {
	case config.Prod:
		logger = log.NewProdLogger(w)
	case config.Dev:
		logger = log.NewDevLogger(w)
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
