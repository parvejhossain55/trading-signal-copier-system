package migrations

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateChannelsCollection creates the channels collection with proper indexes
func CreateChannelsCollection(db *mongo.Database) error {
	collectionName := "channels"
	
	// Create collection with validation schema
	validator := bson.M{
		"$jsonSchema": bson.M{
			"bsonType": "object",
			"required": []string{"user", "name", "channel_id"},
			"properties": bson.M{
				"_id": bson.M{
					"bsonType": "objectId",
					"description": "Channel ID must be an ObjectId",
				},
				"user": bson.M{
					"bsonType": "objectId",
					"description": "User ID reference must be an ObjectId",
				},
				"name": bson.M{
					"bsonType": "string",
					"minLength": 1,
					"maxLength": 255,
					"description": "Channel name is required and must be a string",
				},
				"channel_id": bson.M{
					"bsonType": "string",
					"minLength": 1,
					"maxLength": 255,
					"description": "Channel ID is required and must be a string",
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
			return fmt.Errorf("failed to create channels collection: %v", err)
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
			Keys: bson.D{{Key: "channel_id", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("channel_id_unique"),
		},
		{
			Keys: bson.D{{Key: "user", Value: 1}, {Key: "name", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("user_name_unique"),
		},
		{
			Keys: bson.D{{Key: "created_at", Value: 1}},
			Options: options.Index().SetName("created_at_index"),
		},
	}

	_, err = collection.Indexes().CreateMany(context.Background(), indexes)
	if err != nil {
		return fmt.Errorf("failed to create indexes for channels collection: %v", err)
	}

	return nil
}
