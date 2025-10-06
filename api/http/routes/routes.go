package routes

import (
	"copier/http/handlers"
	"copier/http/middleware"
	"net/http"
)

// Router handles HTTP routing
// SetupRoutes configures all routes with global middleware
func SetupRoutes(mux *http.ServeMux) http.Handler {

	// Initialize middleware manager
	manager := middleware.NewManager()

	// Initialize handlers
	welcomeHandler := handlers.NewWelcomeHandler()
	healthHandler := handlers.NewHealthHandler()
	notFoundHandler := handlers.NewNotFoundHandler()

	// Define routes with middleware
	mux.Handle("GET /api", manager.With(http.HandlerFunc(welcomeHandler.Welcome)))
	mux.Handle("GET /api/v1", manager.With(http.HandlerFunc(welcomeHandler.Welcome)))
	mux.Handle("GET /health", manager.With(http.HandlerFunc(healthHandler.HealthCheck)))

	// Catch-all route for 404 Not Found
	mux.HandleFunc("/", notFoundHandler.NotFound)

	return mux
}
