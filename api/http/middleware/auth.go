package middleware

import (
	"context"
	"copier/config"
	AppError "copier/internal/shared/error"
	"copier/internal/shared/types"
	"copier/internal/shared/utils"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := extractTokenFromHeader(r)
		if err != nil {
			AppError.Unauthorized("Unauthorized: Invalid token format").WriteToResponse(w)
			return
		}

		// Validate and parse the token
		claims, err := validateToken(tokenString)
		if err != nil {
			AppError.Unauthorized("Unauthorized: Invalid token").WriteToResponse(w)
			return
		}

		// Check if user is active
		if !claims.IsActive {
			AppError.Unauthorized("Unauthorized: User is inactive").WriteToResponse(w)
			return
		}

		// Add user information to request context
		ctx := context.WithValue(r.Context(), types.UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func OptionalAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := extractTokenFromHeader(r)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		// Validate and parse the token
		claims, err := validateToken(tokenString)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		// Check if user is active
		if !claims.IsActive {
			next.ServeHTTP(w, r)
			return
		}

		// Add user information to request context
		ctx := context.WithValue(r.Context(), types.UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RoleMiddleware(requiredRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userCtx, ok := utils.GetUserFromRequest(r)
			if !ok {
				AppError.Unauthorized("Unauthorized: User not authenticated").WriteToResponse(w)
				return
			}

			hasRole := slices.Contains(requiredRoles, userCtx.Role)

			if !hasRole {
				AppError.Forbidden("Forbidden: Insufficient permissions").WriteToResponse(w)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func extractTokenFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", AppError.Unauthorized("Unauthorized: No token provided")
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", AppError.Unauthorized("Unauthorized: Invalid token format")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		return "", AppError.Unauthorized("Unauthorized: Invalid token")
	}

	return token, nil
}

func validateToken(tokenString string) (*types.UserClaims, error) {
	cfg := config.GetConfig()

	token, err := jwt.ParseWithClaims(tokenString, &types.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.JwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*types.UserClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, AppError.Unauthorized("Unauthorized: Invalid token")
}
