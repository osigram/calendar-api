package server

import (
	"calendar-api/internal/config"
	"calendar-api/internal/extensions"
	"calendar-api/internal/extensions/khnure"
	"calendar-api/internal/extensions/mapper"
	"calendar-api/internal/log"
	"calendar-api/internal/middlewares/authmock"
	"calendar-api/internal/storage"
	"calendar-api/internal/storage/gormstorage"
	"fmt"
	"github.com/gorilla/schema"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Middleware = func(http.Handler) http.Handler

type App struct {
	L                *slog.Logger
	Storage          storage.Storage
	ExtensionsMapper extensions.Getter
	decoder          *schema.Decoder
	authMiddleware   Middleware
	logWriter        io.WriteCloser
	cfg              *config.Config
}

func (a *App) Logger() *slog.Logger {
	return a.L
}

func (a *App) Decoder() *schema.Decoder {
	return a.decoder
}

func (a *App) Run() {
	defer a.logWriter.Close()

	r := NewRouter(a)

	go func() {
		a.L.Info("Starting server...", slog.String("url", a.cfg.URL))
		if err := http.ListenAndServe(a.cfg.URL, r); err != nil {
			a.L.Error(fmt.Sprint(err))
		}
	}()

	// Graceful shutdown
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	<-exit
}

func NewApp(cfg *config.Config) *App {
	logWriter := log.MustNewLogWriter(cfg.EnableConsoleLogging, cfg.LogFilePath) // must be closed
	logger := log.MustNewLogger(cfg, logWriter)
	logger.Info("Starting application...")

	s, err := gormstorage.NewStorage(cfg.ConnectionString)
	if err != nil {
		panic(err)
	}

	extensionsMapper := mapper.NewExtensionMapper()
	extensionsMapper.RegisterExtension(1, khnure.NewTimeTableExtension())

	authMiddleware := authmock.MockAuthMiddleware(logger, cfg, s)

	// Query params decoder
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	decoder.ZeroEmpty(true)

	return &App{
		L:                logger,
		Storage:          s,
		ExtensionsMapper: extensionsMapper,
		decoder:          decoder,
		authMiddleware:   authMiddleware,
		logWriter:        logWriter,
		cfg:              cfg,
	}
}
