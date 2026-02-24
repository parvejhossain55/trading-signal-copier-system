package cmd

import (
	"copier/config"
	"copier/internal/database/models"
	"copier/internal/logger"
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	RunE:  migrate,
}

func migrate(cmd *cobra.Command, args []string) error {
	// Load configuration
	conf := config.GetConfig()

	// Initialize logger
	logger.SetupLogger(conf.ServiceName, string(conf.Mode))

	// Connect to database
	db, err := config.NewPostgresDB()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	slog.Info("Starting database migrations using GORM AutoMigrate")

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
		return fmt.Errorf("failed to run auto-migration: %w", err)
	}

	slog.Info("All migrations completed successfully")
	return nil
}
