package services

import (
	"context"

	"copier/database/repositories"
	"copier/internal/database/models"

	"github.com/google/uuid"
)

// ChannelService defines channel business logic operations
type ChannelService interface {
	CreateChannel(ctx context.Context, userID uuid.UUID, name, channelID string) (*models.Channel, error)
	GetChannelsByUser(ctx context.Context, userID uuid.UUID) ([]*models.Channel, error)
	GetChannelByID(ctx context.Context, id uuid.UUID) (*models.Channel, error)
	UpdateChannel(ctx context.Context, id uuid.UUID, update *models.Channel) (*models.Channel, error)
	DeleteChannel(ctx context.Context, id uuid.UUID) error
}

type channelService struct {
	channelRepo repositories.ChannelRepository
}

// NewChannelService creates a new channel service instance
func NewChannelService(channelRepo repositories.ChannelRepository) ChannelService {
	return &channelService{
		channelRepo: channelRepo,
	}
}

func (s *channelService) CreateChannel(ctx context.Context, userID uuid.UUID, name, channelID string) (*models.Channel, error) {
	channel := &models.Channel{
		UserID:    userID,
		Name:      name,
		ChannelID: channelID,
	}

	err := s.channelRepo.CreateChannel(ctx, channel)
	if err != nil {
		return nil, err
	}

	return channel, nil
}

func (s *channelService) GetChannelsByUser(ctx context.Context, userID uuid.UUID) ([]*models.Channel, error) {
	return s.channelRepo.FindByUserID(ctx, userID)
}

func (s *channelService) GetChannelByID(ctx context.Context, id uuid.UUID) (*models.Channel, error) {
	return s.channelRepo.FindByIDTyped(ctx, id)
}

func (s *channelService) UpdateChannel(ctx context.Context, id uuid.UUID, update *models.Channel) (*models.Channel, error) {
	err := s.channelRepo.UpdateChannel(ctx, id, update)
	if err != nil {
		return nil, err
	}
	return s.channelRepo.FindByIDTyped(ctx, id)
}

func (s *channelService) DeleteChannel(ctx context.Context, id uuid.UUID) error {
	return s.channelRepo.DeleteChannel(ctx, id)
}
