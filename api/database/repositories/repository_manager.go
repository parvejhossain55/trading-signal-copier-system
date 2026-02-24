package repositories

import (
	"gorm.io/gorm"
)

// RepositoryManager manages all repository instances
type RepositoryManager struct {
	UserRepo             UserRepository
	ChannelRepo          ChannelRepository
	PlatformRepo         PlatformRepository
	TradeSettingsRepo    TradeSettingsRepository
	PackageRepo          PackageRepository
	SubscribePackageRepo SubscribePackageRepository
}

// NewRepositoryManager creates a new repository manager with all repositories
func NewRepositoryManager(db *gorm.DB) *RepositoryManager {
	return &RepositoryManager{
		UserRepo:             NewUserRepository(db),
		ChannelRepo:          NewChannelRepository(db),
		PlatformRepo:         NewPlatformRepository(db),
		TradeSettingsRepo:    NewTradeSettingsRepository(db),
		PackageRepo:          NewPackageRepository(db),
		SubscribePackageRepo: NewSubscribePackageRepository(db),
	}
}

// GetUserRepository returns the user repository
func (rm *RepositoryManager) GetUserRepository() UserRepository {
	return rm.UserRepo
}

// GetChannelRepository returns the channel repository
func (rm *RepositoryManager) GetChannelRepository() ChannelRepository {
	return rm.ChannelRepo
}

// GetPlatformRepository returns the platform repository
func (rm *RepositoryManager) GetPlatformRepository() PlatformRepository {
	return rm.PlatformRepo
}

// GetTradeSettingsRepository returns the trade settings repository
func (rm *RepositoryManager) GetTradeSettingsRepository() TradeSettingsRepository {
	return rm.TradeSettingsRepo
}

// GetPackageRepository returns the package repository
func (rm *RepositoryManager) GetPackageRepository() PackageRepository {
	return rm.PackageRepo
}

// GetSubscribePackageRepository returns the subscribe package repository
func (rm *RepositoryManager) GetSubscribePackageRepository() SubscribePackageRepository {
	return rm.SubscribePackageRepo
}
