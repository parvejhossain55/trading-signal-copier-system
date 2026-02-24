package config

import (
	"fmt"
)

type Mode string

const (
	Development Mode = "development"
	Production  Mode = "production"
	Testing     Mode = "testing"
)

type PostgresDatabase struct {
	Host            string `envconfig:"DB_HOST" default:"localhost"`
	Database        string `envconfig:"DB_NAME"`
	User            string `envconfig:"DB_USER"`
	Password        string `envconfig:"DB_PASSWORD"`
	Port            int    `envconfig:"DB_PORT" default:"5432"`
	SSLMode         string `envconfig:"DB_SSL_MODE" default:"disable"`
	TimeZone        string `envconfig:"DB_TIMEZONE" default:"UTC"`
	MaxIdleConns    int    `envconfig:"DB_MAX_IDLE_CONNS" default:"10"`
	MaxOpenConns    int    `envconfig:"DB_MAX_OPEN_CONNS" default:"100"`
	ConnMaxLifetime int    `envconfig:"DB_CONN_MAX_LIFETIME" default:"3600"`
}

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

	Database PostgresDatabase

	Redis RedisConfig

	LogLevel string `envconfig:"LOG_LEVEL" default:"info"`

	Cors CorsOrigin
}

var config *Config

func GetConfig() *Config {
	if config == nil {
		if err := LoadConfig(); err != nil {
			panic(err)
		}
	}
	return config
}

func (c *Config) IsDevelopment() bool {
	return c.Mode == Development
}

func (c *Config) IsProduction() bool {
	return c.Mode == Production
}

func (c *Config) IsTesting() bool {
	return c.Mode == Testing
}

func (c *Config) GetDatabaseDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		c.Database.Host, c.Database.User, c.Database.Password, c.Database.Database, c.Database.Port, c.Database.SSLMode, c.Database.TimeZone)
}
