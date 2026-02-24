package handlers

import (
	"net/http"

	"copier/internal/database/models"
	"copier/internal/services"
	AppError "copier/internal/shared/error"
	"copier/internal/shared/response"
	"copier/internal/shared/utils"
)

// PlatformHandler handles HTTP requests related to platforms
type PlatformHandler struct {
	platformService services.PlatformService
}

// NewPlatformHandler creates a new PlatformHandler instance
func NewPlatformHandler(platformService services.PlatformService) *PlatformHandler {
	return &PlatformHandler{
		platformService: platformService,
	}
}

// CreatePlatformRequest defines the payload for platform creation
type CreatePlatformRequest struct {
	Name      string `json:"name" validate:"required,min=1,max=255"`
	APIKey    string `json:"api_key" validate:"required"`
	APISecret string `json:"api_secret" validate:"required"`
}

func (h *PlatformHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.ResolveAuthenticatedUser(w, r)
	if !ok {
		return
	}

	var req CreatePlatformRequest
	if !utils.DecodeAndValidate(w, r, &req) {
		return
	}

	platform, err := h.platformService.CreatePlatform(r.Context(), userID, req.Name, req.APIKey, req.APISecret)
	if err != nil {
		AppError.InternalServerErrorWithError("Failed to create platform", err).WriteToResponse(w)
		return
	}

	response.WriteCreated(w, "Platform created successfully", platform)
}

func (h *PlatformHandler) ListByUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.ResolveUserID(w, r, "user_id")
	if !ok {
		return
	}

	platforms, err := h.platformService.GetPlatformsByUser(r.Context(), userID)
	if err != nil {
		AppError.InternalServerErrorWithError("Failed to retrieve platforms", err).WriteToResponse(w)
		return
	}

	response.WriteOK(w, "Platforms retrieved successfully", platforms)
}

func (h *PlatformHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ResolvePathID(w, r, "id")
	if !ok {
		return
	}

	var update models.Platform
	if !utils.DecodeAndValidate(w, r, &update) {
		return
	}

	platform, err := h.platformService.UpdatePlatform(r.Context(), id, &update)
	if err != nil {
		AppError.InternalServerErrorWithError("Failed to update platform", err).WriteToResponse(w)
		return
	}

	response.WriteOK(w, "Platform updated successfully", platform)
}

func (h *PlatformHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ResolvePathID(w, r, "id")
	if !ok {
		return
	}

	if err := h.platformService.DeletePlatform(r.Context(), id); err != nil {
		AppError.InternalServerErrorWithError("Failed to delete platform", err).WriteToResponse(w)
		return
	}

	response.WriteNoContent(w)
}
