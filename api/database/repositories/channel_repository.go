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

// ChannelRepository defines channel-specific repository operations
type ChannelRepository interface {
	BaseRepository
	FindByUserID(ctx context.Context, userID primitive.ObjectID) ([]*Channel, error)
	FindByChannelID(ctx context.Context, channelID string) (*Channel, error)
	FindByUserAndName(ctx context.Context, userID primitive.ObjectID, name string) (*Channel, error)
	FindByIDTyped(ctx context.Context, id primitive.ObjectID) (*Channel, error)
	CreateChannel(ctx context.Context, channel *Channel) (*mongo.InsertOneResult, error)
	UpdateChannel(ctx context.Context, id primitive.ObjectID, update *Channel) error
	DeleteChannel(ctx context.Context, id primitive.ObjectID) error
	DeleteChannelsByUser(ctx context.Context, userID primitive.ObjectID) error
	FindAllByUser(ctx context.Context, userID primitive.ObjectID, skip, limit int64) ([]*Channel, error)
	CountChannelsByUser(ctx context.Context, userID primitive.ObjectID) (int64, error)
}

// channelRepository implements ChannelRepository interface
type channelRepository struct {
	BaseRepository
	collection *mongo.Collection
}

// NewChannelRepository creates a new channel repository instance
func NewChannelRepository(collection *mongo.Collection) ChannelRepository {
	return &channelRepository{
		BaseRepository: NewBaseRepository(collection),
		collection:     collection,
	}
}

// FindByUserID finds all channels for a specific user
func (r *channelRepository) FindByUserID(ctx context.Context, userID primitive.ObjectID) ([]*Channel, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"user": userID})
	if err != nil {
		return nil, fmt.Errorf("failed to find channels by user ID: %w", err)
	}
	defer cursor.Close(ctx)

	var channels []*Channel
	for cursor.Next(ctx) {
		var channel Channel
		if err := cursor.Decode(&channel); err != nil {
			return nil, fmt.Errorf("failed to decode channel: %w", err)
		}
		channels = append(channels, &channel)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return channels, nil
}

// FindByChannelID finds a channel by its channel ID
func (r *channelRepository) FindByChannelID(ctx context.Context, channelID string) (*Channel, error) {
	var channel Channel
	err := r.collection.FindOne(ctx, bson.M{"channel_id": channelID}).Decode(&channel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("channel not found with ID: %s", channelID)
		}
		return nil, fmt.Errorf("failed to find channel by channel ID: %w", err)
	}

	return &channel, nil
}

// FindByUserAndName finds a channel by user ID and name
func (r *channelRepository) FindByUserAndName(ctx context.Context, userID primitive.ObjectID, name string) (*Channel, error) {
	var channel Channel
	err := r.collection.FindOne(ctx, bson.M{"user": userID, "name": name}).Decode(&channel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("channel not found with user ID: %s and name: %s", userID.Hex(), name)
		}
		return nil, fmt.Errorf("failed to find channel by user and name: %w", err)
	}

	return &channel, nil
}

// FindByIDTyped finds a channel by ID and returns typed Channel struct
func (r *channelRepository) FindByIDTyped(ctx context.Context, id primitive.ObjectID) (*Channel, error) {
	var channel Channel
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&channel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("channel not found with ID: %s", id.Hex())
		}
		return nil, fmt.Errorf("failed to find channel by ID: %w", err)
	}

	return &channel, nil
}

// CreateChannel creates a new channel
func (r *channelRepository) CreateChannel(ctx context.Context, channel *Channel) (*mongo.InsertOneResult, error) {
	result, err := r.collection.InsertOne(ctx, channel)
	if err != nil {
		return nil, fmt.Errorf("failed to create channel: %w", err)
	}

	return result, nil
}

// UpdateChannel updates an existing channel
func (r *channelRepository) UpdateChannel(ctx context.Context, id primitive.ObjectID, update *Channel) error {
	update.UpdatedAt = time.Now()
	
	_, err := r.collection.UpdateByID(ctx, id, bson.M{"$set": update})
	if err != nil {
		return fmt.Errorf("failed to update channel: %w", err)
	}

	return nil
}

// DeleteChannel deletes a channel by ID
func (r *channelRepository) DeleteChannel(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete channel: %w", err)
	}

	return nil
}

// DeleteChannelsByUser deletes all channels for a specific user
func (r *channelRepository) DeleteChannelsByUser(ctx context.Context, userID primitive.ObjectID) error {
	_, err := r.collection.DeleteMany(ctx, bson.M{"user": userID})
	if err != nil {
		return fmt.Errorf("failed to delete channels by user: %w", err)
	}

	return nil
}

// FindAllByUser retrieves all channels for a user with pagination
func (r *channelRepository) FindAllByUser(ctx context.Context, userID primitive.ObjectID, skip, limit int64) ([]*Channel, error) {
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
		return nil, fmt.Errorf("failed to find channels by user: %w", err)
	}
	defer cursor.Close(ctx)

	var channels []*Channel
	for cursor.Next(ctx) {
		var channel Channel
		if err := cursor.Decode(&channel); err != nil {
			return nil, fmt.Errorf("failed to decode channel: %w", err)
		}
		channels = append(channels, &channel)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return channels, nil
}

// CountChannelsByUser returns the number of channels for a specific user
func (r *channelRepository) CountChannelsByUser(ctx context.Context, userID primitive.ObjectID) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"user": userID})
	if err != nil {
		return 0, fmt.Errorf("failed to count channels by user: %w", err)
	}

	return count, nil
}
