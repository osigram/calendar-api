package server

import (
	"calendar-api/internal/handlers"
	"calendar-api/internal/services/events"
	"calendar-api/internal/services/extensions"
	"calendar-api/internal/services/tags"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"time"
)

func NewRouter(app *App) http.Handler {
	var (
		logger          = app.L
		storage         = app.Storage
		extensionMapper = app.ExtensionsMapper
		authMiddleware  = app.authMiddleware
		r               = chi.NewRouter()
	)

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/event", func(r chi.Router) {
		r.Use(authMiddleware)

		eventService := events.New(logger, storage, extensionMapper)

		r.Get("/byID", handlers.Get(app, eventService.GetByID))
		r.Get("/byDate", handlers.Get(app, eventService.GetByDate))
		r.Post("/", handlers.Post(app, eventService.Add))
		r.Put("/", handlers.Post(app, eventService.Update))
		r.Delete("/", handlers.Post(app, eventService.Delete))
	})

	r.Route("/tag", func(r chi.Router) {
		r.Use(authMiddleware)

		tagsService := tags.New(logger, storage)

		r.Post("/", handlers.Post(app, tagsService.Add))
		r.Delete("/", handlers.Post(app, tagsService.Delete))
	})

	r.Route("/extension", func(r chi.Router) {
		r.Use(authMiddleware)

		extensionsService := extensions.New(logger, storage, extensionMapper)

		r.Post("/", handlers.Post(app, extensionsService.InstallOrUpdate))
		r.Delete("/", handlers.Post(app, extensionsService.Delete))
	})

	return r
}
