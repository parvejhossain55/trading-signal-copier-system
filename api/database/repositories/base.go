package repositories

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// BaseRepository defines common repository operations
type BaseRepository interface {
	Create(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (interface{}, error)
	FindOne(ctx context.Context, filter bson.M) (interface{}, error)
	FindMany(ctx context.Context, filter bson.M, opts ...*options.FindOptions) (*mongo.Cursor, error)
	UpdateByID(ctx context.Context, id primitive.ObjectID, update bson.M) (*mongo.UpdateResult, error)
	UpdateOne(ctx context.Context, filter bson.M, update bson.M) (*mongo.UpdateResult, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) (*mongo.DeleteResult, error)
	DeleteOne(ctx context.Context, filter bson.M) (*mongo.DeleteResult, error)
	Count(ctx context.Context, filter bson.M) (int64, error)
}

// baseRepository implements common repository operations
type baseRepository struct {
	collection *mongo.Collection
}

// NewBaseRepository creates a new base repository instance
func NewBaseRepository(collection *mongo.Collection) BaseRepository {
	return &baseRepository{
		collection: collection,
	}
}

// Create inserts a new document into the collection
func (r *baseRepository) Create(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error) {
	// Set timestamps
	if doc, ok := document.(bson.M); ok {
		now := time.Now()
		doc["created_at"] = now
		doc["updated_at"] = now
	}

	result, err := r.collection.InsertOne(ctx, document)
	if err != nil {
		return nil, fmt.Errorf("failed to create document: %w", err)
	}

	return result, nil
}

// FindByID finds a document by its ObjectID
func (r *baseRepository) FindByID(ctx context.Context, id primitive.ObjectID) (interface{}, error) {
	var result bson.M
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("document not found")
		}
		return nil, fmt.Errorf("failed to find document by ID: %w", err)
	}

	return result, nil
}

// FindOne finds a single document matching the filter
func (r *baseRepository) FindOne(ctx context.Context, filter bson.M) (interface{}, error) {
	var result bson.M
	err := r.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("document not found")
		}
		return nil, fmt.Errorf("failed to find document: %w", err)
	}

	return result, nil
}

// FindMany finds multiple documents matching the filter
func (r *baseRepository) FindMany(ctx context.Context, filter bson.M, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	cursor, err := r.collection.Find(ctx, filter, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %w", err)
	}

	return cursor, nil
}

// UpdateByID updates a document by its ObjectID
func (r *baseRepository) UpdateByID(ctx context.Context, id primitive.ObjectID, update bson.M) (*mongo.UpdateResult, error) {
	// Add updated_at timestamp
	update["updated_at"] = time.Now()

	result, err := r.collection.UpdateByID(ctx, id, bson.M{"$set": update})
	if err != nil {
		return nil, fmt.Errorf("failed to update document by ID: %w", err)
	}

	return result, nil
}

// UpdateOne updates a single document matching the filter
func (r *baseRepository) UpdateOne(ctx context.Context, filter bson.M, update bson.M) (*mongo.UpdateResult, error) {
	// Add updated_at timestamp
	update["updated_at"] = time.Now()

	result, err := r.collection.UpdateOne(ctx, filter, bson.M{"$set": update})
	if err != nil {
		return nil, fmt.Errorf("failed to update document: %w", err)
	}

	return result, nil
}

// DeleteByID deletes a document by its ObjectID
func (r *baseRepository) DeleteByID(ctx context.Context, id primitive.ObjectID) (*mongo.DeleteResult, error) {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return nil, fmt.Errorf("failed to delete document by ID: %w", err)
	}

	return result, nil
}

// DeleteOne deletes a single document matching the filter
func (r *baseRepository) DeleteOne(ctx context.Context, filter bson.M) (*mongo.DeleteResult, error) {
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to delete document: %w", err)
	}

	return result, nil
}

// Count returns the number of documents matching the filter
func (r *baseRepository) Count(ctx context.Context, filter bson.M) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count documents: %w", err)
	}

	return count, nil
}
