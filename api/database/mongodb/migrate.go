package mongodb

import (
	"context"
	"copier/config"
	"copier/database/mongodb/migrations"
	"copier/internal/logger"
	"fmt"
	"log/slog"
	"os"
	"sort"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Migration represents a MongoDB migration
type Migration struct {
	ID          string    `bson:"_id,omitempty"`
	Filename    string    `bson:"filename"`
	ExecutedAt  time.Time `bson:"executed_at"`
	Description string    `bson:"description,omitempty"`
}

// RunMigrations executes all pending MongoDB migrations
func RunMigrations() error {
	// Load configuration
	conf := config.GetConfig()
	
	// Initialize logger
	logger.SetupLogger(conf.ServiceName, string(conf.Mode))
	
	// Connect to MongoDB
	client, err := config.NewMongoDB()
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	// Get database name from config
	dbName := conf.Database.Database
	if dbName == "" {
		dbName = "copier" // default database name
	}
	
	db := client.Database(dbName)

	// Get migrations directory
	migrationsDir := "./database/mongodb/migrations"

	// Read all migration files
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %v", err)
	}

	// Filter and sort migration files
	var migrationFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".go") && file.Name() != "utils.go" {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}
	sort.Strings(migrationFiles)

	slog.Info("Starting MongoDB migrations", "count", len(migrationFiles))

	// Create migrations collection if it doesn't exist
	migrationsCollection := db.Collection("migrations")
	
	// Create index on filename for uniqueness
	_, err = migrationsCollection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "filename", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		// Index might already exist, continue
		slog.Debug("Index creation result", "error", err)
	}

	// Execute each migration
	for _, filename := range migrationFiles {
		// Check if migration already executed
		var existingMigration Migration
		err := migrationsCollection.FindOne(
			context.Background(),
			bson.M{"filename": filename},
		).Decode(&existingMigration)

		if err == nil {
			slog.Info("Migration already executed", "filename", filename)
			continue
		} else if err != mongo.ErrNoDocuments {
			return fmt.Errorf("failed to check migration status: %v", err)
		}

		// Execute migration
		err = executeMigration(db, filename)
		if err != nil {
			return fmt.Errorf("migration %s failed: %v", filename, err)
		}

		// Record migration execution
		migration := Migration{
			Filename:   filename,
			ExecutedAt: time.Now(),
		}
		
		_, err = migrationsCollection.InsertOne(context.Background(), migration)
		if err != nil {
			return fmt.Errorf("failed to record migration %s: %v", filename, err)
		}

		slog.Info("Migration executed successfully", "filename", filename)
	}

	slog.Info("All MongoDB migrations completed successfully")
	return nil
}

// executeMigration runs a specific migration based on filename
func executeMigration(db *mongo.Database, filename string) error {
	// Remove .go extension to get migration name
	migrationName := strings.TrimSuffix(filename, ".go")
	
	switch migrationName {
	case "001_create_users_collection":
		return migrations.CreateUsersCollection(db)
	case "002_create_channels_collection":
		return migrations.CreateChannelsCollection(db)
	case "003_create_platforms_collection":
		return migrations.CreatePlatformsCollection(db)
	case "004_create_trade_settings_collection":
		return migrations.CreateTradeSettingsCollection(db)
	case "005_create_packages_collection":
		return migrations.CreatePackagesCollection(db)
	case "006_create_subscribe_packages_collection":
		return migrations.CreateSubscribePackagesCollection(db)
	default:
		return fmt.Errorf("unknown migration: %s", migrationName)
	}
}
