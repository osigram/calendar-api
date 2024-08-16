package server

import (
	"calendar-api/internal/config"
	"calendar-api/internal/extensions"
	"calendar-api/internal/extensions/khnure"
	"calendar-api/internal/extensions/mapper"
	"calendar-api/internal/middlewares/authmock"
	"calendar-api/internal/router"
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

type App struct {
	L                *slog.Logger
	Storage          storage.Storage
	ExtensionsMapper extensions.Getter
	Decoder          *schema.Decoder
	router           http.Handler
	logWriter        io.WriteCloser
	cfg              *config.Config
}

func (a *App) Run() {
	defer a.logWriter.Close()

	go func() {
		a.L.Info("Starting server...", slog.String("url", a.cfg.URL))
		if err := http.ListenAndServe(a.cfg.URL, a.router); err != nil {
			a.L.Error(fmt.Sprint(err))
		}
	}()

	// Graceful shutdown
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	<-exit
}

func NewApp() *App {
	cfg := config.MustNewConfig()

	logWriter := mustNewLogWriter(cfg.EnableConsoleLogging, cfg.LogFilePath) // must be closed
	logger := mustNewLogger(cfg, io.Writer(logWriter))
	logger.Info("Starting application...")

	storage, err := gormstorage.NewStorage(cfg.ConnectionString)
	if err != nil {
		panic(err)
	}

	extensionsMapper := mapper.NewExtensionMapper()
	extensionsMapper.RegisterExtension(1, khnure.NewTimeTableExtension())

	authMiddleware := authmock.MockAuthMiddleware(logger, cfg, storage)

	r := router.NewRouter(logger, storage, extensionsMapper, authMiddleware)

	// Query params decoder
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	decoder.ZeroEmpty(true)

	return &App{
		L:                logger,
		Storage:          storage,
		ExtensionsMapper: extensionsMapper,
		Decoder:          decoder,
		router:           r,
		logWriter:        logWriter,
	}
}
