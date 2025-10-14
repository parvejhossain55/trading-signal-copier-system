package migrations

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreatePackagesCollection creates the packages collection with proper indexes
func CreatePackagesCollection(db *mongo.Database) error {
	collectionName := "packages"
	
	// Create collection with validation schema
	validator := bson.M{
		"$jsonSchema": bson.M{
			"bsonType": "object",
			"required": []string{"name", "price", "type"},
			"properties": bson.M{
				"_id": bson.M{
					"bsonType": "objectId",
					"description": "Package ID must be an ObjectId",
				},
				"name": bson.M{
					"bsonType": "string",
					"minLength": 1,
					"maxLength": 255,
					"description": "Package name is required and must be a string",
				},
				"price": bson.M{
					"bsonType": "double",
					"minimum": 0,
					"description": "Package price must be a positive number",
				},
				"type": bson.M{
					"bsonType": "string",
					"enum": []string{"basic", "premium", "enterprise", "trial"},
					"description": "Package type must be one of: basic, premium, enterprise, trial",
				},
				"description": bson.M{
					"bsonType": "string",
					"maxLength": 1000,
					"description": "Package description",
				},
				"features": bson.M{
					"bsonType": "array",
					"items": bson.M{
						"bsonType": "string",
					},
					"description": "Package features list",
				},
				"duration_days": bson.M{
					"bsonType": "int",
					"minimum": 1,
					"description": "Package duration in days",
				},
				"is_active": bson.M{
					"bsonType": "bool",
					"description": "Whether the package is active",
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
			return fmt.Errorf("failed to create packages collection: %v", err)
		}
	}

	// Get collection reference
	collection := db.Collection(collectionName)

	// Create indexes
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "name", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("name_unique"),
		},
		{
			Keys: bson.D{{Key: "type", Value: 1}},
			Options: options.Index().SetName("type_index"),
		},
		{
			Keys: bson.D{{Key: "price", Value: 1}},
			Options: options.Index().SetName("price_index"),
		},
		{
			Keys: bson.D{{Key: "is_active", Value: 1}},
			Options: options.Index().SetName("is_active_index"),
		},
		{
			Keys: bson.D{{Key: "created_at", Value: 1}},
			Options: options.Index().SetName("created_at_index"),
		},
	}

	_, err = collection.Indexes().CreateMany(context.Background(), indexes)
	if err != nil {
		return fmt.Errorf("failed to create indexes for packages collection: %v", err)
	}

	return nil
}
