package handlers

import (
	"copier/internal/shared/response"
	"net/http"
)

type NotFoundHandler struct{}

func NewNotFoundHandler() *NotFoundHandler {
	return &NotFoundHandler{}
}

func (h *NotFoundHandler) NotFound(w http.ResponseWriter, r *http.Request) {
	response.WriteNotFound(w, "The requested resource was not found", "RESOURCE_NOT_FOUND")
}
