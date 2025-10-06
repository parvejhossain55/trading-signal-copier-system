package config

import (
	"fmt"
)

// Mode represents the application mode
type Mode string

const (
	Development Mode = "development"
	Production  Mode = "production"
	Testing     Mode = "testing"
)

// MongoDatabase represents MongoDB database configuration
type MongoDatabase struct {
	Host                string `envconfig:"MONGO_HOST" default:"localhost"`
	Database            string `envconfig:"MONGO_DATABASE"`
	User                string `envconfig:"MONGO_USER"`
	Password            string `envconfig:"MONGO_PASSWORD"`
	Port                int    `envconfig:"MONGO_PORT" default:"27017"`
	AuthSource          string `envconfig:"MONGO_AUTH_SOURCE" default:"admin"`
	SSL                 bool   `envconfig:"MONGO_SSL" default:"false"`
	MaxPoolSize         uint64 `envconfig:"MONGO_MAX_POOL_SIZE" default:"100"`
	MinPoolSize         uint64 `envconfig:"MONGO_MIN_POOL_SIZE" default:"1"`
	MaxConnIdleTimeInMs int    `envconfig:"MONGO_MAX_CONN_IDLE_TIME_MS" default:"300000"`
}

// RedisConfig represents Redis cache configuration
type RedisConfig struct {
	Host     string `envconfig:"REDIS_HOST" default:"localhost"`
	Port     int    `envconfig:"REDIS_PORT" default:"6379"`
	Password string `envconfig:"REDIS_PASSWORD"`
	DB       int    `envconfig:"REDIS_DB" default:"0"`
	PoolSize int    `envconfig:"REDIS_POOL_SIZE" default:"10"`
	MinIdle  int    `envconfig:"REDIS_MIN_IDLE" default:"5"`
}


type CorsOrigin struct {
	Origin []string `envconfig:"CORS_ORIGIN"`
}

// Config represents the application configuration
type Config struct {
	Version          string `envconfig:"VERSION" default:"1.0.0"`
	Mode             Mode   `envconfig:"MODE" default:"development"`
	ServiceName      string `envconfig:"SERVICE_NAME" default:"starter-service"`
	HttpPort         int    `envconfig:"HTTP_PORT" default:"8080"`
	GrpcPort         int    `envconfig:"GRPC_PORT" default:"8081"`
	HealthCheckRoute string `envconfig:"HEALTH_CHECK_ROUTE" default:"/health"`
	ApiVersion       string `envconfig:"API_VERSION" default:"v1"`
	JwtSecret        string `envconfig:"JWT_SECRET" default:"secret"`
	ServiceBasePath  string `envconfig:"SERVICE_BASE_PATH" default:"/api/v1"`

	// Database configuration
	Database MongoDatabase

	// Cache configuration
	Redis RedisConfig

	

	// Logging
	LogLevel string `envconfig:"LOG_LEVEL" default:"info"`

	// CORS configuration
	Cors CorsOrigin
}

var config *Config

// GetConfig returns the application configuration
func GetConfig() *Config {
	if config == nil {
		// Load configuration from .env file and environment variables
		if err := LoadConfig(); err != nil {
			panic(err)
		}
	}
	return config
}

// IsDevelopment checks if the application is in development mode
func (c *Config) IsDevelopment() bool {
	return c.Mode == Development
}

// IsProduction checks if the application is in production mode
func (c *Config) IsProduction() bool {
	return c.Mode == Production
}

// IsTesting checks if the application is in testing mode
func (c *Config) IsTesting() bool {
	return c.Mode == Testing
}

// GetDatabaseURL returns the MongoDB connection string
func (c *Config) GetDatabaseURL() string {
	// Use the new Database configuration if available
	if c.Database.User != "" && c.Database.Password != "" && c.Database.Database != "" {
		sslMode := ""
		if c.Database.SSL {
			sslMode = "&ssl=true"
		}
		return fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?authSource=%s%s",
			c.Database.User, c.Database.Password, c.Database.Host, c.Database.Port, c.Database.Database, c.Database.AuthSource, sslMode)
	}

	return ""
}
