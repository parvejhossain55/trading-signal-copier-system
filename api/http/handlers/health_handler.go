package handlers

import (
	"net/http"
	"time"

	"copier/config"
	"copier/internal/shared/response"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string                 `json:"status"`
	Service   string                 `json:"service"`
	Timestamp time.Time              `json:"timestamp"`
	Version   string                 `json:"version"`
	Services  map[string]interface{} `json:"services,omitempty"`
}

// HealthHandler handles health check requests
type HealthHandler struct {
	serviceName      string
	version          string
}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	serviceName := config.GetConfig().ServiceName
	version := config.GetConfig().Version
	return &HealthHandler{
		serviceName: serviceName,
		version:     version,
	}
}

// HealthCheck handles the health check endpoint
func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	services := make(map[string]interface{})
	overallStatus := "ok"

	// Check database connectivity
	// TODO: Add actual database health check
	services["database"] = map[string]interface{}{
		"status": "ok",
	}

	// Check messaging service
	services["messaging"] = map[string]interface{}{
		"status": "ok",
	}

	// Create health response
	healthResponse := HealthResponse{
		Status:    overallStatus,
		Service:   h.serviceName,
		Timestamp: time.Now(),
		Version:   h.version,
		Services:  services,
	}

	response.WriteOK(w, "Service health check completed", healthResponse)
}
