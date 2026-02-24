package services

import (
	"context"

	"copier/database/repositories"
	"copier/internal/database/models"

	"github.com/google/uuid"
)

// PackageService defines package business logic operations
type PackageService interface {
	GetAllPackages(ctx context.Context, skip, limit int) ([]*models.Package, int64, error)
	GetActivePackages(ctx context.Context) ([]*models.Package, error)
	GetPackageByID(ctx context.Context, id uuid.UUID) (*models.Package, error)
	CreatePackage(ctx context.Context, pkg *models.Package) (*models.Package, error)
	UpdatePackage(ctx context.Context, id uuid.UUID, update *models.Package) (*models.Package, error)
	DeletePackage(ctx context.Context, id uuid.UUID) error
}

type packageService struct {
	packageRepo repositories.PackageRepository
}

// NewPackageService creates a new package service instance
func NewPackageService(packageRepo repositories.PackageRepository) PackageService {
	return &packageService{
		packageRepo: packageRepo,
	}
}

// GetAllPackages retrieves all packages with pagination
func (s *packageService) GetAllPackages(ctx context.Context, skip, limit int) ([]*models.Package, int64, error) {
	packages, err := s.packageRepo.FindAll(ctx, skip, limit)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.packageRepo.CountPackages(ctx)
	if err != nil {
		return nil, 0, err
	}

	return packages, total, nil
}

// GetActivePackages retrieves all packages that are currently active
func (s *packageService) GetActivePackages(ctx context.Context) ([]*models.Package, error) {
	return s.packageRepo.FindActivePackages(ctx)
}

// GetPackageByID retrieves a package by its UUID
func (s *packageService) GetPackageByID(ctx context.Context, id uuid.UUID) (*models.Package, error) {
	return s.packageRepo.FindByIDTyped(ctx, id)
}

// CreatePackage creates a new package
func (s *packageService) CreatePackage(ctx context.Context, pkg *models.Package) (*models.Package, error) {
	err := s.packageRepo.CreatePackage(ctx, pkg)
	if err != nil {
		return nil, err
	}
	return pkg, nil
}

// UpdatePackage updates an existing package
func (s *packageService) UpdatePackage(ctx context.Context, id uuid.UUID, update *models.Package) (*models.Package, error) {
	err := s.packageRepo.UpdatePackage(ctx, id, update)
	if err != nil {
		return nil, err
	}
	return s.packageRepo.FindByIDTyped(ctx, id)
}

// DeletePackage removes a package
func (s *packageService) DeletePackage(ctx context.Context, id uuid.UUID) error {
	return s.packageRepo.DeletePackage(ctx, id)
}
