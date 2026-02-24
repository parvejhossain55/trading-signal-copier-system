package AppError

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type APIError struct {
	StatusCode int                    `json:"-"`
	Status     string                 `json:"status"`
	Message    string                 `json:"message"`
	Details    map[string]interface{} `json:"details,omitempty"`
	Err        error                  `json:"-"`
}

func (e *APIError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s", e.Message, e.Err.Error())
	}
	return e.Message
}

func (e *APIError) Unwrap() error {
	return e.Err
}

func (e *APIError) WriteToResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.StatusCode)

	response := map[string]interface{}{
		"status":    "error",
		"info":      e.Status,
		"message":   e.Message,
		"timestamp": time.Now(),
	}

	if e.Details != nil {
		response["details"] = e.Details
	}

	json.NewEncoder(w).Encode(response)
}

func NewAPIError(statusCode int, code, message string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Status:     code,
		Message:    message,
	}
}

func NewAPIErrorWithDetails(statusCode int, code, message string, details map[string]interface{}) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Status:     code,
		Message:    message,
		Details:    details,
	}
}

func NewAPIErrorWithError(statusCode int, code, message string, err error) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Status:     code,
		Message:    message,
		Err:        err,
	}
}

func BadRequest(message string) *APIError {
	return NewAPIError(http.StatusBadRequest, "BAD_REQUEST", message)
}

func BadRequestWithDetails(message string, details map[string]interface{}) *APIError {
	return NewAPIErrorWithDetails(http.StatusBadRequest, "BAD_REQUEST", message, details)
}

func Unauthorized(message string) *APIError {
	return NewAPIError(http.StatusUnauthorized, "UNAUTHORIZED", message)
}

func Forbidden(message string) *APIError {
	return NewAPIError(http.StatusForbidden, "FORBIDDEN", message)
}

func NotFound(message string) *APIError {
	return NewAPIError(http.StatusNotFound, "NOT_FOUND", message)
}

func NotFoundWithDetails(message string, details map[string]interface{}) *APIError {
	return NewAPIErrorWithDetails(http.StatusNotFound, "NOT_FOUND", message, details)
}

func Conflict(message string) *APIError {
	return NewAPIError(http.StatusConflict, "CONFLICT", message)
}

func UnprocessableEntity(message string, details map[string]interface{}) *APIError {
	return NewAPIErrorWithDetails(http.StatusUnprocessableEntity, "UNPROCESSABLE_ENTITY", message, details)
}

func InternalServerError(message string) *APIError {
	return NewAPIError(http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", message)
}

func InternalServerErrorWithError(message string, err error) *APIError {
	return NewAPIErrorWithError(http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", message, err)
}

func ServiceUnavailable(message string) *APIError {
	return NewAPIError(http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", message)
}

func ValidationError(field, message string) *APIError {
	details := map[string]interface{}{
		"field": field,
		"error": message,
	}
	return NewAPIErrorWithDetails(http.StatusUnprocessableEntity, "VALIDATION_ERROR", "Validation failed", details)
}

func MultipleValidationErrors(errors map[string]string) *APIError {
	details := map[string]interface{}{
		"validation_errors": errors,
	}
	return NewAPIErrorWithDetails(http.StatusUnprocessableEntity, "VALIDATION_ERROR", "Multiple validation errors", details)
}

func ResourceNotFound(resourceType, identifier string) *APIError {
	message := fmt.Sprintf("%s with identifier '%s' not found", resourceType, identifier)
	details := map[string]interface{}{
		"resource_type": resourceType,
		"identifier":    identifier,
	}
	return NewAPIErrorWithDetails(http.StatusNotFound, "RESOURCE_NOT_FOUND", message, details)
}

func DuplicateResource(resourceType, field, value string) *APIError {
	message := fmt.Sprintf("%s with %s '%s' already exists", resourceType, field, value)
	details := map[string]interface{}{
		"resource_type": resourceType,
		"field":         field,
		"value":         value,
	}
	return NewAPIErrorWithDetails(http.StatusConflict, "DUPLICATE_RESOURCE", message, details)
}

func InsufficientPermissions(action, resource string) *APIError {
	message := fmt.Sprintf("Insufficient permissions to %s %s", action, resource)
	details := map[string]interface{}{
		"action":   action,
		"resource": resource,
	}
	return NewAPIErrorWithDetails(http.StatusForbidden, "INSUFFICIENT_PERMISSIONS", message, details)
}

func RateLimitExceeded(message string) *APIError {
	details := map[string]interface{}{
		"message": message,
	}
	return NewAPIErrorWithDetails(http.StatusTooManyRequests, "RATE_LIMIT_EXCEEDED", message, details)
}

func DatabaseError(operation string, err error) *APIError {
	message := fmt.Sprintf("Database %s failed", operation)
	details := map[string]interface{}{
		"operation": operation,
	}
	return NewAPIErrorWithDetails(http.StatusInternalServerError, "DATABASE_ERROR", message, details)
}

func ExternalServiceError(service, operation string, err error) *APIError {
	message := fmt.Sprintf("%s service %s failed", service, operation)
	details := map[string]interface{}{
		"service":   service,
		"operation": operation,
	}
	return NewAPIErrorWithDetails(http.StatusServiceUnavailable, "EXTERNAL_SERVICE_ERROR", message, details)
}

func UnsupportedMediaType(message string) *APIError {
	return NewAPIError(http.StatusUnsupportedMediaType, "UNSUPPORTED_MEDIA_TYPE", message)
}
