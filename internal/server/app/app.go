package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/grnsv/GophKeeper/internal/api"
	"github.com/grnsv/GophKeeper/internal/server/config"
	"github.com/grnsv/GophKeeper/internal/server/handlers"
	"github.com/grnsv/GophKeeper/internal/server/interfaces"
	"github.com/grnsv/GophKeeper/internal/server/service"
	"github.com/grnsv/GophKeeper/internal/server/storage"
)

type application struct {
	Config     *config.Config
	Storage    interfaces.Storage
	JWTService interfaces.JWTService
	Service    interfaces.Service
	Server     *http.Server
}

func New(ctx context.Context, buildVersion, buildDate string) (app *application, err error) {
	app = &application{}
	if app.Config, err = config.Parse(); err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}
	if app.Storage, err = storage.New(ctx, app.Config.DatabaseDSN, app.Config.MigrationsPath); err != nil {
		return nil, fmt.Errorf("storage: %w", err)
	}
	app.JWTService = service.NewJWTService(app.Config.JWTSecret)
	if app.Service, err = service.New(app.Storage, app.JWTService, buildVersion, buildDate); err != nil {
		return nil, fmt.Errorf("service: %w", err)
	}
	server, err := api.NewServer(
		handlers.NewHandler(app.Service),
		handlers.NewSecurityHandler(app.JWTService),
		api.WithErrorHandler(handlers.ErrorHandler),
	)
	if err != nil {
		return nil, fmt.Errorf("server: %w", err)
	}
	app.Server = &http.Server{
		Addr:         app.Config.RunAddress,
		Handler:      server,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return
}

func (app *application) Run() {
	log.Printf("Starting server at %s", app.Server.Addr)
	if err := app.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}
}

func (app *application) Shutdown(ctx context.Context) error {
	if err := app.Server.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown server: %w", err)
	}
	if err := app.Storage.Close(); err != nil {
		return fmt.Errorf("close storage: %w", err)
	}

	return nil
}
