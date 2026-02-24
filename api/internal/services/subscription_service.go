package services

import (
	"context"
	"fmt"
	"time"

	"copier/database/repositories"
	"copier/internal/database/models"

	"github.com/google/uuid"
)

// SubscriptionService defines subscription business logic operations
type SubscriptionService interface {
	Subscribe(ctx context.Context, userID, packageID uuid.UUID, paymentMethod models.PaymentMethod) (*models.SubscribePackage, error)
	VerifyPayment(ctx context.Context, id uuid.UUID, tnxHash *string, binancePayID *string) (*models.SubscribePackage, error)
	GetUserSubscriptions(ctx context.Context, userID uuid.UUID, skip, limit int) ([]*models.SubscribePackage, int64, error)
	GetSubscriptionByID(ctx context.Context, id uuid.UUID) (*models.SubscribePackage, error)
	CancelSubscription(ctx context.Context, id uuid.UUID) error
}

type subscriptionService struct {
	subRepo     repositories.SubscribePackageRepository
	packageRepo repositories.PackageRepository
}

// NewSubscriptionService creates a new subscription service instance
func NewSubscriptionService(subRepo repositories.SubscribePackageRepository, packageRepo repositories.PackageRepository) SubscriptionService {
	return &subscriptionService{
		subRepo:     subRepo,
		packageRepo: packageRepo,
	}
}

// Subscribe handles the initiation of a new subscription
func (s *subscriptionService) Subscribe(ctx context.Context, userID, packageID uuid.UUID, paymentMethod models.PaymentMethod) (*models.SubscribePackage, error) {
	// Find the package to get the current price
	pkg, err := s.packageRepo.FindByIDTyped(ctx, packageID)
	if err != nil {
		return nil, fmt.Errorf("invalid package: %w", err)
	}

	subscription := &models.SubscribePackage{
		UserID:        userID,
		PackageID:     packageID,
		Price:         pkg.Price,
		PaymentMethod: paymentMethod,
		Status:        models.SubscriptionStatusPending,
	}

	err = s.subRepo.CreateSubscribePackage(ctx, subscription)
	if err != nil {
		return nil, err
	}

	return subscription, nil
}

// VerifyPayment handles the verification and completion of a subscription
func (s *subscriptionService) VerifyPayment(ctx context.Context, id uuid.UUID, tnxHash *string, binancePayID *string) (*models.SubscribePackage, error) {
	// In a real app, this would involve calling external payment APIs (Binance, Stripe, etc.)
	// For this implementation, we'll mark it as completed if a transaction identifier is provided.

	sub, err := s.subRepo.FindByIDTyped(ctx, id)
	if err != nil {
		return nil, err
	}

	if tnxHash == nil && binancePayID == nil {
		return nil, fmt.Errorf("transaction hash or Binance Pay ID is required for verification")
	}

	now := time.Now()
	// Default duration is 30 days if not specified in the package
	durationDays := 30
	pkg, err := s.packageRepo.FindByIDTyped(ctx, sub.PackageID)
	if err == nil && pkg.DurationDays != nil {
		durationDays = *pkg.DurationDays
	}

	endDate := now.AddDate(0, 0, durationDays)

	err = s.subRepo.UpdateStatus(ctx, id, models.SubscriptionStatusCompleted)
	if err != nil {
		return nil, err
	}

	err = s.subRepo.UpdatePaymentInfo(ctx, id, tnxHash, binancePayID)
	if err != nil {
		return nil, err
	}

	err = s.subRepo.UpdateSubscriptionDates(ctx, id, &now, &endDate)
	if err != nil {
		return nil, err
	}

	return s.subRepo.FindByIDTyped(ctx, id)
}

// GetUserSubscriptions retrieves all subscriptions for a specific user
func (s *subscriptionService) GetUserSubscriptions(ctx context.Context, userID uuid.UUID, skip, limit int) ([]*models.SubscribePackage, int64, error) {
	subs, err := s.subRepo.FindAllByUser(ctx, userID, skip, limit)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.subRepo.CountSubscriptionsByUser(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	return subs, total, nil
}

// GetSubscriptionByID retrieves a single subscription record
func (s *subscriptionService) GetSubscriptionByID(ctx context.Context, id uuid.UUID) (*models.SubscribePackage, error) {
	return s.subRepo.FindByIDTyped(ctx, id)
}

// CancelSubscription marks a subscription as cancelled
func (s *subscriptionService) CancelSubscription(ctx context.Context, id uuid.UUID) error {
	return s.subRepo.UpdateStatus(ctx, id, models.SubscriptionStatusCancelled)
}
