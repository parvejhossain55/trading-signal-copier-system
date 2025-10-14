package mongodb

import (
	"context"
	"copier/config"
	"fmt"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// SeedData inserts sample data into the database for testing
func SeedData() error {
	// Load configuration
	conf := config.GetConfig()
	
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

	// Seed users
	err = seedUsers(db)
	if err != nil {
		return fmt.Errorf("failed to seed users: %v", err)
	}

	// Seed packages
	err = seedPackages(db)
	if err != nil {
		return fmt.Errorf("failed to seed packages: %v", err)
	}

	slog.Info("Database seeding completed successfully")
	return nil
}

// seedUsers inserts sample user data
func seedUsers(db *mongo.Database) error {
	collection := db.Collection("users")
	
	users := []interface{}{
		map[string]interface{}{
			"name":       "John Doe",
			"email":      "john.doe@example.com",
			"phone":      "+1234567890",
			"session":    "session_123",
			"created_at": time.Now(),
			"updated_at": time.Now(),
		},
		map[string]interface{}{
			"name":       "Jane Smith",
			"email":      "jane.smith@example.com",
			"phone":      "+0987654321",
			"session":    "session_456",
			"created_at": time.Now(),
			"updated_at": time.Now(),
		},
	}

	_, err := collection.InsertMany(context.Background(), users)
	if err != nil {
		return err
	}

	slog.Info("Seeded users collection")
	return nil
}

// seedPackages inserts sample package data
func seedPackages(db *mongo.Database) error {
	collection := db.Collection("packages")
	
	packages := []interface{}{
		map[string]interface{}{
			"name":          "Basic Plan",
			"price":         29.99,
			"type":          "basic",
			"description":   "Basic trading features",
			"features":      []string{"Basic signals", "Email support"},
			"duration_days": 30,
			"is_active":     true,
			"created_at":    time.Now(),
			"updated_at":    time.Now(),
		},
		map[string]interface{}{
			"name":          "Premium Plan",
			"price":         99.99,
			"type":          "premium",
			"description":   "Advanced trading features",
			"features":      []string{"Advanced signals", "Priority support", "Custom strategies"},
			"duration_days": 30,
			"is_active":     true,
			"created_at":    time.Now(),
			"updated_at":    time.Now(),
		},
		map[string]interface{}{
			"name":          "Enterprise Plan",
			"price":         299.99,
			"type":          "enterprise",
			"description":   "Full-featured trading platform",
			"features":      []string{"All signals", "24/7 support", "Custom strategies", "API access"},
			"duration_days": 30,
			"is_active":     true,
			"created_at":    time.Now(),
			"updated_at":    time.Now(),
		},
	}

	_, err := collection.InsertMany(context.Background(), packages)
	if err != nil {
		return err
	}

	slog.Info("Seeded packages collection")
	return nil
}

// GetSampleUserID returns a sample user ID for testing
func GetSampleUserID(db *mongo.Database) (primitive.ObjectID, error) {
	collection := db.Collection("users")
	
	var user struct {
		ID primitive.ObjectID `bson:"_id"`
	}
	
	err := collection.FindOne(context.Background(), map[string]interface{}{}).Decode(&user)
	if err != nil {
		return primitive.NilObjectID, err
	}
	
	return user.ID, nil
}

// GetSamplePackageID returns a sample package ID for testing
func GetSamplePackageID(db *mongo.Database) (primitive.ObjectID, error) {
	collection := db.Collection("packages")
	
	var pkg struct {
		ID primitive.ObjectID `bson:"_id"`
	}
	
	err := collection.FindOne(context.Background(), map[string]interface{}{}).Decode(&pkg)
	if err != nil {
		return primitive.NilObjectID, err
	}
	
	return pkg.ID, nil
}
