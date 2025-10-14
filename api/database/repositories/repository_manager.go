package repositories

import (
	"go.mongodb.org/mongo-driver/mongo"
)

// RepositoryManager manages all repository instances
type RepositoryManager struct {
	UserRepo            UserRepository
	ChannelRepo         ChannelRepository
	PlatformRepo        PlatformRepository
	TradeSettingsRepo   TradeSettingsRepository
	PackageRepo         PackageRepository
	SubscribePackageRepo SubscribePackageRepository
}

// NewRepositoryManager creates a new repository manager with all repositories
func NewRepositoryManager(db *mongo.Database) *RepositoryManager {
	return &RepositoryManager{
		UserRepo:            NewUserRepository(db.Collection("users")),
		ChannelRepo:         NewChannelRepository(db.Collection("channels")),
		PlatformRepo:        NewPlatformRepository(db.Collection("platforms")),
		TradeSettingsRepo:   NewTradeSettingsRepository(db.Collection("trade_settings")),
		PackageRepo:         NewPackageRepository(db.Collection("packages")),
		SubscribePackageRepo: NewSubscribePackageRepository(db.Collection("subscribe_packages")),
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
