package config

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewMongoDB creates and returns a new MongoDB client connection.
func NewMongoDB() (*mongo.Client, error) {
	cfg := GetConfig()

	// Use the GetDatabaseURL method which handles SSL configuration properly
	connStr := cfg.GetDatabaseURL()
	if connStr == "" {
		return nil, fmt.Errorf("no database configuration found")
	}

	// Set client options
	clientOptions := options.Client().ApplyURI(connStr)
	clientOptions.SetMaxPoolSize(uint64(cfg.Database.MaxPoolSize))
	clientOptions.SetMinPoolSize(uint64(cfg.Database.MinPoolSize))
	clientOptions.SetMaxConnIdleTime(time.Duration(cfg.Database.MaxConnIdleTimeInMs) * time.Millisecond)

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Ping the database to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}
