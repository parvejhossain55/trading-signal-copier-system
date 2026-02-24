package repositories

import (
	"context"
	"fmt"

	"copier/internal/database/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRepository defines user-specific repository operations
type UserRepository interface {
	BaseRepository
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByPhone(ctx context.Context, phone string) (*models.User, error)
	FindBySession(ctx context.Context, session string) (*models.User, error)
	UpdateSession(ctx context.Context, userID uuid.UUID, session *string) error
	FindByIDTyped(ctx context.Context, id uuid.UUID) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
	UpdateUser(ctx context.Context, id uuid.UUID, update *models.User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	FindAll(ctx context.Context, skip, limit int) ([]*models.User, error)
	CountUsers(ctx context.Context) (int64, error)
}

// userRepository implements UserRepository interface
type userRepository struct {
	BaseRepository
	db *gorm.DB
}

// NewUserRepository creates a new user repository instance
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		BaseRepository: NewBaseRepository(db),
		db:             db,
	}
}

// FindByEmail finds a user by email address
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found with email: %s", email)
		}
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	return &user, nil
}

// FindByPhone finds a user by phone number
func (r *userRepository) FindByPhone(ctx context.Context, phone string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("phone = ?", phone).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found with phone: %s", phone)
		}
		return nil, fmt.Errorf("failed to find user by phone: %w", err)
	}

	return &user, nil
}

// FindBySession finds a user by session token
func (r *userRepository) FindBySession(ctx context.Context, session string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("session = ?", session).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found with session")
		}
		return nil, fmt.Errorf("failed to find user by session: %w", err)
	}

	return &user, nil
}

// UpdateSession updates user's session token
func (r *userRepository) UpdateSession(ctx context.Context, userID uuid.UUID, session *string) error {
	err := r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Update("session", session).Error
	if err != nil {
		return fmt.Errorf("failed to update user session: %w", err)
	}

	return nil
}

// FindByIDTyped finds a user by ID and returns typed User struct
func (r *userRepository) FindByIDTyped(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found with ID: %s", id)
		}
		return nil, fmt.Errorf("failed to find user by ID: %w", err)
	}

	return &user, nil
}

// CreateUser creates a new user
func (r *userRepository) CreateUser(ctx context.Context, user *models.User) error {
	err := r.db.WithContext(ctx).Create(user).Error
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// UpdateUser updates an existing user
func (r *userRepository) UpdateUser(ctx context.Context, id uuid.UUID, update *models.User) error {
	err := r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", id).Updates(update).Error
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// DeleteUser deletes a user by ID
func (r *userRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	err := r.db.WithContext(ctx).Delete(&models.User{}, id).Error
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// FindAll retrieves all users with pagination
func (r *userRepository) FindAll(ctx context.Context, skip, limit int) ([]*models.User, error) {
	var users []*models.User
	query := r.db.WithContext(ctx).Order("created_at desc")
	if skip > 0 {
		query = query.Offset(skip)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find users: %w", err)
	}

	return users, nil
}

// CountUsers returns the total number of users
func (r *userRepository) CountUsers(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.User{}).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	return count, nil
}
