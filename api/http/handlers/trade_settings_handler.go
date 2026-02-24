package handlers

import (
	"net/http"

	"copier/internal/database/models"
	"copier/internal/services"
	AppError "copier/internal/shared/error"
	"copier/internal/shared/response"
	"copier/internal/shared/utils"
)

// TradeSettingsHandler handles HTTP requests related to trade settings
type TradeSettingsHandler struct {
	settingsService services.TradeSettingsService
}

// NewTradeSettingsHandler creates a new TradeSettingsHandler instance
func NewTradeSettingsHandler(settingsService services.TradeSettingsService) *TradeSettingsHandler {
	return &TradeSettingsHandler{
		settingsService: settingsService,
	}
}

func (h *TradeSettingsHandler) GetByUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.ResolveUserID(w, r, "user_id")
	if !ok {
		return
	}

	settings, err := h.settingsService.GetTradeSettingsByUser(r.Context(), userID)
	if err != nil {
		AppError.ResourceNotFound("TradeSettings", userID.String()).WriteToResponse(w)
		return
	}

	response.WriteOK(w, "Trade settings retrieved successfully", settings)
}

// UpsertTradeSettingsRequest defines the payload for upserting trade settings
type UpsertTradeSettingsRequest struct {
	Settings models.TradeSettings `json:"settings" validate:"required"`
}

func (h *TradeSettingsHandler) Upsert(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.ResolveAuthenticatedUser(w, r)
	if !ok {
		return
	}

	var req UpsertTradeSettingsRequest
	if !utils.DecodeAndValidate(w, r, &req) {
		return
	}

	settings, err := h.settingsService.UpsertTradeSettings(r.Context(), userID, &req.Settings)
	if err != nil {
		AppError.InternalServerErrorWithError("Failed to upsert trade settings", err).WriteToResponse(w)
		return
	}

	response.WriteOK(w, "Trade settings updated successfully", settings)
}

// UpdateStopLossRequest defines the payload for updating stop loss settings
type UpdateStopLossRequest struct {
	Percentage int  `json:"percentage" validate:"required,min=0,max=100"`
	Status     bool `json:"status"`
}

func (h *TradeSettingsHandler) UpdateStopLoss(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.ResolveAuthenticatedUser(w, r)
	if !ok {
		return
	}

	var req UpdateStopLossRequest
	if !utils.DecodeAndValidate(w, r, &req) {
		return
	}

	if err := h.settingsService.UpdateStopLoss(r.Context(), userID, req.Percentage, req.Status); err != nil {
		AppError.InternalServerErrorWithError("Failed to update stop loss settings", err).WriteToResponse(w)
		return
	}

	response.WriteOK(w, "Stop loss settings updated successfully", nil)
}

// UpdateTakeProfitRequest defines the payload for updating take profit settings
type UpdateTakeProfitRequest struct {
	Status      bool      `json:"status"`
	Step        int       `json:"step" validate:"required,min=1"`
	Percentages []float64 `json:"percentages" validate:"required,dive,min=0,max=100"`
}

func (h *TradeSettingsHandler) UpdateTakeProfit(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.ResolveAuthenticatedUser(w, r)
	if !ok {
		return
	}

	var req UpdateTakeProfitRequest
	if !utils.DecodeAndValidate(w, r, &req) {
		return
	}

	if err := h.settingsService.UpdateTakeProfit(r.Context(), userID, req.Status, req.Step, req.Percentages); err != nil {
		AppError.InternalServerErrorWithError("Failed to update take profit settings", err).WriteToResponse(w)
		return
	}

	response.WriteOK(w, "Take profit settings updated successfully", nil)
}
