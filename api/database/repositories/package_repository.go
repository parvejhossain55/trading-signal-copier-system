package repositories

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// PackageRepository defines package-specific repository operations
type PackageRepository interface {
	BaseRepository
	FindByName(ctx context.Context, name string) (*Package, error)
	FindByType(ctx context.Context, packageType PackageType) ([]*Package, error)
	FindActivePackages(ctx context.Context) ([]*Package, error)
	FindByTypeAndActive(ctx context.Context, packageType PackageType, isActive bool) ([]*Package, error)
	FindByPriceRange(ctx context.Context, minPrice, maxPrice float64) ([]*Package, error)
	FindByIDTyped(ctx context.Context, id primitive.ObjectID) (*Package, error)
	CreatePackage(ctx context.Context, pkg *Package) (*mongo.InsertOneResult, error)
	UpdatePackage(ctx context.Context, id primitive.ObjectID, update *Package) error
	DeletePackage(ctx context.Context, id primitive.ObjectID) error
	FindAll(ctx context.Context, skip, limit int64) ([]*Package, error)
	CountPackages(ctx context.Context) (int64, error)
	CountByType(ctx context.Context, packageType PackageType) (int64, error)
	CountActivePackages(ctx context.Context) (int64, error)
	UpdateActiveStatus(ctx context.Context, id primitive.ObjectID, isActive bool) error
	UpdatePrice(ctx context.Context, id primitive.ObjectID, price float64) error
}

// packageRepository implements PackageRepository interface
type packageRepository struct {
	BaseRepository
	collection *mongo.Collection
}

// NewPackageRepository creates a new package repository instance
func NewPackageRepository(collection *mongo.Collection) PackageRepository {
	return &packageRepository{
		BaseRepository: NewBaseRepository(collection),
		collection:     collection,
	}
}

// FindByName finds a package by name
func (r *packageRepository) FindByName(ctx context.Context, name string) (*Package, error) {
	var pkg Package
	err := r.collection.FindOne(ctx, bson.M{"name": name}).Decode(&pkg)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("package not found with name: %s", name)
		}
		return nil, fmt.Errorf("failed to find package by name: %w", err)
	}

	return &pkg, nil
}

// FindByType finds all packages of a specific type
func (r *packageRepository) FindByType(ctx context.Context, packageType PackageType) ([]*Package, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"type": packageType})
	if err != nil {
		return nil, fmt.Errorf("failed to find packages by type: %w", err)
	}
	defer cursor.Close(ctx)

	var packages []*Package
	for cursor.Next(ctx) {
		var pkg Package
		if err := cursor.Decode(&pkg); err != nil {
			return nil, fmt.Errorf("failed to decode package: %w", err)
		}
		packages = append(packages, &pkg)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return packages, nil
}

// FindActivePackages finds all active packages
func (r *packageRepository) FindActivePackages(ctx context.Context) ([]*Package, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"is_active": true})
	if err != nil {
		return nil, fmt.Errorf("failed to find active packages: %w", err)
	}
	defer cursor.Close(ctx)

	var packages []*Package
	for cursor.Next(ctx) {
		var pkg Package
		if err := cursor.Decode(&pkg); err != nil {
			return nil, fmt.Errorf("failed to decode package: %w", err)
		}
		packages = append(packages, &pkg)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return packages, nil
}

// FindByTypeAndActive finds packages by type and active status
func (r *packageRepository) FindByTypeAndActive(ctx context.Context, packageType PackageType, isActive bool) ([]*Package, error) {
	filter := bson.M{
		"type":      packageType,
		"is_active": isActive,
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find packages by type and active status: %w", err)
	}
	defer cursor.Close(ctx)

	var packages []*Package
	for cursor.Next(ctx) {
		var pkg Package
		if err := cursor.Decode(&pkg); err != nil {
			return nil, fmt.Errorf("failed to decode package: %w", err)
		}
		packages = append(packages, &pkg)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return packages, nil
}

// FindByPriceRange finds packages within a price range
func (r *packageRepository) FindByPriceRange(ctx context.Context, minPrice, maxPrice float64) ([]*Package, error) {
	filter := bson.M{
		"price": bson.M{
			"$gte": minPrice,
			"$lte": maxPrice,
		},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find packages by price range: %w", err)
	}
	defer cursor.Close(ctx)

	var packages []*Package
	for cursor.Next(ctx) {
		var pkg Package
		if err := cursor.Decode(&pkg); err != nil {
			return nil, fmt.Errorf("failed to decode package: %w", err)
		}
		packages = append(packages, &pkg)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return packages, nil
}

// FindByIDTyped finds a package by ID and returns typed Package struct
func (r *packageRepository) FindByIDTyped(ctx context.Context, id primitive.ObjectID) (*Package, error) {
	var pkg Package
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&pkg)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("package not found with ID: %s", id.Hex())
		}
		return nil, fmt.Errorf("failed to find package by ID: %w", err)
	}

	return &pkg, nil
}

// CreatePackage creates a new package
func (r *packageRepository) CreatePackage(ctx context.Context, pkg *Package) (*mongo.InsertOneResult, error) {
	result, err := r.collection.InsertOne(ctx, pkg)
	if err != nil {
		return nil, fmt.Errorf("failed to create package: %w", err)
	}

	return result, nil
}

// UpdatePackage updates an existing package
func (r *packageRepository) UpdatePackage(ctx context.Context, id primitive.ObjectID, update *Package) error {
	update.UpdatedAt = time.Now()
	
	_, err := r.collection.UpdateByID(ctx, id, bson.M{"$set": update})
	if err != nil {
		return fmt.Errorf("failed to update package: %w", err)
	}

	return nil
}

// DeletePackage deletes a package by ID
func (r *packageRepository) DeletePackage(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete package: %w", err)
	}

	return nil
}

// FindAll retrieves all packages with pagination
func (r *packageRepository) FindAll(ctx context.Context, skip, limit int64) ([]*Package, error) {
	opts := options.Find()
	if skip > 0 {
		opts.SetSkip(skip)
	}
	if limit > 0 {
		opts.SetLimit(limit)
	}
	opts.SetSort(bson.M{"created_at": -1}) // Sort by creation date descending

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find packages: %w", err)
	}
	defer cursor.Close(ctx)

	var packages []*Package
	for cursor.Next(ctx) {
		var pkg Package
		if err := cursor.Decode(&pkg); err != nil {
			return nil, fmt.Errorf("failed to decode package: %w", err)
		}
		packages = append(packages, &pkg)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return packages, nil
}

// CountPackages returns the total number of packages
func (r *packageRepository) CountPackages(ctx context.Context) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, fmt.Errorf("failed to count packages: %w", err)
	}

	return count, nil
}

// CountByType returns the number of packages of a specific type
func (r *packageRepository) CountByType(ctx context.Context, packageType PackageType) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"type": packageType})
	if err != nil {
		return 0, fmt.Errorf("failed to count packages by type: %w", err)
	}

	return count, nil
}

// CountActivePackages returns the number of active packages
func (r *packageRepository) CountActivePackages(ctx context.Context) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"is_active": true})
	if err != nil {
		return 0, fmt.Errorf("failed to count active packages: %w", err)
	}

	return count, nil
}

// UpdateActiveStatus updates only the active status of a package
func (r *packageRepository) UpdateActiveStatus(ctx context.Context, id primitive.ObjectID, isActive bool) error {
	update := bson.M{
		"is_active":  isActive,
		"updated_at": bson.M{"$currentDate": true},
	}

	_, err := r.collection.UpdateByID(ctx, id, bson.M{"$set": update})
	if err != nil {
		return fmt.Errorf("failed to update package active status: %w", err)
	}

	return nil
}

// UpdatePrice updates only the price of a package
func (r *packageRepository) UpdatePrice(ctx context.Context, id primitive.ObjectID, price float64) error {
	update := bson.M{
		"price":      price,
		"updated_at": bson.M{"$currentDate": true},
	}

	_, err := r.collection.UpdateByID(ctx, id, bson.M{"$set": update})
	if err != nil {
		return fmt.Errorf("failed to update package price: %w", err)
	}

	return nil
}
