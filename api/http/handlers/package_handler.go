package handlers

import (
	"net/http"
	"strconv"

	"copier/internal/services"
	AppError "copier/internal/shared/error"
	"copier/internal/shared/response"

	"github.com/google/uuid"
)

// PackageHandler handles HTTP requests related to packages
type PackageHandler struct {
	packageService services.PackageService
}

// NewPackageHandler creates a new PackageHandler instance
func NewPackageHandler(packageService services.PackageService) *PackageHandler {
	return &PackageHandler{
		packageService: packageService,
	}
}

// List retrieves all packages (with optional pagination)
func (h *PackageHandler) List(w http.ResponseWriter, r *http.Request) {
	skip, _ := strconv.Atoi(r.URL.Query().Get("skip"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 10
	}

	packages, total, err := h.packageService.GetAllPackages(r.Context(), skip, limit)
	if err != nil {
		AppError.InternalServerErrorWithError("Failed to retrieve packages", err).WriteToResponse(w)
		return
	}

	response.WriteOK(w, "Packages retrieved successfully", map[string]interface{}{
		"packages": packages,
		"total":    total,
	})
}

// GetByID retrieves a single package by its ID
func (h *PackageHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		idStr = r.URL.Query().Get("id")
	}

	if idStr == "" {
		AppError.BadRequest("Package ID is required").WriteToResponse(w)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		AppError.BadRequest("Invalid package ID format").WriteToResponse(w)
		return
	}

	pkg, err := h.packageService.GetPackageByID(r.Context(), id)
	if err != nil {
		AppError.ResourceNotFound("Package", idStr).WriteToResponse(w)
		return
	}

	response.WriteOK(w, "Package details retrieved successfully", pkg)
}
