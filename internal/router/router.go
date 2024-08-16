package router

import (
	"calendar-api/internal/app/server"
	pkgextensions "calendar-api/internal/extensions"
	"calendar-api/internal/handlers"
	"calendar-api/internal/services/events"
	"calendar-api/internal/services/extensions"
	"calendar-api/internal/services/tags"
	"calendar-api/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"time"
)

type Middleware = func(http.Handler) http.Handler

func NewRouter(logger *slog.Logger,
	storage storage.Storage,
	extensionMapper pkgextensions.Getter,
	authMiddleware Middleware,
) http.Handler {
	cfg := server.NewApp(logger, storage, extensionMapper)

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
