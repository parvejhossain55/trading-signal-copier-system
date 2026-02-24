package routes

import (
	"copier/http/middleware"
	"copier/internal/di"
	"net/http"
)

// SetupRoutes registers all routes to the mux using the DI container
func SetupRoutes(mux *http.ServeMux, container *di.Container) http.Handler {

	manager := middleware.NewManager()

	// Base Routes
	mux.Handle("GET /api", manager.With(http.HandlerFunc(container.WelcomeHandler.Welcome)))
	mux.Handle("GET /api/v1", manager.With(http.HandlerFunc(container.WelcomeHandler.Welcome)))
	mux.Handle("GET /health", manager.With(http.HandlerFunc(container.HealthHandler.HealthCheck)))

	// User Routes
	mux.Handle("POST /api/v1/users/register", manager.With(http.HandlerFunc(container.UserHandler.Register)))
	mux.Handle("GET /api/v1/users/profile", manager.With(middleware.AuthMiddleware(http.HandlerFunc(container.UserHandler.GetProfile))))

	// Package Routes
	mux.Handle("GET /api/v1/packages", manager.With(http.HandlerFunc(container.PackageHandler.List)))
	mux.Handle("GET /api/v1/packages/{id}", manager.With(http.HandlerFunc(container.PackageHandler.GetByID)))

	// Subscription Routes
	mux.Handle("POST /api/v1/subscriptions", manager.With(middleware.AuthMiddleware(http.HandlerFunc(container.SubscriptionHandler.Create))))
	mux.Handle("POST /api/v1/subscriptions/verify", manager.With(middleware.AuthMiddleware(http.HandlerFunc(container.SubscriptionHandler.Verify))))
	mux.Handle("GET /api/v1/subscriptions/my", manager.With(middleware.AuthMiddleware(http.HandlerFunc(container.SubscriptionHandler.MySubscriptions))))

	// Platform Routes
	mux.Handle("POST /api/v1/platforms", manager.With(middleware.AuthMiddleware(http.HandlerFunc(container.PlatformHandler.Create))))
	mux.Handle("GET /api/v1/platforms", manager.With(middleware.AuthMiddleware(http.HandlerFunc(container.PlatformHandler.ListByUser))))
	mux.Handle("PUT /api/v1/platforms/{id}", manager.With(middleware.AuthMiddleware(http.HandlerFunc(container.PlatformHandler.Update))))
	mux.Handle("DELETE /api/v1/platforms/{id}", manager.With(middleware.AuthMiddleware(http.HandlerFunc(container.PlatformHandler.Delete))))

	// Channel Routes
	mux.Handle("POST /api/v1/channels", manager.With(middleware.AuthMiddleware(http.HandlerFunc(container.ChannelHandler.Create))))
	mux.Handle("GET /api/v1/channels", manager.With(middleware.AuthMiddleware(http.HandlerFunc(container.ChannelHandler.ListByUser))))
	mux.Handle("PUT /api/v1/channels/{id}", manager.With(middleware.AuthMiddleware(http.HandlerFunc(container.ChannelHandler.Update))))
	mux.Handle("DELETE /api/v1/channels/{id}", manager.With(middleware.AuthMiddleware(http.HandlerFunc(container.ChannelHandler.Delete))))

	// Trade Settings Routes
	mux.Handle("GET /api/v1/trade-settings", manager.With(middleware.AuthMiddleware(http.HandlerFunc(container.TradeSettingsHandler.GetByUser))))
	mux.Handle("POST /api/v1/trade-settings", manager.With(middleware.AuthMiddleware(http.HandlerFunc(container.TradeSettingsHandler.Upsert))))
	mux.Handle("PATCH /api/v1/trade-settings/stop-loss", manager.With(middleware.AuthMiddleware(http.HandlerFunc(container.TradeSettingsHandler.UpdateStopLoss))))
	mux.Handle("PATCH /api/v1/trade-settings/take-profit", manager.With(middleware.AuthMiddleware(http.HandlerFunc(container.TradeSettingsHandler.UpdateTakeProfit))))

	mux.HandleFunc("/", container.NotFoundHandler.NotFound)

	return mux
}
