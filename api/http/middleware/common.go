package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"copier/config"
	"copier/internal/logger"
	AppError "copier/internal/shared/error"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		requestID := getRequestID(r.Context())

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)

		slog.Info("HTTP request",
			logger.Method(r.Method),
			logger.Path(r.URL.Path),
			logger.Status(wrapped.statusCode),
			logger.Latency(duration),
			logger.RequestID(requestID),
			logger.Ip(getClientIP(r)),
			logger.UserAgent(r.UserAgent()),
		)
	})
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", strings.Join(config.GetConfig().Cors.Origin, ","))
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Request-ID") == "" {
			requestID := generateRequestID()
			r.Header.Set("X-Request-ID", requestID)
		}

		w.Header().Set("X-Request-ID", r.Header.Get("X-Request-ID"))

		type ctxKeyRequestID struct{}
		ctx := context.WithValue(r.Context(), ctxKeyRequestID{}, r.Header.Get("X-Request-ID"))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("HTTP request panic recovered",
					logger.Method(r.Method),
					logger.Path(r.URL.Path),
					logger.Extra(fmt.Sprintf("panic: %v", err)),
					logger.RequestID(getRequestID(r.Context())),
					"stack_trace", string(debug.Stack()),
				)

				AppError.InternalServerError("Internal Server Error").WriteToResponse(w)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")

		next.ServeHTTP(w, r)
	})
}

func TimeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func generateRequestID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func getRequestID(ctx context.Context) string {
	type ctxKeyRequestID struct{}
	if id, ok := ctx.Value(ctxKeyRequestID{}).(string); ok {
		return id
	}
	return generateRequestID()
}

func getClientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}

	return r.RemoteAddr
}

func RateLimitMiddleware(requestsPerMinute int) func(http.Handler) http.Handler {
	clients := make(map[string][]time.Time)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientIP := getClientIP(r)
			now := time.Now()

			if times, exists := clients[clientIP]; exists {
				var validTimes []time.Time
				for _, t := range times {
					if now.Sub(t) < time.Minute {
						validTimes = append(validTimes, t)
					}
				}
				clients[clientIP] = validTimes
			}

			if len(clients[clientIP]) >= requestsPerMinute {
				AppError.RateLimitExceeded("Rate Limit Exceeded").WriteToResponse(w)
				return
			}

			clients[clientIP] = append(clients[clientIP], now)

			remaining := requestsPerMinute - len(clients[clientIP])
			w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", requestsPerMinute))
			w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
			w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", now.Add(time.Minute).Unix()))

			next.ServeHTTP(w, r)
		})
	}
}

func ContentTypeMiddleware(contentType string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
				if r.Header.Get("Content-Type") != contentType {
					AppError.UnsupportedMediaType("Unsupported Media Type").WriteToResponse(w)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
