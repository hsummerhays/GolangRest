package main

import (
	"context"
	"golangrest/internal/config"
	"golangrest/internal/handlers"
	"golangrest/internal/repository"
	"golangrest/internal/services"
	"golangrest/pkg/middleware"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"golang.org/x/time/rate"
	
	_ "golangrest/docs"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// @title           Golang REST API
// @version         1.0
// @description     A simple Golang REST application featuring products and health check.
// @host            localhost:8080
// @BasePath        /
func main() {
	// Initialize slog
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize dependencies
	productRepo, err := repository.NewSQLiteProductRepository(cfg.DBPath)
	if err != nil {
		slog.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	// Initialize router
	r := chi.NewRouter()

	// Global middleware
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	rateLimiter := middleware.NewIPRateLimiter(rate.Limit(cfg.RateLimit), cfg.RateLimitBurst)
	r.Use(middleware.NewRateLimitMiddleware(rateLimiter))

	// Routes
	r.Get("/health", handlers.HealthHandler)
	r.Get("/products", productHandler.GetProducts)
	r.Post("/products", productHandler.CreateProduct)

	// Swagger documentation route
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	// Setup HTTP server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	// Channel to listen for errors coming from the listener.
	serverErrors := make(chan error, 1)

	// Start server in a goroutine
	go func() {
		slog.Info("Server listening", "port", cfg.Port)
		serverErrors <- srv.ListenAndServe()
	}()

	// Channel to listen for an interrupt or terminate signal from the OS.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		slog.Error("Error starting server", "error", err)

	case sig := <-shutdown:
		slog.Info("Starting shutdown", "signal", sig)

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Asking listener to shutdown and shed load.
		if err := srv.Shutdown(ctx); err != nil {
			slog.Error("Graceful shutdown did not complete in time", "error", err)
			if err := srv.Close(); err != nil {
				slog.Error("Error killing server", "error", err)
			}
		}
	}
}
