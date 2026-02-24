package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	AppError "copier/internal/shared/error"

	"github.com/google/uuid"
)

// DecodeAndValidate decodes the request body and validates the struct
func DecodeAndValidate(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		AppError.BadRequest("Invalid request body").WriteToResponse(w)
		return false
	}

	if errs := ValidateStruct(v); errs != nil {
		AppError.MultipleValidationErrors(errs).WriteToResponse(w)
		return false
	}

	return true
}

// ResolveUserID extract user ID from context and performs optional ID parsing and ownership checks
func ResolveUserID(w http.ResponseWriter, r *http.Request, queryParam string) (uuid.UUID, bool) {
	claims, ok := GetUserFromRequest(r)
	if !ok {
		AppError.Unauthorized("Unauthorized").WriteToResponse(w)
		return uuid.Nil, false
	}

	// If no query parameter specified, return current user's ID
	if queryParam == "" {
		id, err := uuid.Parse(claims.UserID)
		if err != nil {
			AppError.InternalServerError("Invalid user ID in token").WriteToResponse(w)
			return uuid.Nil, false
		}
		return id, true
	}

	// Parse ID from query parameter
	idStr := r.URL.Query().Get(queryParam)
	if idStr == "" {
		// Default to authenticated user
		id, err := uuid.Parse(claims.UserID)
		if err != nil {
			AppError.InternalServerError("Invalid user ID in token").WriteToResponse(w)
			return uuid.Nil, false
		}
		return id, true
	}

	parsedID, err := uuid.Parse(idStr)
	if err != nil {
		AppError.BadRequest(fmt.Sprintf("Invalid %s format", queryParam)).WriteToResponse(w)
		return uuid.Nil, false
	}

	// Security check: only admins or the user themselves can access this
	if claims.Role != "admin" && claims.UserID != idStr {
		AppError.Forbidden("Insufficient permissions").WriteToResponse(w)
		return uuid.Nil, false
	}

	return parsedID, true
}

// ResolvePathID extracts a UUID from a path parameter
func ResolvePathID(w http.ResponseWriter, r *http.Request, key string) (uuid.UUID, bool) {
	idStr := r.PathValue(key)
	if idStr == "" {
		AppError.BadRequest(fmt.Sprintf("%s is required", key)).WriteToResponse(w)
		return uuid.Nil, false
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		AppError.BadRequest(fmt.Sprintf("Invalid %s format", key)).WriteToResponse(w)
		return uuid.Nil, false
	}

	return id, true
}

// ResolveAuthenticatedUser extracts the user ID from the token context
func ResolveAuthenticatedUser(w http.ResponseWriter, r *http.Request) (uuid.UUID, bool) {
	claims, ok := GetUserFromRequest(r)
	if !ok {
		AppError.Unauthorized("Unauthorized").WriteToResponse(w)
		return uuid.Nil, false
	}

	id, err := uuid.Parse(claims.UserID)
	if err != nil {
		AppError.InternalServerError("Invalid user ID in token").WriteToResponse(w)
		return uuid.Nil, false
	}

	return id, true
}
