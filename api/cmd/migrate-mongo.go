package cmd

import (
	"copier/database/mongodb"
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
)

var migrateMongoCmd = &cobra.Command{
	Use:   "migrate-mongo",
	Short: "Run MongoDB database migrations",
	RunE:  migrateMongo,
}

func migrateMongo(cmd *cobra.Command, args []string) error {
	slog.Info("Starting MongoDB migrations...")
	
	err := mongodb.RunMigrations()
	if err != nil {
		return fmt.Errorf("MongoDB migration failed: %v", err)
	}
	
	slog.Info("MongoDB migrations completed successfully")
	return nil
}

func init() {
	RootCmd.AddCommand(migrateMongoCmd)
}
