package cmd

import (
	"context"
	"copier/config"
	"copier/internal/logger"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

var serveRestCmd = &cobra.Command{
	Use:   "serve-rest",
	Short: "Serve the REST API",
	RunE:  serveRest,
}

func serveRest(cmd *cobra.Command, args []string) error {
	// Load configuration
	conf := config.GetConfig()
	
	// Initialize logger
	logger.SetupLogger(conf.ServiceName, string(conf.Mode))
	
	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + strconv.Itoa(conf.HttpPort),
		Handler:      setupRoutes(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	serverErrors := make(chan error, 1)
	go func() {
		slog.Info("Starting HTTP server", "port", conf.HttpPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrors <- fmt.Errorf("server error: %w", err)
		}
	}()

	// Wait for either server error or shutdown signal
	select {
	case err := <-serverErrors:
		slog.Error("Server error occurred", "error", err)
		return err
	case <-shutdown():
		slog.Info("Shutdown signal received, starting graceful shutdown")
		return gracefulShutdown(server)
	}
}

func setupRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	
	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"copier","timestamp":"` + time.Now().Format(time.RFC3339) + `"}`))
	})
	
	// Welcome endpoint
	mux.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"Welcome to Signal Copier API","version":"1.0.0"}`))
	})
	
	// API v1 endpoint
	mux.HandleFunc("/api/v1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"Welcome to Signal Copier API v1","version":"1.0.0"}`))
	})
	
	// Catch-all for 404
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"Not Found","message":"The requested resource was not found"}`))
	})
	
	return mux
}

func shutdown() chan os.Signal {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown,
		os.Interrupt,    // Ctrl+C
		syscall.SIGTERM, // Termination signal
		syscall.SIGINT,  // Interrupt signal
		syscall.SIGQUIT, // Quit signal
	)
	return shutdown
}

func gracefulShutdown(server *http.Server) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	slog.Info("Initiating graceful shutdown", "timeout", "30s")

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Graceful shutdown failed, forcing close", "error", err)
		return fmt.Errorf("graceful shutdown failed: %w", err)
	}

	slog.Info("Server shutdown completed successfully")
	return nil
}
