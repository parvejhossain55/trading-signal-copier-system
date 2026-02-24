package repositories

import (
	"context"
	"fmt"

	"copier/internal/database/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PlatformRepository defines platform-specific repository operations
type PlatformRepository interface {
	BaseRepository
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Platform, error)
	FindByUserAndName(ctx context.Context, userID uuid.UUID, name string) (*models.Platform, error)
	FindByAPIKey(ctx context.Context, apiKey string) (*models.Platform, error)
	FindByIDTyped(ctx context.Context, id uuid.UUID) (*models.Platform, error)
	CreatePlatform(ctx context.Context, platform *models.Platform) error
	UpdatePlatform(ctx context.Context, id uuid.UUID, update *models.Platform) error
	DeletePlatform(ctx context.Context, id uuid.UUID) error
	DeletePlatformsByUser(ctx context.Context, userID uuid.UUID) error
	FindAllByUser(ctx context.Context, userID uuid.UUID, skip, limit int) ([]*models.Platform, error)
	CountPlatformsByUser(ctx context.Context, userID uuid.UUID) (int64, error)
	UpdateAPIKeys(ctx context.Context, id uuid.UUID, apiKey, apiSecret string) error
}

// platformRepository implements PlatformRepository interface
type platformRepository struct {
	BaseRepository
	db *gorm.DB
}

// NewPlatformRepository creates a new platform repository instance
func NewPlatformRepository(db *gorm.DB) PlatformRepository {
	return &platformRepository{
		BaseRepository: NewBaseRepository(db),
		db:             db,
	}
}

// FindByUserID finds all platforms for a specific user
func (r *platformRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Platform, error) {
	var platforms []*models.Platform
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&platforms).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find platforms by user ID: %w", err)
	}

	return platforms, nil
}

// FindByUserAndName finds a platform by user ID and name
func (r *platformRepository) FindByUserAndName(ctx context.Context, userID uuid.UUID, name string) (*models.Platform, error) {
	var platform models.Platform
	err := r.db.WithContext(ctx).Where("user_id = ? AND name = ?", userID, name).First(&platform).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("platform not found with user ID: %s and name: %s", userID, name)
		}
		return nil, fmt.Errorf("failed to find platform by user and name: %w", err)
	}

	return &platform, nil
}

// FindByAPIKey finds a platform by API key
func (r *platformRepository) FindByAPIKey(ctx context.Context, apiKey string) (*models.Platform, error) {
	var platform models.Platform
	err := r.db.WithContext(ctx).Where("api_key = ?", apiKey).First(&platform).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("platform not found with API key")
		}
		return nil, fmt.Errorf("failed to find platform by API key: %w", err)
	}

	return &platform, nil
}

// FindByIDTyped finds a platform by ID and returns typed Platform struct
func (r *platformRepository) FindByIDTyped(ctx context.Context, id uuid.UUID) (*models.Platform, error) {
	var platform models.Platform
	err := r.db.WithContext(ctx).First(&platform, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("platform not found with ID: %s", id)
		}
		return nil, fmt.Errorf("failed to find platform by ID: %w", err)
	}

	return &platform, nil
}

// CreatePlatform creates a new platform
func (r *platformRepository) CreatePlatform(ctx context.Context, platform *models.Platform) error {
	err := r.db.WithContext(ctx).Create(platform).Error
	if err != nil {
		return fmt.Errorf("failed to create platform: %w", err)
	}

	return nil
}

// UpdatePlatform updates an existing platform
func (r *platformRepository) UpdatePlatform(ctx context.Context, id uuid.UUID, update *models.Platform) error {
	err := r.db.WithContext(ctx).Model(&models.Platform{}).Where("id = ?", id).Updates(update).Error
	if err != nil {
		return fmt.Errorf("failed to update platform: %w", err)
	}

	return nil
}

// DeletePlatform deletes a platform by ID
func (r *platformRepository) DeletePlatform(ctx context.Context, id uuid.UUID) error {
	err := r.db.WithContext(ctx).Delete(&models.Platform{}, id).Error
	if err != nil {
		return fmt.Errorf("failed to delete platform: %w", err)
	}

	return nil
}

// DeletePlatformsByUser deletes all platforms for a specific user
func (r *platformRepository) DeletePlatformsByUser(ctx context.Context, userID uuid.UUID) error {
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.Platform{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete platforms by user: %w", err)
	}

	return nil
}

// FindAllByUser retrieves all platforms for a user with pagination
func (r *platformRepository) FindAllByUser(ctx context.Context, userID uuid.UUID, skip, limit int) ([]*models.Platform, error) {
	var platforms []*models.Platform
	query := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at desc")
	if skip > 0 {
		query = query.Offset(skip)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&platforms).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find platforms by user: %w", err)
	}

	return platforms, nil
}

// CountPlatformsByUser returns the number of platforms for a specific user
func (r *platformRepository) CountPlatformsByUser(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Platform{}).Where("user_id = ?", userID).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count platforms by user: %w", err)
	}

	return count, nil
}

// UpdateAPIKeys updates only the API key and secret for a platform
func (r *platformRepository) UpdateAPIKeys(ctx context.Context, id uuid.UUID, apiKey, apiSecret string) error {
	err := r.db.WithContext(ctx).Model(&models.Platform{}).Where("id = ?", id).Updates(map[string]interface{}{
		"api_key":    apiKey,
		"api_secret": apiSecret,
	}).Error
	if err != nil {
		return fmt.Errorf("failed to update platform API keys: %w", err)
	}

	return nil
}
