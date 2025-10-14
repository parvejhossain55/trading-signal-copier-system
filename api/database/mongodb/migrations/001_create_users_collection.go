package migrations

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateUsersCollection creates the users collection with proper indexes
func CreateUsersCollection(db *mongo.Database) error {
	collectionName := "users"
	
	// Create collection with validation schema
	validator := bson.M{
		"$jsonSchema": bson.M{
			"bsonType": "object",
			"required": []string{"name", "email"},
			"properties": bson.M{
				"_id": bson.M{
					"bsonType": "objectId",
					"description": "User ID must be an ObjectId",
				},
				"name": bson.M{
					"bsonType": "string",
					"minLength": 1,
					"maxLength": 255,
					"description": "User name is required and must be a string",
				},
				"email": bson.M{
					"bsonType": "string",
					"pattern": "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$",
					"description": "Email must be a valid email address",
				},
				"phone": bson.M{
					"bsonType": "string",
					"description": "Phone number",
				},
				"session": bson.M{
					"bsonType": "string",
					"description": "User session token",
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
			return fmt.Errorf("failed to create users collection: %v", err)
		}
	}

	// Get collection reference
	collection := db.Collection(collectionName)

	// Create indexes
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("email_unique"),
		},
		{
			Keys: bson.D{{Key: "phone", Value: 1}},
			Options: options.Index().SetUnique(true).SetSparse(true).SetName("phone_unique"),
		},
		{
			Keys: bson.D{{Key: "session", Value: 1}},
			Options: options.Index().SetSparse(true).SetName("session_index"),
		},
		{
			Keys: bson.D{{Key: "created_at", Value: 1}},
			Options: options.Index().SetName("created_at_index"),
		},
	}

	_, err = collection.Indexes().CreateMany(context.Background(), indexes)
	if err != nil {
		return fmt.Errorf("failed to create indexes for users collection: %v", err)
	}

	return nil
}

