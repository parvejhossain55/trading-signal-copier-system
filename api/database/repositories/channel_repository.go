package repositories

import (
	"context"
	"fmt"

	"copier/internal/database/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ChannelRepository defines channel-specific repository operations
type ChannelRepository interface {
	BaseRepository
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Channel, error)
	FindByChannelID(ctx context.Context, channelID string) (*models.Channel, error)
	FindByUserAndName(ctx context.Context, userID uuid.UUID, name string) (*models.Channel, error)
	FindByIDTyped(ctx context.Context, id uuid.UUID) (*models.Channel, error)
	CreateChannel(ctx context.Context, channel *models.Channel) error
	UpdateChannel(ctx context.Context, id uuid.UUID, update *models.Channel) error
	DeleteChannel(ctx context.Context, id uuid.UUID) error
	DeleteChannelsByUser(ctx context.Context, userID uuid.UUID) error
	FindAllByUser(ctx context.Context, userID uuid.UUID, skip, limit int) ([]*models.Channel, error)
	CountChannelsByUser(ctx context.Context, userID uuid.UUID) (int64, error)
}

// channelRepository implements ChannelRepository interface
type channelRepository struct {
	BaseRepository
	db *gorm.DB
}

// NewChannelRepository creates a new channel repository instance
func NewChannelRepository(db *gorm.DB) ChannelRepository {
	return &channelRepository{
		BaseRepository: NewBaseRepository(db),
		db:             db,
	}
}

// FindByUserID finds all channels for a specific user
func (r *channelRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Channel, error) {
	var channels []*models.Channel
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&channels).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find channels by user ID: %w", err)
	}

	return channels, nil
}

// FindByChannelID finds a channel by its channel ID
func (r *channelRepository) FindByChannelID(ctx context.Context, channelID string) (*models.Channel, error) {
	var channel models.Channel
	err := r.db.WithContext(ctx).Where("channel_id = ?", channelID).First(&channel).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("channel not found with ID: %s", channelID)
		}
		return nil, fmt.Errorf("failed to find channel by channel ID: %w", err)
	}

	return &channel, nil
}

// FindByUserAndName finds a channel by user ID and name
func (r *channelRepository) FindByUserAndName(ctx context.Context, userID uuid.UUID, name string) (*models.Channel, error) {
	var channel models.Channel
	err := r.db.WithContext(ctx).Where("user_id = ? AND name = ?", userID, name).First(&channel).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("channel not found with user ID: %s and name: %s", userID, name)
		}
		return nil, fmt.Errorf("failed to find channel by user and name: %w", err)
	}

	return &channel, nil
}

// FindByIDTyped finds a channel by ID and returns typed Channel struct
func (r *channelRepository) FindByIDTyped(ctx context.Context, id uuid.UUID) (*models.Channel, error) {
	var channel models.Channel
	err := r.db.WithContext(ctx).First(&channel, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("channel not found with ID: %s", id)
		}
		return nil, fmt.Errorf("failed to find channel by ID: %w", err)
	}

	return &channel, nil
}

// CreateChannel creates a new channel
func (r *channelRepository) CreateChannel(ctx context.Context, channel *models.Channel) error {
	err := r.db.WithContext(ctx).Create(channel).Error
	if err != nil {
		return fmt.Errorf("failed to create channel: %w", err)
	}

	return nil
}

// UpdateChannel updates an existing channel
func (r *channelRepository) UpdateChannel(ctx context.Context, id uuid.UUID, update *models.Channel) error {
	err := r.db.WithContext(ctx).Model(&models.Channel{}).Where("id = ?", id).Updates(update).Error
	if err != nil {
		return fmt.Errorf("failed to update channel: %w", err)
	}

	return nil
}

// DeleteChannel deletes a channel by ID
func (r *channelRepository) DeleteChannel(ctx context.Context, id uuid.UUID) error {
	err := r.db.WithContext(ctx).Delete(&models.Channel{}, id).Error
	if err != nil {
		return fmt.Errorf("failed to delete channel: %w", err)
	}

	return nil
}

// DeleteChannelsByUser deletes all channels for a specific user
func (r *channelRepository) DeleteChannelsByUser(ctx context.Context, userID uuid.UUID) error {
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.Channel{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete channels by user: %w", err)
	}

	return nil
}

// FindAllByUser retrieves all channels for a user with pagination
func (r *channelRepository) FindAllByUser(ctx context.Context, userID uuid.UUID, skip, limit int) ([]*models.Channel, error) {
	var channels []*models.Channel
	query := r.db.WithContext(ctx).Where("user_id = ? ", userID).Order("created_at desc")
	if skip > 0 {
		query = query.Offset(skip)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&channels).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find channels by user: %w", err)
	}

	return channels, nil
}

// CountChannelsByUser returns the number of channels for a specific user
func (r *channelRepository) CountChannelsByUser(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Channel{}).Where("user_id = ?", userID).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count channels by user: %w", err)
	}

	return count, nil
}
