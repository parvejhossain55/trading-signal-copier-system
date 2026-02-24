package repositories

import (
	"context"
	"fmt"
	"time"

	"copier/internal/database/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SubscribePackageRepository defines subscribe package-specific repository operations
type SubscribePackageRepository interface {
	BaseRepository
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*models.SubscribePackage, error)
	FindByPackageID(ctx context.Context, packageID uuid.UUID) ([]*models.SubscribePackage, error)
	FindByUserAndPackage(ctx context.Context, userID, packageID uuid.UUID) (*models.SubscribePackage, error)
	FindByTnxHash(ctx context.Context, tnxHash string) (*models.SubscribePackage, error)
	FindByBinancePayID(ctx context.Context, binancePayID string) (*models.SubscribePackage, error)
	FindByStatus(ctx context.Context, status models.SubscriptionStatus) ([]*models.SubscribePackage, error)
	FindByPaymentMethod(ctx context.Context, paymentMethod models.PaymentMethod) ([]*models.SubscribePackage, error)
	FindActiveSubscriptions(ctx context.Context) ([]*models.SubscribePackage, error)
	FindExpiredSubscriptions(ctx context.Context) ([]*models.SubscribePackage, error)
	FindByIDTyped(ctx context.Context, id uuid.UUID) (*models.SubscribePackage, error)
	CreateSubscribePackage(ctx context.Context, subscription *models.SubscribePackage) error
	UpdateSubscribePackage(ctx context.Context, id uuid.UUID, update *models.SubscribePackage) error
	DeleteSubscribePackage(ctx context.Context, id uuid.UUID) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
	FindAllByUser(ctx context.Context, userID uuid.UUID, skip, limit int) ([]*models.SubscribePackage, error)
	CountSubscriptionsByUser(ctx context.Context, userID uuid.UUID) (int64, error)
	CountByStatus(ctx context.Context, status models.SubscriptionStatus) (int64, error)
	CountByPaymentMethod(ctx context.Context, paymentMethod models.PaymentMethod) (int64, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.SubscriptionStatus) error
	UpdatePaymentInfo(ctx context.Context, id uuid.UUID, tnxHash, binancePayID *string) error
	UpdateSubscriptionDates(ctx context.Context, id uuid.UUID, startDate, endDate *time.Time) error
}

// subscribePackageRepository implements SubscribePackageRepository interface
type subscribePackageRepository struct {
	BaseRepository
	db *gorm.DB
}

// NewSubscribePackageRepository creates a new subscribe package repository instance
func NewSubscribePackageRepository(db *gorm.DB) SubscribePackageRepository {
	return &subscribePackageRepository{
		BaseRepository: NewBaseRepository(db),
		db:             db,
	}
}

// FindByUserID finds all subscriptions for a specific user
func (r *subscribePackageRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*models.SubscribePackage, error) {
	var subscriptions []*models.SubscribePackage
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&subscriptions).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find subscriptions by user ID: %w", err)
	}

	return subscriptions, nil
}

// FindByPackageID finds all subscriptions for a specific package
func (r *subscribePackageRepository) FindByPackageID(ctx context.Context, packageID uuid.UUID) ([]*models.SubscribePackage, error) {
	var subscriptions []*models.SubscribePackage
	err := r.db.WithContext(ctx).Where("package_id = ?", packageID).Find(&subscriptions).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find subscriptions by package ID: %w", err)
	}

	return subscriptions, nil
}

// FindByUserAndPackage finds a subscription by user and package
func (r *subscribePackageRepository) FindByUserAndPackage(ctx context.Context, userID, packageID uuid.UUID) (*models.SubscribePackage, error) {
	var subscription models.SubscribePackage
	err := r.db.WithContext(ctx).Where("user_id = ? AND package_id = ?", userID, packageID).First(&subscription).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("subscription not found for user ID: %s and package ID: %s", userID, packageID)
		}
		return nil, fmt.Errorf("failed to find subscription by user and package: %w", err)
	}

	return &subscription, nil
}

// FindByTnxHash finds a subscription by transaction hash
func (r *subscribePackageRepository) FindByTnxHash(ctx context.Context, tnxHash string) (*models.SubscribePackage, error) {
	var subscription models.SubscribePackage
	err := r.db.WithContext(ctx).Where("tnx_hash = ?", tnxHash).First(&subscription).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("subscription not found with transaction hash: %s", tnxHash)
		}
		return nil, fmt.Errorf("failed to find subscription by transaction hash: %w", err)
	}

	return &subscription, nil
}

// FindByBinancePayID finds a subscription by Binance Pay ID
func (r *subscribePackageRepository) FindByBinancePayID(ctx context.Context, binancePayID string) (*models.SubscribePackage, error) {
	var subscription models.SubscribePackage
	err := r.db.WithContext(ctx).Where("binance_pay_id = ?", binancePayID).First(&subscription).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("subscription not found with Binance Pay ID: %s", binancePayID)
		}
		return nil, fmt.Errorf("failed to find subscription by Binance Pay ID: %w", err)
	}

	return &subscription, nil
}

// FindByStatus finds all subscriptions with a specific status
func (r *subscribePackageRepository) FindByStatus(ctx context.Context, status models.SubscriptionStatus) ([]*models.SubscribePackage, error) {
	var subscriptions []*models.SubscribePackage
	err := r.db.WithContext(ctx).Where("status = ?", status).Find(&subscriptions).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find subscriptions by status: %w", err)
	}

	return subscriptions, nil
}

// FindByPaymentMethod finds all subscriptions with a specific payment method
func (r *subscribePackageRepository) FindByPaymentMethod(ctx context.Context, paymentMethod models.PaymentMethod) ([]*models.SubscribePackage, error) {
	var subscriptions []*models.SubscribePackage
	err := r.db.WithContext(ctx).Where("payment_method = ?", paymentMethod).Find(&subscriptions).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find subscriptions by payment method: %w", err)
	}

	return subscriptions, nil
}

// FindActiveSubscriptions finds all active subscriptions (completed status with valid dates)
func (r *subscribePackageRepository) FindActiveSubscriptions(ctx context.Context) ([]*models.SubscribePackage, error) {
	var subscriptions []*models.SubscribePackage
	err := r.db.WithContext(ctx).Where("status = ? AND end_date > ?", models.SubscriptionStatusCompleted, time.Now()).Find(&subscriptions).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find active subscriptions: %w", err)
	}

	return subscriptions, nil
}

// FindExpiredSubscriptions finds all expired subscriptions
func (r *subscribePackageRepository) FindExpiredSubscriptions(ctx context.Context) ([]*models.SubscribePackage, error) {
	var subscriptions []*models.SubscribePackage
	err := r.db.WithContext(ctx).Where("status = ? AND end_date < ?", models.SubscriptionStatusCompleted, time.Now()).Find(&subscriptions).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find expired subscriptions: %w", err)
	}

	return subscriptions, nil
}

// FindByIDTyped finds a subscription by ID and returns typed SubscribePackage struct
func (r *subscribePackageRepository) FindByIDTyped(ctx context.Context, id uuid.UUID) (*models.SubscribePackage, error) {
	var subscription models.SubscribePackage
	err := r.db.WithContext(ctx).First(&subscription, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("subscription not found with ID: %s", id)
		}
		return nil, fmt.Errorf("failed to find subscription by ID: %w", err)
	}

	return &subscription, nil
}

// CreateSubscribePackage creates a new subscription
func (r *subscribePackageRepository) CreateSubscribePackage(ctx context.Context, subscription *models.SubscribePackage) error {
	err := r.db.WithContext(ctx).Create(subscription).Error
	if err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	return nil
}

// UpdateSubscribePackage updates an existing subscription
func (r *subscribePackageRepository) UpdateSubscribePackage(ctx context.Context, id uuid.UUID, update *models.SubscribePackage) error {
	err := r.db.WithContext(ctx).Model(&models.SubscribePackage{}).Where("id = ?", id).Updates(update).Error
	if err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	return nil
}

// DeleteSubscribePackage deletes a subscription by ID
func (r *subscribePackageRepository) DeleteSubscribePackage(ctx context.Context, id uuid.UUID) error {
	err := r.db.WithContext(ctx).Delete(&models.SubscribePackage{}, id).Error
	if err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	return nil
}

// DeleteByUserID deletes all subscriptions for a specific user
func (r *subscribePackageRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.SubscribePackage{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete subscriptions by user ID: %w", err)
	}

	return nil
}

// FindAllByUser retrieves all subscriptions for a user with pagination
func (r *subscribePackageRepository) FindAllByUser(ctx context.Context, userID uuid.UUID, skip, limit int) ([]*models.SubscribePackage, error) {
	var subscriptions []*models.SubscribePackage
	query := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at desc")
	if skip > 0 {
		query = query.Offset(skip)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&subscriptions).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find subscriptions by user: %w", err)
	}

	return subscriptions, nil
}

// CountSubscriptionsByUser returns the number of subscriptions for a specific user
func (r *subscribePackageRepository) CountSubscriptionsByUser(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.SubscribePackage{}).Where("user_id = ?", userID).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count subscriptions by user: %w", err)
	}

	return count, nil
}

// CountByStatus returns the number of subscriptions with a specific status
func (r *subscribePackageRepository) CountByStatus(ctx context.Context, status models.SubscriptionStatus) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.SubscribePackage{}).Where("status = ?", status).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count subscriptions by status: %w", err)
	}

	return count, nil
}

// CountByPaymentMethod returns the number of subscriptions with a specific payment method
func (r *subscribePackageRepository) CountByPaymentMethod(ctx context.Context, paymentMethod models.PaymentMethod) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.SubscribePackage{}).Where("payment_method = ?", paymentMethod).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count subscriptions by payment method: %w", err)
	}

	return count, nil
}

// UpdateStatus updates only the status of a subscription
func (r *subscribePackageRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.SubscriptionStatus) error {
	err := r.db.WithContext(ctx).Model(&models.SubscribePackage{}).Where("id = ?", id).Update("status", status).Error
	if err != nil {
		return fmt.Errorf("failed to update subscription status: %w", err)
	}

	return nil
}

// UpdatePaymentInfo updates payment-related information
func (r *subscribePackageRepository) UpdatePaymentInfo(ctx context.Context, id uuid.UUID, tnxHash, binancePayID *string) error {
	updates := make(map[string]interface{})
	if tnxHash != nil {
		updates["tnx_hash"] = tnxHash
	}
	if binancePayID != nil {
		updates["binance_pay_id"] = binancePayID
	}

	if len(updates) == 0 {
		return nil
	}

	err := r.db.WithContext(ctx).Model(&models.SubscribePackage{}).Where("id = ?", id).Updates(updates).Error
	if err != nil {
		return fmt.Errorf("failed to update subscription payment info: %w", err)
	}

	return nil
}

// UpdateSubscriptionDates updates start and end dates of a subscription
func (r *subscribePackageRepository) UpdateSubscriptionDates(ctx context.Context, id uuid.UUID, startDate, endDate *time.Time) error {
	updates := make(map[string]interface{})
	if startDate != nil {
		updates["start_date"] = startDate
	}
	if endDate != nil {
		updates["end_date"] = endDate
	}

	if len(updates) == 0 {
		return nil
	}

	err := r.db.WithContext(ctx).Model(&models.SubscribePackage{}).Where("id = ?", id).Updates(updates).Error
	if err != nil {
		return fmt.Errorf("failed to update subscription dates: %w", err)
	}

	return nil
}
