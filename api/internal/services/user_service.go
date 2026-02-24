package services

import (
	"context"

	"copier/database/repositories"
	"copier/internal/database/models"

	"github.com/google/uuid"
)

// UserService defines user business logic operations
type UserService interface {
	RegisterUser(ctx context.Context, name, email string, phone *string) (*models.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, update *models.User) (*models.User, error)
	GetAllUsers(ctx context.Context, skip, limit int) ([]*models.User, int64, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type userService struct {
	userRepo repositories.UserRepository
}

// NewUserService creates a new user service instance
func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// RegisterUser handles the business logic for registering a new user
func (s *userService) RegisterUser(ctx context.Context, name, email string, phone *string) (*models.User, error) {
	// Note: In a real app, we would also check for existing entries and use business-specific validations.
	// We'll rely on repository errors for now (e.g., unique email constraint) to keep it simple but functional.

	user := &models.User{
		Name:  name,
		Email: email,
		Phone: phone,
	}

	err := s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByID retrieves a user by their UUID
func (s *userService) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return s.userRepo.FindByIDTyped(ctx, id)
}

// GetUserByEmail retrieves a user by their email address
func (s *userService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.userRepo.FindByEmail(ctx, email)
}

// UpdateUser handles user profile updates
func (s *userService) UpdateUser(ctx context.Context, id uuid.UUID, update *models.User) (*models.User, error) {
	err := s.userRepo.UpdateUser(ctx, id, update)
	if err != nil {
		return nil, err
	}
	return s.userRepo.FindByIDTyped(ctx, id)
}

// GetAllUsers retrieves a paginated list of users and the total count
func (s *userService) GetAllUsers(ctx context.Context, skip, limit int) ([]*models.User, int64, error) {
	users, err := s.userRepo.FindAll(ctx, skip, limit)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.userRepo.CountUsers(ctx)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// DeleteUser removes a user from the system
func (s *userService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.userRepo.DeleteUser(ctx, id)
}
