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

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/notification/internal/api"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/notification/internal/api/handler"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/notification/internal/config"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/notification/internal/database"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/notification/internal/repository/postgres"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/notification/internal/service"
)

func main() {
	cfg := config.Load()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, err := database.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("notification service: %v", err)
	}
	defer db.Close()

	store := postgres.New(db)
	notificationService := service.NewNotificationService(store)

	handlers := handler.New(notificationService)
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
		log.Printf("notification service listening on :%s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("notification service: %v", err)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("notification service shutdown: %v", err)
	}
}
