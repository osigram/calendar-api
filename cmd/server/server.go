package main

import (
	"calendar-api/internal/config"
	pkgextensions "calendar-api/internal/extensions"
	"calendar-api/internal/extensions/khnure"
	"calendar-api/internal/extensions/mapper"
	"calendar-api/internal/handlers"
	"calendar-api/internal/log"
	"calendar-api/internal/middlewares/authmock"
	"calendar-api/internal/services/events"
	"calendar-api/internal/services/extensions"
	"calendar-api/internal/services/tags"
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

	extensionMapper := mapper.NewExtensionMapper()
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
	extensionMapper pkgextensions.Getter,
	authMiddleware Middleware,
) http.Handler {
	cfg := handlers.NewConfiguration(logger, storage, extensionMapper)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/event", func(r chi.Router) {
		r.Use(authMiddleware)

		eventService := events.New(logger, storage, extensionMapper)

		r.Get("/byID", handlers.Get(cfg, eventService.GetByID))
		r.Get("/byDate", handlers.Get(cfg, eventService.GetByDate))
		r.Post("/", handlers.Post(cfg, eventService.Add))
		r.Put("/", handlers.Post(cfg, eventService.Update))
		r.Delete("/", handlers.Post(cfg, eventService.Delete))
	})

	r.Route("/tag", func(r chi.Router) {
		r.Use(authMiddleware)

		tagsService := tags.New(logger, storage)

		r.Post("/", handlers.Post(cfg, tagsService.Add))
		r.Delete("/", handlers.Post(cfg, tagsService.Delete))
	})

	r.Route("/extension", func(r chi.Router) {
		r.Use(authMiddleware)

		extensionsService := extensions.New(logger, storage, extensionMapper)

		r.Post("/", handlers.Post(cfg, extensionsService.InstallOrUpdate))
		r.Delete("/", handlers.Post(cfg, extensionsService.Delete))
	})

	return r
}

func mustNewLogger(cfg config.Config, w io.Writer) *slog.Logger {
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
