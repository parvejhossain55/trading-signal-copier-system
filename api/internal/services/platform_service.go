package services

import (
	"context"

	"copier/database/repositories"
	"copier/internal/database/models"

	"github.com/google/uuid"
)

// PlatformService defines platform business logic operations
type PlatformService interface {
	CreatePlatform(ctx context.Context, userID uuid.UUID, name, apiKey, apiSecret string) (*models.Platform, error)
	GetPlatformsByUser(ctx context.Context, userID uuid.UUID) ([]*models.Platform, error)
	GetPlatformByID(ctx context.Context, id uuid.UUID) (*models.Platform, error)
	UpdatePlatform(ctx context.Context, id uuid.UUID, update *models.Platform) (*models.Platform, error)
	DeletePlatform(ctx context.Context, id uuid.UUID) error
}

type platformService struct {
	platformRepo repositories.PlatformRepository
}

// NewPlatformService creates a new platform service instance
func NewPlatformService(platformRepo repositories.PlatformRepository) PlatformService {
	return &platformService{
		platformRepo: platformRepo,
	}
}

func (s *platformService) CreatePlatform(ctx context.Context, userID uuid.UUID, name, apiKey, apiSecret string) (*models.Platform, error) {
	platform := &models.Platform{
		UserID:    userID,
		Name:      name,
		APIKey:    apiKey,
		APISecret: apiSecret,
	}

	err := s.platformRepo.CreatePlatform(ctx, platform)
	if err != nil {
		return nil, err
	}

	return platform, nil
}

func (s *platformService) GetPlatformsByUser(ctx context.Context, userID uuid.UUID) ([]*models.Platform, error) {
	return s.platformRepo.FindByUserID(ctx, userID)
}

func (s *platformService) GetPlatformByID(ctx context.Context, id uuid.UUID) (*models.Platform, error) {
	return s.platformRepo.FindByIDTyped(ctx, id)
}

func (s *platformService) UpdatePlatform(ctx context.Context, id uuid.UUID, update *models.Platform) (*models.Platform, error) {
	err := s.platformRepo.UpdatePlatform(ctx, id, update)
	if err != nil {
		return nil, err
	}
	return s.platformRepo.FindByIDTyped(ctx, id)
}

func (s *platformService) DeletePlatform(ctx context.Context, id uuid.UUID) error {
	return s.platformRepo.DeletePlatform(ctx, id)
}
