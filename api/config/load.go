package config

import (
	"log/slog"
	"os"

	"github.com/go-playground/validator"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func LoadConfig() error {
	exit := func(err error) {
		slog.Error(err.Error())
		os.Exit(1)
	}

	_ = godotenv.Load()

	viper.SetDefault("VERSION", "1.0.0")
	viper.SetDefault("MODE", "development")
	viper.SetDefault("SERVICE_NAME", "copier")
	viper.SetDefault("HTTP_PORT", 8080)
	viper.SetDefault("GRPC_PORT", 8081)
	viper.SetDefault("HEALTH_CHECK_ROUTE", "/health")
	viper.SetDefault("API_VERSION", "v1")
	viper.SetDefault("JWT_SECRET", "secret")
	viper.SetDefault("SERVICE_BASE_PATH", "/api/v1")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", 5432)
	viper.SetDefault("DB_USER", "admin")
	viper.SetDefault("DB_PASSWORD", "password")
	viper.SetDefault("DB_NAME", "copier_db")
	viper.SetDefault("DB_SSL_MODE", "disable")
	viper.SetDefault("DB_TIMEZONE", "UTC")
	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", 6379)

	viper.AutomaticEnv()

	var err error

	config = &Config{
		Version:          viper.GetString("VERSION"),
		Mode:             Mode(viper.GetString("MODE")),
		ServiceName:      viper.GetString("SERVICE_NAME"),
		HttpPort:         viper.GetInt("HTTP_PORT"),
		GrpcPort:         viper.GetInt("GRPC_PORT"),
		HealthCheckRoute: viper.GetString("HEALTH_CHECK_ROUTE"),
		ApiVersion:       viper.GetString("API_VERSION"),
		JwtSecret:        viper.GetString("JWT_SECRET"),
		ServiceBasePath:  viper.GetString("SERVICE_BASE_PATH"),

		Redis: RedisConfig{
			Host:     viper.GetString("REDIS_HOST"),
			Port:     viper.GetInt("REDIS_PORT"),
			Password: viper.GetString("REDIS_PASSWORD"),
			DB:       viper.GetInt("REDIS_DB"),
			PoolSize: viper.GetInt("REDIS_POOL_SIZE"),
			MinIdle:  viper.GetInt("REDIS_MIN_IDLE"),
		},

		Database: PostgresDatabase{
			Host:            viper.GetString("DB_HOST"),
			Port:            viper.GetInt("DB_PORT"),
			User:            viper.GetString("DB_USER"),
			Password:        viper.GetString("DB_PASSWORD"),
			Database:        viper.GetString("DB_NAME"),
			SSLMode:         viper.GetString("DB_SSL_MODE"),
			TimeZone:        viper.GetString("DB_TIMEZONE"),
			MaxIdleConns:    viper.GetInt("DB_MAX_IDLE_CONNS"),
			MaxOpenConns:    viper.GetInt("DB_MAX_OPEN_CONNS"),
			ConnMaxLifetime: viper.GetInt("DB_CONN_MAX_LIFETIME"),
		},

		Cors: CorsOrigin{
			Origin: viper.GetStringSlice("CORS_ORIGIN"),
		},
	}

	v := validator.New()
	if err = v.Struct(config); err != nil {
		exit(err)
	}

	return nil
}
