package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nibble/mock-fps/internal/config"
	"github.com/nibble/mock-fps/internal/handlers"
	"github.com/nibble/mock-fps/internal/jsonapi"
	"github.com/nibble/mock-fps/internal/lifecycle"
	"github.com/nibble/mock-fps/internal/store"
	"github.com/nibble/mock-fps/internal/webhook"
)

func main() {
	cfg := config.Load()

	memStore := store.NewMemoryStore()

	dispatcher := webhook.NewDispatcher(memStore, cfg.WebhookBufferSize, cfg.WebhookWorkers)

	engine := lifecycle.NewEngine(cfg.LifecycleStepDelayMs, func(resourceType, resourceID, newStatus string) {
		dispatcher.Notify(resourceType, resourceID, newStatus)
	})

	mux := http.NewServeMux()
	handlers.RegisterRoutes(mux, memStore, engine)

	// Apply middleware chain: recovery -> logging -> content-type -> routes
	handler := handlers.Recovery(handlers.Logging(jsonapi.EnforceContentType(mux)))

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		log.Println("shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		dispatcher.Close()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("shutdown error: %v", err)
		}
	}()

	log.Printf("mock-fps server starting on :%s", cfg.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
	log.Println("server stopped")
}
