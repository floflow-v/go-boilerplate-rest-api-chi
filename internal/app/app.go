package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"

	"go-rest-api-chi-example/internal/config"
	"go-rest-api-chi-example/internal/database"
	"go-rest-api-chi-example/internal/logger"
	"go-rest-api-chi-example/migrations"
)

type Application struct {
	apiServer *http.Server
	db        *database.Database
	logger    zerolog.Logger
}

func New() (*Application, error) {
	newConfig, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	newLogger, err := logger.NewLogger(&newConfig)
	if err != nil {
		return nil, err
	}

	db, err := database.Init(newConfig, newLogger)
	if err != nil {
		return nil, err
	}

	if err := migrations.Run(db.DB); err != nil {
		return nil, err
	}

	// API
	apiServer := buildAPI(newConfig, newLogger, db)

	return &Application{
		apiServer: apiServer,
		db:        db,
		logger:    newLogger,
	}, nil
}

func (a *Application) Run() error {
	go func() {
		a.logger.Info().Msgf("API server started and listening at http://%s", a.apiServer.Addr)
		if err := a.apiServer.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			a.logger.Fatal().Err(err).Msg("API server failed")
		}
	}()

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	a.logger.Info().Msg("Shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := a.apiServer.Shutdown(shutdownCtx); err != nil {
		a.logger.Error().Err(err).Msg("API shutdown error")
	}

	if err := a.db.Close(); err != nil {
		a.logger.Error().Err(err).Msg("Database close error")
	}

	a.logger.Info().Msg("Shutdown complete")
	return nil
}
