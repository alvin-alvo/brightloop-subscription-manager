package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/internal/api"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/internal/api/handler"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/internal/api/middleware"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/internal/auth"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/internal/config"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/internal/database"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/internal/repository/postgres"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/internal/service"
)

func main() {
	if err := run(); err != nil {
		slog.Error("application stopped", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	if err := cfg.Validate(); err != nil {
		return err
	}

	db, err := database.Open(context.Background(), cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	users := postgres.NewUserRepository(db)
	tokens := auth.NewTokenManager(cfg.JWTSecret, cfg.JWTTTL)
	authService := service.NewAuthService(users, tokens)

	router := api.NewRouter(api.Dependencies{
		Health: handler.NewHealthHandler(db),
		Auth:   handler.NewAuthHandler(authService),
		AuthMW: middleware.NewAuthenticator(tokens),
	})

	server := &http.Server{
		Addr:              ":" + cfg.HTTPPort,
		Handler:           middleware.Recover(middleware.Logging(router)),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		slog.Info("http server listening", "addr", server.Addr, "env", cfg.Env)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	stopCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-errCh:
		return fmt.Errorf("serve http: %w", err)
	case <-stopCtx.Done():
	}

	slog.Info("shutting down http server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown http server: %w", err)
	}

	return nil
}
