package handlers

import (
	"net/http"

	"copier/internal/database/models"
	"copier/internal/services"
	AppError "copier/internal/shared/error"
	"copier/internal/shared/response"
	"copier/internal/shared/utils"
)

// ChannelHandler handles HTTP requests related to channels
type ChannelHandler struct {
	channelService services.ChannelService
}

// NewChannelHandler creates a new ChannelHandler instance
func NewChannelHandler(channelService services.ChannelService) *ChannelHandler {
	return &ChannelHandler{
		channelService: channelService,
	}
}

// CreateChannelRequest defines the payload for channel creation
type CreateChannelRequest struct {
	Name      string `json:"name" validate:"required,min=1,max=255"`
	ChannelID string `json:"channel_id" validate:"required"`
}

func (h *ChannelHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.ResolveAuthenticatedUser(w, r)
	if !ok {
		return
	}

	var req CreateChannelRequest
	if !utils.DecodeAndValidate(w, r, &req) {
		return
	}

	channel, err := h.channelService.CreateChannel(r.Context(), userID, req.Name, req.ChannelID)
	if err != nil {
		AppError.InternalServerErrorWithError("Failed to create channel", err).WriteToResponse(w)
		return
	}

	response.WriteCreated(w, "Channel created successfully", channel)
}

func (h *ChannelHandler) ListByUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.ResolveUserID(w, r, "user_id")
	if !ok {
		return
	}

	channels, err := h.channelService.GetChannelsByUser(r.Context(), userID)
	if err != nil {
		AppError.InternalServerErrorWithError("Failed to retrieve channels", err).WriteToResponse(w)
		return
	}

	response.WriteOK(w, "Channels retrieved successfully", channels)
}

func (h *ChannelHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ResolvePathID(w, r, "id")
	if !ok {
		return
	}

	var update models.Channel
	if !utils.DecodeAndValidate(w, r, &update) {
		return
	}

	channel, err := h.channelService.UpdateChannel(r.Context(), id, &update)
	if err != nil {
		AppError.InternalServerErrorWithError("Failed to update channel", err).WriteToResponse(w)
		return
	}

	response.WriteOK(w, "Channel updated successfully", channel)
}

func (h *ChannelHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ResolvePathID(w, r, "id")
	if !ok {
		return
	}

	if err := h.channelService.DeleteChannel(r.Context(), id); err != nil {
		AppError.InternalServerErrorWithError("Failed to delete channel", err).WriteToResponse(w)
		return
	}

	response.WriteNoContent(w)
}
