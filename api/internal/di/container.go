package di

import (
	"copier/database/repositories"
	"copier/http/handlers"
	"copier/internal/services"

	"gorm.io/gorm"
)

// Container holds all application dependencies
type Container struct {
	DB *gorm.DB

	// Repositories
	UserRepo             repositories.UserRepository
	ChannelRepo          repositories.ChannelRepository
	PlatformRepo         repositories.PlatformRepository
	TradeSettingsRepo    repositories.TradeSettingsRepository
	PackageRepo          repositories.PackageRepository
	SubscribePackageRepo repositories.SubscribePackageRepository

	// Services
	UserService          services.UserService
	PackageService       services.PackageService
	SubscriptionService  services.SubscriptionService
	PlatformService      services.PlatformService
	ChannelService       services.ChannelService
	TradeSettingsService services.TradeSettingsService

	// Handlers
	UserHandler          *handlers.UserHandler
	PackageHandler       *handlers.PackageHandler
	SubscriptionHandler  *handlers.SubscriptionHandler
	PlatformHandler      *handlers.PlatformHandler
	ChannelHandler       *handlers.ChannelHandler
	TradeSettingsHandler *handlers.TradeSettingsHandler
	WelcomeHandler       *handlers.WelcomeHandler
	HealthHandler        *handlers.HealthHandler
	NotFoundHandler      *handlers.NotFoundHandler
}

// NewContainer initializes all dependencies and returns a new Container instance
func NewContainer(db *gorm.DB) *Container {
	// 1. Repositories
	userRepo := repositories.NewUserRepository(db)
	channelRepo := repositories.NewChannelRepository(db)
	platformRepo := repositories.NewPlatformRepository(db)
	tradeSettingsRepo := repositories.NewTradeSettingsRepository(db)
	packageRepo := repositories.NewPackageRepository(db)
	subscribePackageRepo := repositories.NewSubscribePackageRepository(db)

	// 2. Services
	userService := services.NewUserService(userRepo)
	packageService := services.NewPackageService(packageRepo)
	subscriptionService := services.NewSubscriptionService(subscribePackageRepo, packageRepo)
	platformService := services.NewPlatformService(platformRepo)
	channelService := services.NewChannelService(channelRepo)
	tradeSettingsService := services.NewTradeSettingsService(tradeSettingsRepo)

	// 3. Handlers
	userHandler := handlers.NewUserHandler(userService)
	packageHandler := handlers.NewPackageHandler(packageService)
	subscriptionHandler := handlers.NewSubscriptionHandler(subscriptionService)
	platformHandler := handlers.NewPlatformHandler(platformService)
	channelHandler := handlers.NewChannelHandler(channelService)
	tradeSettingsHandler := handlers.NewTradeSettingsHandler(tradeSettingsService)
	welcomeHandler := handlers.NewWelcomeHandler()
	healthHandler := handlers.NewHealthHandler()
	notFoundHandler := handlers.NewNotFoundHandler()

	return &Container{
		DB: db,

		// Repositories
		UserRepo:             userRepo,
		ChannelRepo:          channelRepo,
		PlatformRepo:         platformRepo,
		TradeSettingsRepo:    tradeSettingsRepo,
		PackageRepo:          packageRepo,
		SubscribePackageRepo: subscribePackageRepo,

		// Services
		UserService:          userService,
		PackageService:       packageService,
		SubscriptionService:  subscriptionService,
		PlatformService:      platformService,
		ChannelService:       channelService,
		TradeSettingsService: tradeSettingsService,

		// Handlers
		UserHandler:          userHandler,
		PackageHandler:       packageHandler,
		SubscriptionHandler:  subscriptionHandler,
		PlatformHandler:      platformHandler,
		ChannelHandler:       channelHandler,
		TradeSettingsHandler: tradeSettingsHandler,
		WelcomeHandler:       welcomeHandler,
		HealthHandler:        healthHandler,
		NotFoundHandler:      notFoundHandler,
	}
}
