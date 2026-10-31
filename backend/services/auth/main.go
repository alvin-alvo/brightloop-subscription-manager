package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/auth/internal/api"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/auth/internal/api/handler"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/auth/internal/auth"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/auth/internal/config"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/auth/internal/database"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/auth/internal/repository/postgres"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/auth/internal/service"
)

func main() {
	cfg := config.Load()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, err := database.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("auth service: %v", err)
	}
	defer db.Close()

	store := postgres.New(db)
	tokens := auth.NewTokenManager(cfg.EffectiveSecret(), cfg.JWTIssuer, cfg.AccessTokenTTL)
	authService := service.NewAuthService(store, store, tokens, cfg.PasswordResetTTL)

	handlers := handler.New(authService)
	router := api.NewRouter(cfg, handlers)

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		log.Printf("auth service listening on :%s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("auth service: %v", err)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("auth service shutdown: %v", err)
	}
}
