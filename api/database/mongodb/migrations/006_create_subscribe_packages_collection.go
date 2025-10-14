package migrations

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateSubscribePackagesCollection creates the subscribe_packages collection with proper indexes
func CreateSubscribePackagesCollection(db *mongo.Database) error {
	collectionName := "subscribe_packages"
	
	// Create collection with validation schema
	validator := bson.M{
		"$jsonSchema": bson.M{
			"bsonType": "object",
			"required": []string{"user", "package_id", "price", "payment_method"},
			"properties": bson.M{
				"_id": bson.M{
					"bsonType": "objectId",
					"description": "Subscription ID must be an ObjectId",
				},
				"user": bson.M{
					"bsonType": "objectId",
					"description": "User ID reference must be an ObjectId",
				},
				"package_id": bson.M{
					"bsonType": "objectId",
					"description": "Package ID reference must be an ObjectId",
				},
				"price": bson.M{
					"bsonType": "double",
					"minimum": 0,
					"description": "Subscription price must be a positive number",
				},
				"tnx_hash": bson.M{
					"bsonType": "string",
					"maxLength": 255,
					"description": "Transaction hash",
				},
				"binance_pay_id": bson.M{
					"bsonType": "string",
					"maxLength": 255,
					"description": "Binance Pay ID",
				},
				"payment_method": bson.M{
					"bsonType": "string",
					"enum": []string{"crypto", "binance_pay", "credit_card", "paypal"},
					"description": "Payment method must be one of: crypto, binance_pay, credit_card, paypal",
				},
				"status": bson.M{
					"bsonType": "string",
					"enum": []string{"pending", "completed", "failed", "cancelled", "refunded"},
					"description": "Payment status",
				},
				"start_date": bson.M{
					"bsonType": "date",
					"description": "Subscription start date",
				},
				"end_date": bson.M{
					"bsonType": "date",
					"description": "Subscription end date",
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
			return fmt.Errorf("failed to create subscribe_packages collection: %v", err)
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
			Keys: bson.D{{Key: "package_id", Value: 1}},
			Options: options.Index().SetName("package_id_index"),
		},
		{
			Keys: bson.D{{Key: "tnx_hash", Value: 1}},
			Options: options.Index().SetUnique(true).SetSparse(true).SetName("tnx_hash_unique"),
		},
		{
			Keys: bson.D{{Key: "binance_pay_id", Value: 1}},
			Options: options.Index().SetUnique(true).SetSparse(true).SetName("binance_pay_id_unique"),
		},
		{
			Keys: bson.D{{Key: "payment_method", Value: 1}},
			Options: options.Index().SetName("payment_method_index"),
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}},
			Options: options.Index().SetName("status_index"),
		},
		{
			Keys: bson.D{{Key: "user", Value: 1}, {Key: "package_id", Value: 1}},
			Options: options.Index().SetName("user_package_index"),
		},
		{
			Keys: bson.D{{Key: "start_date", Value: 1}},
			Options: options.Index().SetName("start_date_index"),
		},
		{
			Keys: bson.D{{Key: "end_date", Value: 1}},
			Options: options.Index().SetName("end_date_index"),
		},
		{
			Keys: bson.D{{Key: "created_at", Value: 1}},
			Options: options.Index().SetName("created_at_index"),
		},
	}

	_, err = collection.Indexes().CreateMany(context.Background(), indexes)
	if err != nil {
		return fmt.Errorf("failed to create indexes for subscribe_packages collection: %v", err)
	}

	return nil
}
