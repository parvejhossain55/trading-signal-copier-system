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

// UserRepository defines user-specific repository operations
type UserRepository interface {
	BaseRepository
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByPhone(ctx context.Context, phone string) (*User, error)
	FindBySession(ctx context.Context, session string) (*User, error)
	UpdateSession(ctx context.Context, userID primitive.ObjectID, session *string) error
	FindByIDTyped(ctx context.Context, id primitive.ObjectID) (*User, error)
	CreateUser(ctx context.Context, user *User) (*mongo.InsertOneResult, error)
	UpdateUser(ctx context.Context, id primitive.ObjectID, update *User) error
	DeleteUser(ctx context.Context, id primitive.ObjectID) error
	FindAll(ctx context.Context, skip, limit int64) ([]*User, error)
	CountUsers(ctx context.Context) (int64, error)
}

// userRepository implements UserRepository interface
type userRepository struct {
	BaseRepository
	collection *mongo.Collection
}

// NewUserRepository creates a new user repository instance
func NewUserRepository(collection *mongo.Collection) UserRepository {
	return &userRepository{
		BaseRepository: NewBaseRepository(collection),
		collection:     collection,
	}
}

// FindByEmail finds a user by email address
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found with email: %s", email)
		}
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	return &user, nil
}

// FindByPhone finds a user by phone number
func (r *userRepository) FindByPhone(ctx context.Context, phone string) (*User, error) {
	var user User
	err := r.collection.FindOne(ctx, bson.M{"phone": phone}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found with phone: %s", phone)
		}
		return nil, fmt.Errorf("failed to find user by phone: %w", err)
	}

	return &user, nil
}

// FindBySession finds a user by session token
func (r *userRepository) FindBySession(ctx context.Context, session string) (*User, error) {
	var user User
	err := r.collection.FindOne(ctx, bson.M{"session": session}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found with session")
		}
		return nil, fmt.Errorf("failed to find user by session: %w", err)
	}

	return &user, nil
}

// UpdateSession updates user's session token
func (r *userRepository) UpdateSession(ctx context.Context, userID primitive.ObjectID, session *string) error {
	update := bson.M{
		"session":    session,
		"updated_at": bson.M{"$currentDate": true},
	}

	_, err := r.collection.UpdateByID(ctx, userID, bson.M{"$set": update})
	if err != nil {
		return fmt.Errorf("failed to update user session: %w", err)
	}

	return nil
}

// FindByIDTyped finds a user by ID and returns typed User struct
func (r *userRepository) FindByIDTyped(ctx context.Context, id primitive.ObjectID) (*User, error) {
	var user User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found with ID: %s", id.Hex())
		}
		return nil, fmt.Errorf("failed to find user by ID: %w", err)
	}

	return &user, nil
}

// CreateUser creates a new user
func (r *userRepository) CreateUser(ctx context.Context, user *User) (*mongo.InsertOneResult, error) {
	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return result, nil
}

// UpdateUser updates an existing user
func (r *userRepository) UpdateUser(ctx context.Context, id primitive.ObjectID, update *User) error {
	update.UpdatedAt = time.Now()
	
	_, err := r.collection.UpdateByID(ctx, id, bson.M{"$set": update})
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// DeleteUser deletes a user by ID
func (r *userRepository) DeleteUser(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// FindAll retrieves all users with pagination
func (r *userRepository) FindAll(ctx context.Context, skip, limit int64) ([]*User, error) {
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
		return nil, fmt.Errorf("failed to find users: %w", err)
	}
	defer cursor.Close(ctx)

	var users []*User
	for cursor.Next(ctx) {
		var user User
		if err := cursor.Decode(&user); err != nil {
			return nil, fmt.Errorf("failed to decode user: %w", err)
		}
		users = append(users, &user)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return users, nil
}

// CountUsers returns the total number of users
func (r *userRepository) CountUsers(ctx context.Context) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	return count, nil
}
