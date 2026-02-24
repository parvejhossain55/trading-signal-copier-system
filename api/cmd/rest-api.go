package cmd

import (
	"context"
	"copier/config"
	"copier/http/routes"
	"copier/internal/database/models"
	"copier/internal/di"
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
	conf := config.GetConfig()

	logger.SetupLogger(conf.ServiceName, string(conf.Mode))

	db, err := config.NewPostgresDB()
	if err != nil {
		fmt.Println("Failed to connect to PostgreSQL", "error", err)
		return fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	// AutoMigrate models
	err = db.AutoMigrate(
		&models.User{},
		&models.Channel{},
		&models.Platform{},
		&models.TradeSettings{},
		&models.Package{},
		&models.SubscribePackage{},
	)
	if err != nil {
		slog.Error("Failed to run auto-migration", "error", err)
		return fmt.Errorf("failed to run auto-migration: %w", err)
	}

	fmt.Println("PostgreSQL connection established successfully",
		"host", conf.Database.Host,
		"port", conf.Database.Port,
		"database", conf.Database.Database)

	// Initialize DI Container
	container := di.NewContainer(db)

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:         ":" + strconv.Itoa(conf.HttpPort),
		Handler:      routes.SetupRoutes(mux, container),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	serverErrors := make(chan error, 1)
	go func() {
		slog.Info("Starting HTTP server", "port", conf.HttpPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrors <- fmt.Errorf("server error: %w", err)
		}
	}()

	select {
	case err := <-serverErrors:
		slog.Error("Server error occurred", "error", err)
		return err
	case <-shutdown():
		slog.Info("Shutdown signal received, starting graceful shutdown")
		return gracefulShutdown(server)
	}
}

func shutdown() chan os.Signal {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown,
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
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
