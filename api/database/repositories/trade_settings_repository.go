package repositories

import (
	"context"
	"fmt"

	"copier/internal/database/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TradeSettingsRepository defines trade settings-specific repository operations
type TradeSettingsRepository interface {
	BaseRepository
	FindByUserID(ctx context.Context, userID uuid.UUID) (*models.TradeSettings, error)
	FindByIDTyped(ctx context.Context, id uuid.UUID) (*models.TradeSettings, error)
	CreateTradeSettings(ctx context.Context, settings *models.TradeSettings) error
	UpdateTradeSettings(ctx context.Context, id uuid.UUID, update *models.TradeSettings) error
	UpdateByUserID(ctx context.Context, userID uuid.UUID, update *models.TradeSettings) error
	DeleteTradeSettings(ctx context.Context, id uuid.UUID) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
	FindAllWithStopLoss(ctx context.Context, skip, limit int) ([]*models.TradeSettings, error)
	FindAllWithTakeProfit(ctx context.Context, skip, limit int) ([]*models.TradeSettings, error)
	CountTradeSettings(ctx context.Context) (int64, error)
	UpdateStopLossSettings(ctx context.Context, userID uuid.UUID, percentage int, status bool) error
	UpdateTakeProfitSettings(ctx context.Context, userID uuid.UUID, status bool, step int, percentages []float64) error
}

// tradeSettingsRepository implements TradeSettingsRepository interface
type tradeSettingsRepository struct {
	BaseRepository
	db *gorm.DB
}

// NewTradeSettingsRepository creates a new trade settings repository instance
func NewTradeSettingsRepository(db *gorm.DB) TradeSettingsRepository {
	return &tradeSettingsRepository{
		BaseRepository: NewBaseRepository(db),
		db:             db,
	}
}

// FindByUserID finds trade settings for a specific user
func (r *tradeSettingsRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (*models.TradeSettings, error) {
	var settings models.TradeSettings
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&settings).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("trade settings not found for user ID: %s", userID)
		}
		return nil, fmt.Errorf("failed to find trade settings by user ID: %w", err)
	}

	return &settings, nil
}

// FindByIDTyped finds trade settings by ID and returns typed TradeSettings struct
func (r *tradeSettingsRepository) FindByIDTyped(ctx context.Context, id uuid.UUID) (*models.TradeSettings, error) {
	var settings models.TradeSettings
	err := r.db.WithContext(ctx).First(&settings, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("trade settings not found with ID: %s", id)
		}
		return nil, fmt.Errorf("failed to find trade settings by ID: %w", err)
	}

	return &settings, nil
}

// CreateTradeSettings creates new trade settings
func (r *tradeSettingsRepository) CreateTradeSettings(ctx context.Context, settings *models.TradeSettings) error {
	err := r.db.WithContext(ctx).Create(settings).Error
	if err != nil {
		return fmt.Errorf("failed to create trade settings: %w", err)
	}

	return nil
}

// UpdateTradeSettings updates existing trade settings by ID
func (r *tradeSettingsRepository) UpdateTradeSettings(ctx context.Context, id uuid.UUID, update *models.TradeSettings) error {
	err := r.db.WithContext(ctx).Model(&models.TradeSettings{}).Where("id = ?", id).Updates(update).Error
	if err != nil {
		return fmt.Errorf("failed to update trade settings: %w", err)
	}

	return nil
}

// UpdateByUserID updates trade settings by user ID
func (r *tradeSettingsRepository) UpdateByUserID(ctx context.Context, userID uuid.UUID, update *models.TradeSettings) error {
	err := r.db.WithContext(ctx).Model(&models.TradeSettings{}).Where("user_id = ?", userID).Updates(update).Error
	if err != nil {
		return fmt.Errorf("failed to update trade settings by user ID: %w", err)
	}

	return nil
}

// DeleteTradeSettings deletes trade settings by ID
func (r *tradeSettingsRepository) DeleteTradeSettings(ctx context.Context, id uuid.UUID) error {
	err := r.db.WithContext(ctx).Delete(&models.TradeSettings{}, id).Error
	if err != nil {
		return fmt.Errorf("failed to delete trade settings: %w", err)
	}

	return nil
}

// DeleteByUserID deletes trade settings by user ID
func (r *tradeSettingsRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.TradeSettings{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete trade settings by user ID: %w", err)
	}

	return nil
}

// FindAllWithStopLoss finds all trade settings with stop loss enabled
func (r *tradeSettingsRepository) FindAllWithStopLoss(ctx context.Context, skip, limit int) ([]*models.TradeSettings, error) {
	var settings []*models.TradeSettings
	query := r.db.WithContext(ctx).Where("stop_loss_status = ?", true).Order("created_at desc")
	if skip > 0 {
		query = query.Offset(skip)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&settings).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find trade settings with stop loss: %w", err)
	}

	return settings, nil
}

// FindAllWithTakeProfit finds all trade settings with take profit enabled
func (r *tradeSettingsRepository) FindAllWithTakeProfit(ctx context.Context, skip, limit int) ([]*models.TradeSettings, error) {
	var settings []*models.TradeSettings
	query := r.db.WithContext(ctx).Where("take_profit_status = ?", true).Order("created_at desc")
	if skip > 0 {
		query = query.Offset(skip)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&settings).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find trade settings with take profit: %w", err)
	}

	return settings, nil
}

// CountTradeSettings returns the total number of trade settings
func (r *tradeSettingsRepository) CountTradeSettings(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.TradeSettings{}).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count trade settings: %w", err)
	}

	return count, nil
}

// UpdateStopLossSettings updates only stop loss related settings
func (r *tradeSettingsRepository) UpdateStopLossSettings(ctx context.Context, userID uuid.UUID, percentage int, status bool) error {
	err := r.db.WithContext(ctx).Model(&models.TradeSettings{}).Where("user_id = ?", userID).Updates(map[string]interface{}{
		"stop_loss_percentage": percentage,
		"stop_loss_status":     status,
	}).Error
	if err != nil {
		return fmt.Errorf("failed to update stop loss settings: %w", err)
	}

	return nil
}

// UpdateTakeProfitSettings updates only take profit related settings
func (r *tradeSettingsRepository) UpdateTakeProfitSettings(ctx context.Context, userID uuid.UUID, status bool, step int, percentages []float64) error {
	err := r.db.WithContext(ctx).Model(&models.TradeSettings{}).Where("user_id = ?", userID).Updates(map[string]interface{}{
		"take_profit_status": status,
		"take_profit_step":   step,
		"tp_percentage":      percentages,
	}).Error
	if err != nil {
		return fmt.Errorf("failed to update take profit settings: %w", err)
	}

	return nil
}
