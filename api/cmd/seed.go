package cmd

import (
	"copier/database/mongodb"
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
)

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed the database with sample data",
	RunE:  seed,
}

func seed(cmd *cobra.Command, args []string) error {
	slog.Info("Starting database seeding...")
	
	err := mongodb.SeedData()
	if err != nil {
		return fmt.Errorf("database seeding failed: %v", err)
	}
	
	slog.Info("Database seeding completed successfully")
	return nil
}

func init() {
	RootCmd.AddCommand(seedCmd)
}
