package server

import (
	"calendar-api/internal/config"
	"calendar-api/internal/extensions/khnure"
	"calendar-api/internal/extensions/mapper"
	"calendar-api/internal/log"
	"calendar-api/internal/storage"
	"calendar-api/internal/storage/gormstorage"
	"calendar-api/pkg/extensions"
	"context"
	"github.com/go-chi/jwtauth/v5"
	"github.com/gorilla/schema"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"golang.org/x/sync/errgroup"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	logger           *slog.Logger
	storage          storage.Storage
	extensionsMapper extensions.Getter
	decoder          *schema.Decoder
	accessTokenAuth  *jwtauth.JWTAuth
	refreshTokenAuth *jwtauth.JWTAuth
	logWriter        io.WriteCloser
	cfg              *config.Config
}

func (a *App) Logger() *slog.Logger {
	return a.logger
}

func (a *App) Decoder() *schema.Decoder {
	return a.decoder
}

func (a *App) Run() {
	defer a.logWriter.Close()

	r := NewRouter(a)

	// Graceful shutdown context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Init server
	server := &http.Server{
		Addr:    a.cfg.URL,
		Handler: r,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}

	a.logger.Info("running server", slog.String("url", a.cfg.URL))
	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return server.ListenAndServe()
	})

	g.Go(func() error {
		<-gCtx.Done()
		return server.Shutdown(context.Background())
	})

	if err := g.Wait(); err != nil {
		a.logger.Error("server shutdown", slog.String("err", err.Error()))
	}
}

func NewApp(cfg *config.Config) *App {
	logWriter := log.MustNewLogWriter(cfg.EnableConsoleLogging, cfg.LogFilePath) // must be closed
	logger := log.MustNewLogger(cfg, logWriter)

	logger.Info("starting application")

	s, err := gormstorage.NewStorage(cfg.ConnectionString)
	if err != nil {
		panic(err)
	}

	extensionsMapper := mapper.NewExtensionMapper()
	extensionsMapper.RegisterExtension(1, khnure.NewTimeTableExtension())

	accessTokenAuth := jwtauth.New("HS512", []byte(cfg.AuthSecret), nil, jwt.WithAcceptableSkew(time.Duration(cfg.AccessTokenExpInMinutes)*time.Minute))
	refreshTokenAuth := jwtauth.New("HS512", []byte(cfg.AuthSecret), nil, jwt.WithAcceptableSkew(time.Duration(cfg.RefreshTokenExpInDays)*24*time.Hour))

	// Query params decoder
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	decoder.ZeroEmpty(true)

	return &App{
		logger:           logger,
		storage:          s,
		extensionsMapper: extensionsMapper,
		decoder:          decoder,
		accessTokenAuth:  accessTokenAuth,
		refreshTokenAuth: refreshTokenAuth,
		logWriter:        logWriter,
		cfg:              cfg,
	}
}
