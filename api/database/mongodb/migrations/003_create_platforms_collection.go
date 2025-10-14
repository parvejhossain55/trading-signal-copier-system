package migrations

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreatePlatformsCollection creates the platforms collection with proper indexes
func CreatePlatformsCollection(db *mongo.Database) error {
	collectionName := "platforms"
	
	// Create collection with validation schema
	validator := bson.M{
		"$jsonSchema": bson.M{
			"bsonType": "object",
			"required": []string{"user", "name", "api_key", "api_secret"},
			"properties": bson.M{
				"_id": bson.M{
					"bsonType": "objectId",
					"description": "Platform ID must be an ObjectId",
				},
				"user": bson.M{
					"bsonType": "objectId",
					"description": "User ID reference must be an ObjectId",
				},
				"name": bson.M{
					"bsonType": "string",
					"minLength": 1,
					"maxLength": 255,
					"description": "Platform name is required and must be a string",
				},
				"api_key": bson.M{
					"bsonType": "string",
					"minLength": 1,
					"maxLength": 500,
					"description": "API key is required and must be a string",
				},
				"api_secret": bson.M{
					"bsonType": "string",
					"minLength": 1,
					"maxLength": 500,
					"description": "API secret is required and must be a string",
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
			return fmt.Errorf("failed to create platforms collection: %v", err)
		}
	}

	// Get collection reference
	collection := db.Collection(collectionName)

	// Create indexes
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "user", Value: 1}},
			Options: options.Index().SetName("user_index"),
		},
		{
			Keys: bson.D{{Key: "name", Value: 1}},
			Options: options.Index().SetName("name_index"),
		},
		{
			Keys: bson.D{{Key: "user", Value: 1}, {Key: "name", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("user_platform_unique"),
		},
		{
			Keys: bson.D{{Key: "api_key", Value: 1}},
			Options: options.Index().SetUnique(true).SetSparse(true).SetName("api_key_unique"),
		},
		{
			Keys: bson.D{{Key: "created_at", Value: 1}},
			Options: options.Index().SetName("created_at_index"),
		},
	}

	_, err = collection.Indexes().CreateMany(context.Background(), indexes)
	if err != nil {
		return fmt.Errorf("failed to create indexes for platforms collection: %v", err)
	}

	return nil
}
