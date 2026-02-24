package middleware

import (
	"net/http"
	"time"
)

type Middleware func(http.Handler) http.Handler

type Manager struct {
	globalMiddlewares []Middleware
}

func NewManager() *Manager {
	manager := &Manager{
		globalMiddlewares: make([]Middleware, 0),
	}

	manager.UseDefault()

	return manager
}

func (m *Manager) UseDefault() *Manager {
	return m.Use(
		RecoveryMiddleware,                        // Panic recovery
		RequestIDMiddleware,                       // Request ID generation
		LoggingMiddleware,                         // Request logging
		SecurityHeadersMiddleware,                 // Security headers
		CORSMiddleware,                            // CORS handling
		TimeoutMiddleware(30*time.Second),         // Request timeout
		RateLimitMiddleware(60),                   // Rate limiting (60 requests per minute)
		ContentTypeMiddleware("application/json"), // Content type validation
	)
}

func (m *Manager) Use(middlewares ...Middleware) *Manager {
	m.globalMiddlewares = append(m.globalMiddlewares, middlewares...)
	return m
}

func (m *Manager) With(handler http.Handler, middlewares ...Middleware) http.Handler {
	var h http.Handler
	h = handler

	for _, m := range middlewares {
		h = m(h)
	}

	for _, m := range m.globalMiddlewares {
		h = m(h)
	}

	return h
}
