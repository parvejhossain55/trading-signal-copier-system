package handlers

import (
	"net/http"

	"copier/internal/database/models"
	"copier/internal/services"
	AppError "copier/internal/shared/error"
	"copier/internal/shared/response"
	"copier/internal/shared/utils"

	"github.com/google/uuid"
)

// SubscriptionHandler handles HTTP requests related to subscriptions
type SubscriptionHandler struct {
	subService services.SubscriptionService
}

// NewSubscriptionHandler creates a new SubscriptionHandler instance
func NewSubscriptionHandler(subService services.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{
		subService: subService,
	}
}

// CreateSubscriptionRequest defines the payload for initiating a subscription
type CreateSubscriptionRequest struct {
	PackageID     uuid.UUID            `json:"package_id" validate:"required"`
	PaymentMethod models.PaymentMethod `json:"payment_method" validate:"required,oneof=crypto binance_pay credit_card paypal"`
}

// Create initiates a new subscription
func (h *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.ResolveAuthenticatedUser(w, r)
	if !ok {
		return
	}

	var req CreateSubscriptionRequest
	if !utils.DecodeAndValidate(w, r, &req) {
		return
	}

	sub, err := h.subService.Subscribe(r.Context(), userID, req.PackageID, req.PaymentMethod)
	if err != nil {
		AppError.InternalServerErrorWithError("Failed to create subscription", err).WriteToResponse(w)
		return
	}

	response.WriteCreated(w, "Subscription initiated successfully", sub)
}

// VerifySubscriptionRequest defines the payload for payment verification
type VerifySubscriptionRequest struct {
	SubscriptionID uuid.UUID `json:"subscription_id" validate:"required"`
	TnxHash        *string   `json:"tnx_hash" validate:"required_without=BinancePayID"`
	BinancePayID   *string   `json:"binance_pay_id" validate:"required_without=TnxHash"`
}

// Verify verifies a payment and completes the subscription
func (h *SubscriptionHandler) Verify(w http.ResponseWriter, r *http.Request) {
	var req VerifySubscriptionRequest

	if !utils.DecodeAndValidate(w, r, &req) {
		return
	}

	sub, err := h.subService.VerifyPayment(r.Context(), req.SubscriptionID, req.TnxHash, req.BinancePayID)
	if err != nil {
		AppError.InternalServerErrorWithError("Failed to verify payment", err).WriteToResponse(w)
		return
	}

	response.WriteOK(w, "Payment verified and subscription activated", sub)
}

// MySubscriptions retrieves the subscriptions of the current user
func (h *SubscriptionHandler) MySubscriptions(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.ResolveUserID(w, r, "user_id")
	if !ok {
		return
	}

	subs, total, err := h.subService.GetUserSubscriptions(r.Context(), userID, 0, 100)
	if err != nil {
		AppError.InternalServerErrorWithError("Failed to retrieve subscriptions", err).WriteToResponse(w)
		return
	}

	response.WriteOK(w, "Subscriptions retrieved successfully", map[string]interface{}{
		"subscriptions": subs,
		"total":         total,
	})
}
