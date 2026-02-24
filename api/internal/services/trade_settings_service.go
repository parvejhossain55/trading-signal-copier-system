package services

import (
	"context"

	"copier/database/repositories"
	"copier/internal/database/models"

	"github.com/google/uuid"
)

// TradeSettingsService defines trade settings business logic operations
type TradeSettingsService interface {
	GetTradeSettingsByUser(ctx context.Context, userID uuid.UUID) (*models.TradeSettings, error)
	UpsertTradeSettings(ctx context.Context, userID uuid.UUID, settings *models.TradeSettings) (*models.TradeSettings, error)
	UpdateStopLoss(ctx context.Context, userID uuid.UUID, percentage int, status bool) error
	UpdateTakeProfit(ctx context.Context, userID uuid.UUID, status bool, step int, percentages []float64) error
}

type tradeSettingsService struct {
	settingsRepo repositories.TradeSettingsRepository
}

// NewTradeSettingsService creates a new trade settings service instance
func NewTradeSettingsService(settingsRepo repositories.TradeSettingsRepository) TradeSettingsService {
	return &tradeSettingsService{
		settingsRepo: settingsRepo,
	}
}

func (s *tradeSettingsService) GetTradeSettingsByUser(ctx context.Context, userID uuid.UUID) (*models.TradeSettings, error) {
	return s.settingsRepo.FindByUserID(ctx, userID)
}

func (s *tradeSettingsService) UpsertTradeSettings(ctx context.Context, userID uuid.UUID, settings *models.TradeSettings) (*models.TradeSettings, error) {
	_, err := s.settingsRepo.FindByUserID(ctx, userID)
	if err != nil {
		// New settings
		settings.UserID = userID
		err = s.settingsRepo.CreateTradeSettings(ctx, settings)
		if err != nil {
			return nil, err
		}
		return settings, nil
	}

	// Update existing
	err = s.settingsRepo.UpdateByUserID(ctx, userID, settings)
	if err != nil {
		return nil, err
	}
	return s.settingsRepo.FindByUserID(ctx, userID)
}

func (s *tradeSettingsService) UpdateStopLoss(ctx context.Context, userID uuid.UUID, percentage int, status bool) error {
	return s.settingsRepo.UpdateStopLossSettings(ctx, userID, percentage, status)
}

func (s *tradeSettingsService) UpdateTakeProfit(ctx context.Context, userID uuid.UUID, status bool, step int, percentages []float64) error {
	return s.settingsRepo.UpdateTakeProfitSettings(ctx, userID, status, step, percentages)
}
