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

// TradeSettingsRepository defines trade settings-specific repository operations
type TradeSettingsRepository interface {
	BaseRepository
	FindByUserID(ctx context.Context, userID primitive.ObjectID) (*TradeSettings, error)
	FindByIDTyped(ctx context.Context, id primitive.ObjectID) (*TradeSettings, error)
	CreateTradeSettings(ctx context.Context, settings *TradeSettings) (*mongo.InsertOneResult, error)
	UpdateTradeSettings(ctx context.Context, id primitive.ObjectID, update *TradeSettings) error
	UpdateByUserID(ctx context.Context, userID primitive.ObjectID, update *TradeSettings) error
	DeleteTradeSettings(ctx context.Context, id primitive.ObjectID) error
	DeleteByUserID(ctx context.Context, userID primitive.ObjectID) error
	FindAllWithStopLoss(ctx context.Context, skip, limit int64) ([]*TradeSettings, error)
	FindAllWithTakeProfit(ctx context.Context, skip, limit int64) ([]*TradeSettings, error)
	CountTradeSettings(ctx context.Context) (int64, error)
	UpdateStopLossSettings(ctx context.Context, userID primitive.ObjectID, percentage int, status bool) error
	UpdateTakeProfitSettings(ctx context.Context, userID primitive.ObjectID, status bool, step int, percentages []float64) error
}

// tradeSettingsRepository implements TradeSettingsRepository interface
type tradeSettingsRepository struct {
	BaseRepository
	collection *mongo.Collection
}

// NewTradeSettingsRepository creates a new trade settings repository instance
func NewTradeSettingsRepository(collection *mongo.Collection) TradeSettingsRepository {
	return &tradeSettingsRepository{
		BaseRepository: NewBaseRepository(collection),
		collection:     collection,
	}
}

// FindByUserID finds trade settings for a specific user
func (r *tradeSettingsRepository) FindByUserID(ctx context.Context, userID primitive.ObjectID) (*TradeSettings, error) {
	var settings TradeSettings
	err := r.collection.FindOne(ctx, bson.M{"user": userID}).Decode(&settings)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("trade settings not found for user ID: %s", userID.Hex())
		}
		return nil, fmt.Errorf("failed to find trade settings by user ID: %w", err)
	}

	return &settings, nil
}

// FindByIDTyped finds trade settings by ID and returns typed TradeSettings struct
func (r *tradeSettingsRepository) FindByIDTyped(ctx context.Context, id primitive.ObjectID) (*TradeSettings, error) {
	var settings TradeSettings
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&settings)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("trade settings not found with ID: %s", id.Hex())
		}
		return nil, fmt.Errorf("failed to find trade settings by ID: %w", err)
	}

	return &settings, nil
}

// CreateTradeSettings creates new trade settings
func (r *tradeSettingsRepository) CreateTradeSettings(ctx context.Context, settings *TradeSettings) (*mongo.InsertOneResult, error) {
	result, err := r.collection.InsertOne(ctx, settings)
	if err != nil {
		return nil, fmt.Errorf("failed to create trade settings: %w", err)
	}

	return result, nil
}

// UpdateTradeSettings updates existing trade settings by ID
func (r *tradeSettingsRepository) UpdateTradeSettings(ctx context.Context, id primitive.ObjectID, update *TradeSettings) error {
	update.UpdatedAt = time.Now()
	
	_, err := r.collection.UpdateByID(ctx, id, bson.M{"$set": update})
	if err != nil {
		return fmt.Errorf("failed to update trade settings: %w", err)
	}

	return nil
}

// UpdateByUserID updates trade settings by user ID
func (r *tradeSettingsRepository) UpdateByUserID(ctx context.Context, userID primitive.ObjectID, update *TradeSettings) error {
	update.UpdatedAt = time.Now()
	
	_, err := r.collection.UpdateOne(ctx, bson.M{"user": userID}, bson.M{"$set": update})
	if err != nil {
		return fmt.Errorf("failed to update trade settings by user ID: %w", err)
	}

	return nil
}

// DeleteTradeSettings deletes trade settings by ID
func (r *tradeSettingsRepository) DeleteTradeSettings(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete trade settings: %w", err)
	}

	return nil
}

// DeleteByUserID deletes trade settings by user ID
func (r *tradeSettingsRepository) DeleteByUserID(ctx context.Context, userID primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"user": userID})
	if err != nil {
		return fmt.Errorf("failed to delete trade settings by user ID: %w", err)
	}

	return nil
}

// FindAllWithStopLoss finds all trade settings with stop loss enabled
func (r *tradeSettingsRepository) FindAllWithStopLoss(ctx context.Context, skip, limit int64) ([]*TradeSettings, error) {
	opts := options.Find()
	if skip > 0 {
		opts.SetSkip(skip)
	}
	if limit > 0 {
		opts.SetLimit(limit)
	}
	opts.SetSort(bson.M{"created_at": -1})

	cursor, err := r.collection.Find(ctx, bson.M{"stop_loss_status": true}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find trade settings with stop loss: %w", err)
	}
	defer cursor.Close(ctx)

	var settings []*TradeSettings
	for cursor.Next(ctx) {
		var setting TradeSettings
		if err := cursor.Decode(&setting); err != nil {
			return nil, fmt.Errorf("failed to decode trade settings: %w", err)
		}
		settings = append(settings, &setting)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return settings, nil
}

// FindAllWithTakeProfit finds all trade settings with take profit enabled
func (r *tradeSettingsRepository) FindAllWithTakeProfit(ctx context.Context, skip, limit int64) ([]*TradeSettings, error) {
	opts := options.Find()
	if skip > 0 {
		opts.SetSkip(skip)
	}
	if limit > 0 {
		opts.SetLimit(limit)
	}
	opts.SetSort(bson.M{"created_at": -1})

	cursor, err := r.collection.Find(ctx, bson.M{"take_profit_status": true}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find trade settings with take profit: %w", err)
	}
	defer cursor.Close(ctx)

	var settings []*TradeSettings
	for cursor.Next(ctx) {
		var setting TradeSettings
		if err := cursor.Decode(&setting); err != nil {
			return nil, fmt.Errorf("failed to decode trade settings: %w", err)
		}
		settings = append(settings, &setting)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return settings, nil
}

// CountTradeSettings returns the total number of trade settings
func (r *tradeSettingsRepository) CountTradeSettings(ctx context.Context) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, fmt.Errorf("failed to count trade settings: %w", err)
	}

	return count, nil
}

// UpdateStopLossSettings updates only stop loss related settings
func (r *tradeSettingsRepository) UpdateStopLossSettings(ctx context.Context, userID primitive.ObjectID, percentage int, status bool) error {
	update := bson.M{
		"stop_loss_percentage": percentage,
		"stop_loss_status":     status,
		"updated_at":           bson.M{"$currentDate": true},
	}

	_, err := r.collection.UpdateOne(ctx, bson.M{"user": userID}, bson.M{"$set": update})
	if err != nil {
		return fmt.Errorf("failed to update stop loss settings: %w", err)
	}

	return nil
}

// UpdateTakeProfitSettings updates only take profit related settings
func (r *tradeSettingsRepository) UpdateTakeProfitSettings(ctx context.Context, userID primitive.ObjectID, status bool, step int, percentages []float64) error {
	update := bson.M{
		"take_profit_status": status,
		"take_profit_step":   step,
		"tp_percentage":      percentages,
		"updated_at":         bson.M{"$currentDate": true},
	}

	_, err := r.collection.UpdateOne(ctx, bson.M{"user": userID}, bson.M{"$set": update})
	if err != nil {
		return fmt.Errorf("failed to update take profit settings: %w", err)
	}

	return nil
}
