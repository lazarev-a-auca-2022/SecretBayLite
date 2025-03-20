// Package main is the entry point for the SecretBay VPN configuration server.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"secretbay/internal/api"
	"secretbay/internal/config"
)

func main() {
	// Initialize logger
	logger := log.New(os.Stdout, "SecretBay: ", log.LstdFlags|log.Lshortfile)

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize API server
	server := api.NewServer(cfg, logger)

	// Create HTTP server
	httpServer := &http.Server{
		Addr:         cfg.ServerAddress,
		Handler:      server.Router(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 120 * time.Second, // As per requirements: 120 seconds max processing time
	}

	// Start server in a goroutine
	go func() {
		// SSL certificate paths
		certFile := filepath.Join("/etc/letsencrypt/live/localhost/fullchain.pem")
		keyFile := filepath.Join("/etc/letsencrypt/live/localhost/privkey.pem")

		// Try to start with HTTPS
		if _, err := os.Stat(certFile); err == nil {
			if _, err := os.Stat(keyFile); err == nil {
				logger.Printf("Starting HTTPS server on %s", cfg.ServerAddress)
				if err := httpServer.ListenAndServeTLS(certFile, keyFile); err != nil && err != http.ErrServerClosed {
					logger.Printf("HTTPS server failed: %v, falling back to HTTP", err)
				}
				return
			}
		}

		// Fall back to HTTP if certificates are not available
		logger.Printf("Starting HTTP server on %s", cfg.ServerAddress)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Printf("Server forced to shutdown: %v", err)
	}

	logger.Println("Server gracefully stopped")
}
