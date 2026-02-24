package handlers

import (
	"net/http"

	"copier/internal/services"
	AppError "copier/internal/shared/error"
	"copier/internal/shared/response"
	"copier/internal/shared/utils"
)

// UserHandler handles HTTP requests related to users
type UserHandler struct {
	userService services.UserService
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// UserRegisterRequest defines the payload for user registration
type UserRegisterRequest struct {
	Name  string  `json:"name" validate:"required,min=2,max=100"`
	Email string  `json:"email" validate:"required,email"`
	Phone *string `json:"phone" validate:"omitempty,min=10,max=20"`
}

// Register handles user registration
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req UserRegisterRequest

	if !utils.DecodeAndValidate(w, r, &req) {
		return
	}

	user, err := h.userService.RegisterUser(r.Context(), req.Name, req.Email, req.Phone)
	if err != nil {
		AppError.InternalServerErrorWithError("Failed to register user", err).WriteToResponse(w)
		return
	}

	response.WriteCreated(w, "User registered successfully", user)
}

// GetProfile retrieves the user profile
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.ResolveUserID(w, r, "id")
	if !ok {
		return
	}

	user, err := h.userService.GetUserByID(r.Context(), userID)
	if err != nil {
		AppError.ResourceNotFound("User", userID.String()).WriteToResponse(w)
		return
	}

	response.WriteOK(w, "User profile retrieved successfully", user)
}
