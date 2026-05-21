package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/CB-InsuranceStack/InsuranceStack/apps/pricing-engine/internal/features"
	"github.com/CB-InsuranceStack/InsuranceStack/apps/pricing-engine/internal/handlers"
	"github.com/CB-InsuranceStack/InsuranceStack/apps/pricing-engine/internal/middleware"
	"github.com/CB-InsuranceStack/InsuranceStack/apps/pricing-engine/internal/repository"
	"github.com/CB-InsuranceStack/InsuranceStack/apps/pricing-engine/internal/services"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func main() {
	// Initialize logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Configure log level from environment
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logger.Warnf("Invalid log level '%s', defaulting to info", logLevel)
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	logger.Info("Starting Pricing Engine service...")

	// Get configuration from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8003"
	}

	dataPath := os.Getenv("DATA_PATH")
	if dataPath == "" {
		// Default to relative path from the project root
		dataPath = filepath.Join("..", "..", "data", "seed")
	}

	cloudBeesAPIKey := os.Getenv("CLOUDBEES_FM_API_KEY")
	if cloudBeesAPIKey == "" {
		logger.Warn("CLOUDBEES_FM_API_KEY not set, feature flags will use defaults")
		// Use a placeholder for development
		cloudBeesAPIKey = "dev-mode"
	}

	// Initialize CloudBees Feature Management
	flags, err := features.Initialize(cloudBeesAPIKey, logger)
	if err != nil {
		logger.WithError(err).Fatal("Failed to initialize feature management")
	}
	defer features.Shutdown()

	// Log feature flag status
	logger.WithFields(logrus.Fields{
		"dynamicRates": flags.IsDynamicRatesEnabled(),
	}).Info("Feature flags initialized")

	// Initialize repository
	repo, err := repository.NewRepository(dataPath, logger)
	if err != nil {
		logger.WithError(err).Fatal("Failed to initialize repository")
	}

	// Initialize services
	pricingService := services.NewPricingService(repo, flags, logger)

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler("pricing-engine")
	pricingHandler := handlers.NewPricingHandler(pricingService, logger)

	// Setup router
	router := mux.NewRouter()

	// Apply global middleware
	router.Use(middleware.LoggingMiddleware(logger))
	router.Use(middleware.AuthMiddleware(logger))

	// Setup CORS
	corsHandler := middleware.NewCORS()

	// Register routes
	router.Handle("/healthz", healthHandler).Methods("GET")
	router.HandleFunc("/quote", pricingHandler.GetQuote).Methods("POST")
	router.HandleFunc("/rates", pricingHandler.GetRates).Methods("GET")

	// Wrap router with CORS
	handler := corsHandler.Handler(router)

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Infof("Server listening on port %s", port)
		logger.Info("API Endpoints:")
		logger.Info("  GET  /healthz - Health check")
		logger.Info("  POST /quote - Calculate insurance quote")
		logger.Info("  GET  /rates - Get current base rates")
		logger.Info("")
		logger.Info("Feature Flags:")
		logger.Infof("  pricing.dynamicRates: %v (enables real-time rate adjustments)", flags.IsDynamicRatesEnabled())

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Fatal("Server failed to start")
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.WithError(err).Error("Server forced to shutdown")
	}

	logger.Info("Server stopped gracefully")
}
