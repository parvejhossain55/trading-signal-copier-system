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

// PlatformRepository defines platform-specific repository operations
type PlatformRepository interface {
	BaseRepository
	FindByUserID(ctx context.Context, userID primitive.ObjectID) ([]*Platform, error)
	FindByUserAndName(ctx context.Context, userID primitive.ObjectID, name string) (*Platform, error)
	FindByAPIKey(ctx context.Context, apiKey string) (*Platform, error)
	FindByIDTyped(ctx context.Context, id primitive.ObjectID) (*Platform, error)
	CreatePlatform(ctx context.Context, platform *Platform) (*mongo.InsertOneResult, error)
	UpdatePlatform(ctx context.Context, id primitive.ObjectID, update *Platform) error
	DeletePlatform(ctx context.Context, id primitive.ObjectID) error
	DeletePlatformsByUser(ctx context.Context, userID primitive.ObjectID) error
	FindAllByUser(ctx context.Context, userID primitive.ObjectID, skip, limit int64) ([]*Platform, error)
	CountPlatformsByUser(ctx context.Context, userID primitive.ObjectID) (int64, error)
	UpdateAPIKeys(ctx context.Context, id primitive.ObjectID, apiKey, apiSecret string) error
}

// platformRepository implements PlatformRepository interface
type platformRepository struct {
	BaseRepository
	collection *mongo.Collection
}

// NewPlatformRepository creates a new platform repository instance
func NewPlatformRepository(collection *mongo.Collection) PlatformRepository {
	return &platformRepository{
		BaseRepository: NewBaseRepository(collection),
		collection:     collection,
	}
}

// FindByUserID finds all platforms for a specific user
func (r *platformRepository) FindByUserID(ctx context.Context, userID primitive.ObjectID) ([]*Platform, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"user": userID})
	if err != nil {
		return nil, fmt.Errorf("failed to find platforms by user ID: %w", err)
	}
	defer cursor.Close(ctx)

	var platforms []*Platform
	for cursor.Next(ctx) {
		var platform Platform
		if err := cursor.Decode(&platform); err != nil {
			return nil, fmt.Errorf("failed to decode platform: %w", err)
		}
		platforms = append(platforms, &platform)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return platforms, nil
}

// FindByUserAndName finds a platform by user ID and name
func (r *platformRepository) FindByUserAndName(ctx context.Context, userID primitive.ObjectID, name string) (*Platform, error) {
	var platform Platform
	err := r.collection.FindOne(ctx, bson.M{"user": userID, "name": name}).Decode(&platform)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("platform not found with user ID: %s and name: %s", userID.Hex(), name)
		}
		return nil, fmt.Errorf("failed to find platform by user and name: %w", err)
	}

	return &platform, nil
}

// FindByAPIKey finds a platform by API key
func (r *platformRepository) FindByAPIKey(ctx context.Context, apiKey string) (*Platform, error) {
	var platform Platform
	err := r.collection.FindOne(ctx, bson.M{"api_key": apiKey}).Decode(&platform)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("platform not found with API key")
		}
		return nil, fmt.Errorf("failed to find platform by API key: %w", err)
	}

	return &platform, nil
}

// FindByIDTyped finds a platform by ID and returns typed Platform struct
func (r *platformRepository) FindByIDTyped(ctx context.Context, id primitive.ObjectID) (*Platform, error) {
	var platform Platform
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&platform)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("platform not found with ID: %s", id.Hex())
		}
		return nil, fmt.Errorf("failed to find platform by ID: %w", err)
	}

	return &platform, nil
}

// CreatePlatform creates a new platform
func (r *platformRepository) CreatePlatform(ctx context.Context, platform *Platform) (*mongo.InsertOneResult, error) {
	result, err := r.collection.InsertOne(ctx, platform)
	if err != nil {
		return nil, fmt.Errorf("failed to create platform: %w", err)
	}

	return result, nil
}

// UpdatePlatform updates an existing platform
func (r *platformRepository) UpdatePlatform(ctx context.Context, id primitive.ObjectID, update *Platform) error {
	update.UpdatedAt = time.Now()
	
	_, err := r.collection.UpdateByID(ctx, id, bson.M{"$set": update})
	if err != nil {
		return fmt.Errorf("failed to update platform: %w", err)
	}

	return nil
}

// DeletePlatform deletes a platform by ID
func (r *platformRepository) DeletePlatform(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete platform: %w", err)
	}

	return nil
}

// DeletePlatformsByUser deletes all platforms for a specific user
func (r *platformRepository) DeletePlatformsByUser(ctx context.Context, userID primitive.ObjectID) error {
	_, err := r.collection.DeleteMany(ctx, bson.M{"user": userID})
	if err != nil {
		return fmt.Errorf("failed to delete platforms by user: %w", err)
	}

	return nil
}

// FindAllByUser retrieves all platforms for a user with pagination
func (r *platformRepository) FindAllByUser(ctx context.Context, userID primitive.ObjectID, skip, limit int64) ([]*Platform, error) {
	opts := options.Find()
	if skip > 0 {
		opts.SetSkip(skip)
	}
	if limit > 0 {
		opts.SetLimit(limit)
	}
	opts.SetSort(bson.M{"created_at": -1}) // Sort by creation date descending

	cursor, err := r.collection.Find(ctx, bson.M{"user": userID}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find platforms by user: %w", err)
	}
	defer cursor.Close(ctx)

	var platforms []*Platform
	for cursor.Next(ctx) {
		var platform Platform
		if err := cursor.Decode(&platform); err != nil {
			return nil, fmt.Errorf("failed to decode platform: %w", err)
		}
		platforms = append(platforms, &platform)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return platforms, nil
}

// CountPlatformsByUser returns the number of platforms for a specific user
func (r *platformRepository) CountPlatformsByUser(ctx context.Context, userID primitive.ObjectID) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"user": userID})
	if err != nil {
		return 0, fmt.Errorf("failed to count platforms by user: %w", err)
	}

	return count, nil
}

// UpdateAPIKeys updates only the API key and secret for a platform
func (r *platformRepository) UpdateAPIKeys(ctx context.Context, id primitive.ObjectID, apiKey, apiSecret string) error {
	update := bson.M{
		"api_key":    apiKey,
		"api_secret": apiSecret,
		"updated_at": bson.M{"$currentDate": true},
	}

	_, err := r.collection.UpdateByID(ctx, id, bson.M{"$set": update})
	if err != nil {
		return fmt.Errorf("failed to update platform API keys: %w", err)
	}

	return nil
}
