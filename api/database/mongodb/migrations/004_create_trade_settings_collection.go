package migrations

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateTradeSettingsCollection creates the trade_settings collection with proper indexes
func CreateTradeSettingsCollection(db *mongo.Database) error {
	collectionName := "trade_settings"
	
	// Create collection with validation schema
	validator := bson.M{
		"$jsonSchema": bson.M{
			"bsonType": "object",
			"required": []string{"user", "per_trade_amount", "stop_loss_percentage", "stop_loss_status", "over_range_percentage", "take_profit_status", "take_profit_step", "tp_percentage"},
			"properties": bson.M{
				"_id": bson.M{
					"bsonType": "objectId",
					"description": "Trade settings ID must be an ObjectId",
				},
				"user": bson.M{
					"bsonType": "objectId",
					"description": "User ID reference must be an ObjectId",
				},
				"per_trade_amount": bson.M{
					"bsonType": "double",
					"minimum": 0,
					"description": "Per trade amount must be a positive number",
				},
				"stop_loss_percentage": bson.M{
					"bsonType": "int",
					"minimum": 0,
					"maximum": 100,
					"description": "Stop loss percentage must be between 0 and 100",
				},
				"stop_loss_status": bson.M{
					"bsonType": "bool",
					"description": "Stop loss status must be a boolean",
				},
				"over_range_percentage": bson.M{
					"bsonType": "double",
					"minimum": 0,
					"description": "Over range percentage must be a positive number",
				},
				"take_profit_status": bson.M{
					"bsonType": "bool",
					"description": "Take profit status must be a boolean",
				},
				"take_profit_step": bson.M{
					"bsonType": "int",
					"minimum": 1,
					"description": "Take profit step must be a positive integer",
				},
				"tp_percentage": bson.M{
					"bsonType": "array",
					"items": bson.M{
						"bsonType": "double",
						"minimum": 0,
						"maximum": 100,
					},
					"description": "Take profit percentages must be an array of numbers between 0 and 100",
				},
				"created_at": bson.M{
					"bsonType": "date",
					"description": "Creation timestamp",
				},
				"updated_at": bson.M{
					"bsonType": "date",
					"description": "Last update timestamp",
				},
			},
		},
	}

	// Create collection options
	opts := options.CreateCollection().SetValidator(validator)
	
	// Create the collection
	err := db.CreateCollection(context.Background(), collectionName, opts)
	if err != nil {
		// Collection might already exist
		if !isCollectionExistsError(err) {
			return fmt.Errorf("failed to create trade_settings collection: %v", err)
		}
	}

	// Get collection reference
	collection := db.Collection(collectionName)

	// Create indexes
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "user", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("user_unique"),
		},
		{
			Keys: bson.D{{Key: "created_at", Value: 1}},
			Options: options.Index().SetName("created_at_index"),
		},
		{
			Keys: bson.D{{Key: "stop_loss_status", Value: 1}},
			Options: options.Index().SetName("stop_loss_status_index"),
		},
		{
			Keys: bson.D{{Key: "take_profit_status", Value: 1}},
			Options: options.Index().SetName("take_profit_status_index"),
		},
	}

	_, err = collection.Indexes().CreateMany(context.Background(), indexes)
	if err != nil {
		return fmt.Errorf("failed to create indexes for trade_settings collection: %v", err)
	}

	return nil
}
