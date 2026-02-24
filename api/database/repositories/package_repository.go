package repositories

import (
	"context"
	"fmt"

	"copier/internal/database/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PackageRepository defines package-specific repository operations
type PackageRepository interface {
	BaseRepository
	FindByName(ctx context.Context, name string) (*models.Package, error)
	FindByType(ctx context.Context, packageType models.PackageType) ([]*models.Package, error)
	FindActivePackages(ctx context.Context) ([]*models.Package, error)
	FindByTypeAndActive(ctx context.Context, packageType models.PackageType, isActive bool) ([]*models.Package, error)
	FindByPriceRange(ctx context.Context, minPrice, maxPrice float64) ([]*models.Package, error)
	FindByIDTyped(ctx context.Context, id uuid.UUID) (*models.Package, error)
	CreatePackage(ctx context.Context, pkg *models.Package) error
	UpdatePackage(ctx context.Context, id uuid.UUID, update *models.Package) error
	DeletePackage(ctx context.Context, id uuid.UUID) error
	FindAll(ctx context.Context, skip, limit int) ([]*models.Package, error)
	CountPackages(ctx context.Context) (int64, error)
	CountByType(ctx context.Context, packageType models.PackageType) (int64, error)
	CountActivePackages(ctx context.Context) (int64, error)
	UpdateActiveStatus(ctx context.Context, id uuid.UUID, isActive bool) error
	UpdatePrice(ctx context.Context, id uuid.UUID, price float64) error
}

// packageRepository implements PackageRepository interface
type packageRepository struct {
	BaseRepository
	db *gorm.DB
}

// NewPackageRepository creates a new package repository instance
func NewPackageRepository(db *gorm.DB) PackageRepository {
	return &packageRepository{
		BaseRepository: NewBaseRepository(db),
		db:             db,
	}
}

// FindByName finds a package by name
func (r *packageRepository) FindByName(ctx context.Context, name string) (*models.Package, error) {
	var pkg models.Package
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&pkg).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("package not found with name: %s", name)
		}
		return nil, fmt.Errorf("failed to find package by name: %w", err)
	}

	return &pkg, nil
}

// FindByType finds all packages of a specific type
func (r *packageRepository) FindByType(ctx context.Context, packageType models.PackageType) ([]*models.Package, error) {
	var packages []*models.Package
	err := r.db.WithContext(ctx).Where("type = ?", packageType).Find(&packages).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find packages by type: %w", err)
	}

	return packages, nil
}

// FindActivePackages finds all active packages
func (r *packageRepository) FindActivePackages(ctx context.Context) ([]*models.Package, error) {
	var packages []*models.Package
	err := r.db.WithContext(ctx).Where("is_active = ?", true).Find(&packages).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find active packages: %w", err)
	}

	return packages, nil
}

// FindByTypeAndActive finds packages by type and active status
func (r *packageRepository) FindByTypeAndActive(ctx context.Context, packageType models.PackageType, isActive bool) ([]*models.Package, error) {
	var packages []*models.Package
	err := r.db.WithContext(ctx).Where("type = ? AND is_active = ?", packageType, isActive).Find(&packages).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find packages by type and active status: %w", err)
	}

	return packages, nil
}

// FindByPriceRange finds packages within a price range
func (r *packageRepository) FindByPriceRange(ctx context.Context, minPrice, maxPrice float64) ([]*models.Package, error) {
	var packages []*models.Package
	err := r.db.WithContext(ctx).Where("price >= ? AND price <= ?", minPrice, maxPrice).Find(&packages).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find packages by price range: %w", err)
	}

	return packages, nil
}

// FindByIDTyped finds a package by ID and returns typed Package struct
func (r *packageRepository) FindByIDTyped(ctx context.Context, id uuid.UUID) (*models.Package, error) {
	var pkg models.Package
	err := r.db.WithContext(ctx).First(&pkg, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("package not found with ID: %s", id)
		}
		return nil, fmt.Errorf("failed to find package by ID: %w", err)
	}

	return &pkg, nil
}

// CreatePackage creates a new package
func (r *packageRepository) CreatePackage(ctx context.Context, pkg *models.Package) error {
	err := r.db.WithContext(ctx).Create(pkg).Error
	if err != nil {
		return fmt.Errorf("failed to create package: %w", err)
	}

	return nil
}

// UpdatePackage updates an existing package
func (r *packageRepository) UpdatePackage(ctx context.Context, id uuid.UUID, update *models.Package) error {
	err := r.db.WithContext(ctx).Model(&models.Package{}).Where("id = ?", id).Updates(update).Error
	if err != nil {
		return fmt.Errorf("failed to update package: %w", err)
	}

	return nil
}

// DeletePackage deletes a package by ID
func (r *packageRepository) DeletePackage(ctx context.Context, id uuid.UUID) error {
	err := r.db.WithContext(ctx).Delete(&models.Package{}, id).Error
	if err != nil {
		return fmt.Errorf("failed to delete package: %w", err)
	}

	return nil
}

// FindAll retrieves all packages with pagination
func (r *packageRepository) FindAll(ctx context.Context, skip, limit int) ([]*models.Package, error) {
	var packages []*models.Package
	query := r.db.WithContext(ctx).Order("created_at desc")
	if skip > 0 {
		query = query.Offset(skip)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&packages).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find packages: %w", err)
	}

	return packages, nil
}

// CountPackages returns the total number of packages
func (r *packageRepository) CountPackages(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Package{}).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count packages: %w", err)
	}

	return count, nil
}

// CountByType returns the number of packages of a specific type
func (r *packageRepository) CountByType(ctx context.Context, packageType models.PackageType) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Package{}).Where("type = ?", packageType).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count packages by type: %w", err)
	}

	return count, nil
}

// CountActivePackages returns the number of active packages
func (r *packageRepository) CountActivePackages(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Package{}).Where("is_active = ?", true).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count active packages: %w", err)
	}

	return count, nil
}

// UpdateActiveStatus updates only the active status of a package
func (r *packageRepository) UpdateActiveStatus(ctx context.Context, id uuid.UUID, isActive bool) error {
	err := r.db.WithContext(ctx).Model(&models.Package{}).Where("id = ?", id).Update("is_active", isActive).Error
	if err != nil {
		return fmt.Errorf("failed to update package active status: %w", err)
	}

	return nil
}

// UpdatePrice updates only the price of a package
func (r *packageRepository) UpdatePrice(ctx context.Context, id uuid.UUID, price float64) error {
	err := r.db.WithContext(ctx).Model(&models.Package{}).Where("id = ?", id).Update("price", price).Error
	if err != nil {
		return fmt.Errorf("failed to update package price: %w", err)
	}

	return nil
}
