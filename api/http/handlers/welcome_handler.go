package handlers

import (
	"copier/config"
	"copier/internal/shared/response"
	"net/http"
	"time"
)

type WelcomeResponse struct {
	Message        string    `json:"message"`
	Service        string    `json:"service"`
	Timestamp      time.Time `json:"timestamp"`
	ProjectVersion string    `json:"project_version"`
}

type WelcomeHandler struct {
	serviceName string
}

func NewWelcomeHandler() *WelcomeHandler {
	return &WelcomeHandler{
		serviceName: config.GetConfig().ServiceName,
	}
}

func (h *WelcomeHandler) Welcome(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"message":         "Welcome to the " + h.serviceName + " API!",
		"service":         h.serviceName,
		"project_version": config.GetConfig().Version,
	}

	response.WriteOK(w, data)
}
