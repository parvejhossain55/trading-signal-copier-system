package seed

import (
	"copier/config"
	"copier/internal/database/models"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SeedData() error {

	db, err := config.NewPostgresDB()
	if err != nil {
		return fmt.Errorf("failed to connect to PostgreSQL: %v", err)
	}

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	err = seedUsers(db)
	if err != nil {
		return fmt.Errorf("failed to seed users: %v", err)
	}

	err = seedPlatforms(db)
	if err != nil {
		return fmt.Errorf("failed to seed platforms: %v", err)
	}

	err = seedChannels(db)
	if err != nil {
		return fmt.Errorf("failed to seed channels: %v", err)
	}

	err = seedTradeSettings(db)
	if err != nil {
		return fmt.Errorf("failed to seed trade settings: %v", err)
	}

	slog.Info("Database seeding completed successfully")
	return nil
}

func seedUsers(db *gorm.DB) error {
	users := []models.User{
		{
			ID:        uuid.New(),
			Name:      "John Doe",
			Email:     "john.doe@example.com",
			Phone:     stringPtr("+1234567890"),
			Session:   stringPtr("session_123"),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			Name:      "Jane Smith",
			Email:     "jane.smith@example.com",
			Phone:     stringPtr("+0987654321"),
			Session:   stringPtr("session_456"),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	for _, user := range users {
		var existing models.User
		if err := db.Where("email = ?", user.Email).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&user).Error; err != nil {
					return fmt.Errorf("failed to create user %s: %v", user.Email, err)
				}
			} else {
				return err
			}
		}
	}

	slog.Info("Seeded users collection")
	return nil
}

func seedPackages(db *gorm.DB) error {
	packages := []models.Package{
		{
			ID:           uuid.New(),
			Name:         "Basic Plan",
			Price:        29.99,
			Type:         models.PackageTypeBasic,
			Description:  stringPtr("Basic trading features"),
			Features:     []string{"Basic signals", "Email support"},
			DurationDays: intPtr(30),
			IsActive:     true,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			ID:           uuid.New(),
			Name:         "Premium Plan",
			Price:        99.99,
			Type:         models.PackageTypePremium,
			Description:  stringPtr("Advanced trading features"),
			Features:     []string{"Advanced signals", "Priority support", "Custom strategies"},
			DurationDays: intPtr(30),
			IsActive:     true,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			ID:           uuid.New(),
			Name:         "Enterprise Plan",
			Price:        299.99,
			Type:         models.PackageTypeEnterprise,
			Description:  stringPtr("Full-featured trading platform"),
			Features:     []string{"All signals", "24/7 support", "Custom strategies", "API access"},
			DurationDays: intPtr(30),
			IsActive:     true,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}

	for _, pkg := range packages {
		var existing models.Package
		if err := db.Where("name = ?", pkg.Name).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&pkg).Error; err != nil {
					return fmt.Errorf("failed to create package %s: %v", pkg.Name, err)
				}
			} else {
				return err
			}
		}
	}

	slog.Info("Seeded packages collection")
	return nil
}

func seedPlatforms(db *gorm.DB) error {
	userID, err := GetSampleUserID(db)
	if err != nil {
		return err
	}

	platforms := []models.Platform{
		{
			ID:        uuid.New(),
			UserID:    userID,
			Name:      "MT4-Demo",
			APIKey:    "sample_api_key",
			APISecret: "sample_api_secret",
		},
	}

	for _, platform := range platforms {
		var existing models.Platform
		if err := db.Where("name = ? AND user_id = ?", platform.Name, userID).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				db.Create(&platform)
			}
		}
	}

	slog.Info("Seeded platforms")
	return nil
}

func seedChannels(db *gorm.DB) error {
	userID, err := GetSampleUserID(db)
	if err != nil {
		return err
	}

	channels := []models.Channel{
		{
			ID:        uuid.New(),
			UserID:    userID,
			Name:      "Telegram Signals",
			ChannelID: "-100123456789",
		},
	}

	for _, channel := range channels {
		var existing models.Channel
		if err := db.Where("name = ? AND user_id = ?", channel.Name, userID).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				db.Create(&channel)
			}
		}
	}

	slog.Info("Seeded channels")
	return nil
}

func seedTradeSettings(db *gorm.DB) error {
	userID, err := GetSampleUserID(db)
	if err != nil {
		return err
	}

	settings := models.TradeSettings{
		ID:                  uuid.New(),
		UserID:              userID,
		PerTradeAmount:      1000.0,
		StopLossStatus:      true,
		StopLossPercentage:  5,
		OverRangePercentage: 2.0,
		TakeProfitStatus:    true,
		TakeProfitStep:      3,
		TPPercentage:        []float64{1.5, 3.0, 5.0},
	}

	var existing models.TradeSettings
	if err := db.Where("user_id = ?", userID).First(&existing).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			db.Create(&settings)
		}
	}

	slog.Info("Seeded trade settings")
	return nil
}

func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func GetSampleUserID(db *gorm.DB) (uuid.UUID, error) {
	var user models.User
	err := db.First(&user).Error
	if err != nil {
		return uuid.Nil, err
	}
	return user.ID, nil
}

func GetSamplePackageID(db *gorm.DB) (uuid.UUID, error) {
	var pkg models.Package
	err := db.First(&pkg).Error
	if err != nil {
		return uuid.Nil, err
	}
	return pkg.ID, nil
}
