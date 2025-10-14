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

// SubscribePackageRepository defines subscribe package-specific repository operations
type SubscribePackageRepository interface {
	BaseRepository
	FindByUserID(ctx context.Context, userID primitive.ObjectID) ([]*SubscribePackage, error)
	FindByPackageID(ctx context.Context, packageID primitive.ObjectID) ([]*SubscribePackage, error)
	FindByUserAndPackage(ctx context.Context, userID, packageID primitive.ObjectID) (*SubscribePackage, error)
	FindByTnxHash(ctx context.Context, tnxHash string) (*SubscribePackage, error)
	FindByBinancePayID(ctx context.Context, binancePayID string) (*SubscribePackage, error)
	FindByStatus(ctx context.Context, status SubscriptionStatus) ([]*SubscribePackage, error)
	FindByPaymentMethod(ctx context.Context, paymentMethod PaymentMethod) ([]*SubscribePackage, error)
	FindActiveSubscriptions(ctx context.Context) ([]*SubscribePackage, error)
	FindExpiredSubscriptions(ctx context.Context) ([]*SubscribePackage, error)
	FindByIDTyped(ctx context.Context, id primitive.ObjectID) (*SubscribePackage, error)
	CreateSubscribePackage(ctx context.Context, subscription *SubscribePackage) (*mongo.InsertOneResult, error)
	UpdateSubscribePackage(ctx context.Context, id primitive.ObjectID, update *SubscribePackage) error
	DeleteSubscribePackage(ctx context.Context, id primitive.ObjectID) error
	DeleteByUserID(ctx context.Context, userID primitive.ObjectID) error
	FindAllByUser(ctx context.Context, userID primitive.ObjectID, skip, limit int64) ([]*SubscribePackage, error)
	CountSubscriptionsByUser(ctx context.Context, userID primitive.ObjectID) (int64, error)
	CountByStatus(ctx context.Context, status SubscriptionStatus) (int64, error)
	CountByPaymentMethod(ctx context.Context, paymentMethod PaymentMethod) (int64, error)
	UpdateStatus(ctx context.Context, id primitive.ObjectID, status SubscriptionStatus) error
	UpdatePaymentInfo(ctx context.Context, id primitive.ObjectID, tnxHash, binancePayID *string) error
	UpdateSubscriptionDates(ctx context.Context, id primitive.ObjectID, startDate, endDate *primitive.DateTime) error
}

// subscribePackageRepository implements SubscribePackageRepository interface
type subscribePackageRepository struct {
	BaseRepository
	collection *mongo.Collection
}

// NewSubscribePackageRepository creates a new subscribe package repository instance
func NewSubscribePackageRepository(collection *mongo.Collection) SubscribePackageRepository {
	return &subscribePackageRepository{
		BaseRepository: NewBaseRepository(collection),
		collection:     collection,
	}
}

// FindByUserID finds all subscriptions for a specific user
func (r *subscribePackageRepository) FindByUserID(ctx context.Context, userID primitive.ObjectID) ([]*SubscribePackage, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"user": userID})
	if err != nil {
		return nil, fmt.Errorf("failed to find subscriptions by user ID: %w", err)
	}
	defer cursor.Close(ctx)

	var subscriptions []*SubscribePackage
	for cursor.Next(ctx) {
		var subscription SubscribePackage
		if err := cursor.Decode(&subscription); err != nil {
			return nil, fmt.Errorf("failed to decode subscription: %w", err)
		}
		subscriptions = append(subscriptions, &subscription)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return subscriptions, nil
}

// FindByPackageID finds all subscriptions for a specific package
func (r *subscribePackageRepository) FindByPackageID(ctx context.Context, packageID primitive.ObjectID) ([]*SubscribePackage, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"package_id": packageID})
	if err != nil {
		return nil, fmt.Errorf("failed to find subscriptions by package ID: %w", err)
	}
	defer cursor.Close(ctx)

	var subscriptions []*SubscribePackage
	for cursor.Next(ctx) {
		var subscription SubscribePackage
		if err := cursor.Decode(&subscription); err != nil {
			return nil, fmt.Errorf("failed to decode subscription: %w", err)
		}
		subscriptions = append(subscriptions, &subscription)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return subscriptions, nil
}

// FindByUserAndPackage finds a subscription by user and package
func (r *subscribePackageRepository) FindByUserAndPackage(ctx context.Context, userID, packageID primitive.ObjectID) (*SubscribePackage, error) {
	var subscription SubscribePackage
	err := r.collection.FindOne(ctx, bson.M{"user": userID, "package_id": packageID}).Decode(&subscription)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("subscription not found for user ID: %s and package ID: %s", userID.Hex(), packageID.Hex())
		}
		return nil, fmt.Errorf("failed to find subscription by user and package: %w", err)
	}

	return &subscription, nil
}

// FindByTnxHash finds a subscription by transaction hash
func (r *subscribePackageRepository) FindByTnxHash(ctx context.Context, tnxHash string) (*SubscribePackage, error) {
	var subscription SubscribePackage
	err := r.collection.FindOne(ctx, bson.M{"tnx_hash": tnxHash}).Decode(&subscription)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("subscription not found with transaction hash: %s", tnxHash)
		}
		return nil, fmt.Errorf("failed to find subscription by transaction hash: %w", err)
	}

	return &subscription, nil
}

// FindByBinancePayID finds a subscription by Binance Pay ID
func (r *subscribePackageRepository) FindByBinancePayID(ctx context.Context, binancePayID string) (*SubscribePackage, error) {
	var subscription SubscribePackage
	err := r.collection.FindOne(ctx, bson.M{"binance_pay_id": binancePayID}).Decode(&subscription)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("subscription not found with Binance Pay ID: %s", binancePayID)
		}
		return nil, fmt.Errorf("failed to find subscription by Binance Pay ID: %w", err)
	}

	return &subscription, nil
}

// FindByStatus finds all subscriptions with a specific status
func (r *subscribePackageRepository) FindByStatus(ctx context.Context, status SubscriptionStatus) ([]*SubscribePackage, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"status": status})
	if err != nil {
		return nil, fmt.Errorf("failed to find subscriptions by status: %w", err)
	}
	defer cursor.Close(ctx)

	var subscriptions []*SubscribePackage
	for cursor.Next(ctx) {
		var subscription SubscribePackage
		if err := cursor.Decode(&subscription); err != nil {
			return nil, fmt.Errorf("failed to decode subscription: %w", err)
		}
		subscriptions = append(subscriptions, &subscription)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return subscriptions, nil
}

// FindByPaymentMethod finds all subscriptions with a specific payment method
func (r *subscribePackageRepository) FindByPaymentMethod(ctx context.Context, paymentMethod PaymentMethod) ([]*SubscribePackage, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"payment_method": paymentMethod})
	if err != nil {
		return nil, fmt.Errorf("failed to find subscriptions by payment method: %w", err)
	}
	defer cursor.Close(ctx)

	var subscriptions []*SubscribePackage
	for cursor.Next(ctx) {
		var subscription SubscribePackage
		if err := cursor.Decode(&subscription); err != nil {
			return nil, fmt.Errorf("failed to decode subscription: %w", err)
		}
		subscriptions = append(subscriptions, &subscription)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return subscriptions, nil
}

// FindActiveSubscriptions finds all active subscriptions (completed status with valid dates)
func (r *subscribePackageRepository) FindActiveSubscriptions(ctx context.Context) ([]*SubscribePackage, error) {
	filter := bson.M{
		"status": SubscriptionStatusCompleted,
		"end_date": bson.M{
			"$gt": time.Now(),
		},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find active subscriptions: %w", err)
	}
	defer cursor.Close(ctx)

	var subscriptions []*SubscribePackage
	for cursor.Next(ctx) {
		var subscription SubscribePackage
		if err := cursor.Decode(&subscription); err != nil {
			return nil, fmt.Errorf("failed to decode subscription: %w", err)
		}
		subscriptions = append(subscriptions, &subscription)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return subscriptions, nil
}

// FindExpiredSubscriptions finds all expired subscriptions
func (r *subscribePackageRepository) FindExpiredSubscriptions(ctx context.Context) ([]*SubscribePackage, error) {
	filter := bson.M{
		"status": SubscriptionStatusCompleted,
		"end_date": bson.M{
			"$lt": time.Now(),
		},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find expired subscriptions: %w", err)
	}
	defer cursor.Close(ctx)

	var subscriptions []*SubscribePackage
	for cursor.Next(ctx) {
		var subscription SubscribePackage
		if err := cursor.Decode(&subscription); err != nil {
			return nil, fmt.Errorf("failed to decode subscription: %w", err)
		}
		subscriptions = append(subscriptions, &subscription)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return subscriptions, nil
}

// FindByIDTyped finds a subscription by ID and returns typed SubscribePackage struct
func (r *subscribePackageRepository) FindByIDTyped(ctx context.Context, id primitive.ObjectID) (*SubscribePackage, error) {
	var subscription SubscribePackage
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&subscription)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("subscription not found with ID: %s", id.Hex())
		}
		return nil, fmt.Errorf("failed to find subscription by ID: %w", err)
	}

	return &subscription, nil
}

// CreateSubscribePackage creates a new subscription
func (r *subscribePackageRepository) CreateSubscribePackage(ctx context.Context, subscription *SubscribePackage) (*mongo.InsertOneResult, error) {
	result, err := r.collection.InsertOne(ctx, subscription)
	if err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	return result, nil
}

// UpdateSubscribePackage updates an existing subscription
func (r *subscribePackageRepository) UpdateSubscribePackage(ctx context.Context, id primitive.ObjectID, update *SubscribePackage) error {
	update.UpdatedAt = time.Now()
	
	_, err := r.collection.UpdateByID(ctx, id, bson.M{"$set": update})
	if err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	return nil
}

// DeleteSubscribePackage deletes a subscription by ID
func (r *subscribePackageRepository) DeleteSubscribePackage(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	return nil
}

// DeleteByUserID deletes all subscriptions for a specific user
func (r *subscribePackageRepository) DeleteByUserID(ctx context.Context, userID primitive.ObjectID) error {
	_, err := r.collection.DeleteMany(ctx, bson.M{"user": userID})
	if err != nil {
		return fmt.Errorf("failed to delete subscriptions by user ID: %w", err)
	}

	return nil
}

// FindAllByUser retrieves all subscriptions for a user with pagination
func (r *subscribePackageRepository) FindAllByUser(ctx context.Context, userID primitive.ObjectID, skip, limit int64) ([]*SubscribePackage, error) {
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
		return nil, fmt.Errorf("failed to find subscriptions by user: %w", err)
	}
	defer cursor.Close(ctx)

	var subscriptions []*SubscribePackage
	for cursor.Next(ctx) {
		var subscription SubscribePackage
		if err := cursor.Decode(&subscription); err != nil {
			return nil, fmt.Errorf("failed to decode subscription: %w", err)
		}
		subscriptions = append(subscriptions, &subscription)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return subscriptions, nil
}

// CountSubscriptionsByUser returns the number of subscriptions for a specific user
func (r *subscribePackageRepository) CountSubscriptionsByUser(ctx context.Context, userID primitive.ObjectID) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"user": userID})
	if err != nil {
		return 0, fmt.Errorf("failed to count subscriptions by user: %w", err)
	}

	return count, nil
}

// CountByStatus returns the number of subscriptions with a specific status
func (r *subscribePackageRepository) CountByStatus(ctx context.Context, status SubscriptionStatus) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"status": status})
	if err != nil {
		return 0, fmt.Errorf("failed to count subscriptions by status: %w", err)
	}

	return count, nil
}

// CountByPaymentMethod returns the number of subscriptions with a specific payment method
func (r *subscribePackageRepository) CountByPaymentMethod(ctx context.Context, paymentMethod PaymentMethod) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"payment_method": paymentMethod})
	if err != nil {
		return 0, fmt.Errorf("failed to count subscriptions by payment method: %w", err)
	}

	return count, nil
}

// UpdateStatus updates only the status of a subscription
func (r *subscribePackageRepository) UpdateStatus(ctx context.Context, id primitive.ObjectID, status SubscriptionStatus) error {
	update := bson.M{
		"status":     status,
		"updated_at": bson.M{"$currentDate": true},
	}

	_, err := r.collection.UpdateByID(ctx, id, bson.M{"$set": update})
	if err != nil {
		return fmt.Errorf("failed to update subscription status: %w", err)
	}

	return nil
}

// UpdatePaymentInfo updates payment-related information
func (r *subscribePackageRepository) UpdatePaymentInfo(ctx context.Context, id primitive.ObjectID, tnxHash, binancePayID *string) error {
	update := bson.M{
		"updated_at": bson.M{"$currentDate": true},
	}

	if tnxHash != nil {
		update["tnx_hash"] = tnxHash
	}
	if binancePayID != nil {
		update["binance_pay_id"] = binancePayID
	}

	_, err := r.collection.UpdateByID(ctx, id, bson.M{"$set": update})
	if err != nil {
		return fmt.Errorf("failed to update subscription payment info: %w", err)
	}

	return nil
}

// UpdateSubscriptionDates updates start and end dates of a subscription
func (r *subscribePackageRepository) UpdateSubscriptionDates(ctx context.Context, id primitive.ObjectID, startDate, endDate *primitive.DateTime) error {
	update := bson.M{
		"updated_at": bson.M{"$currentDate": true},
	}

	if startDate != nil {
		update["start_date"] = startDate
	}
	if endDate != nil {
		update["end_date"] = endDate
	}

	_, err := r.collection.UpdateByID(ctx, id, bson.M{"$set": update})
	if err != nil {
		return fmt.Errorf("failed to update subscription dates: %w", err)
	}

	return nil
}
