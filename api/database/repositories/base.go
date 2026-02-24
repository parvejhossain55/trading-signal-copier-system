package repositories

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BaseRepository defines common repository operations
type BaseRepository interface {
	Create(ctx context.Context, model interface{}) error
	First(ctx context.Context, out interface{}, conds ...interface{}) error
	Find(ctx context.Context, out interface{}, conds ...interface{}) error
	Save(ctx context.Context, model interface{}) error
	Updates(ctx context.Context, model interface{}, values interface{}) error
	Delete(ctx context.Context, model interface{}, conds ...interface{}) error
	DeleteByID(ctx context.Context, model interface{}, id uuid.UUID) error
	Count(ctx context.Context, model interface{}, count *int64) error
	WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error
}

// baseRepository implements common repository operations
type baseRepository struct {
	db *gorm.DB
}

// NewBaseRepository creates a new base repository instance
func NewBaseRepository(db *gorm.DB) BaseRepository {
	return &baseRepository{
		db: db,
	}
}

func (r *baseRepository) Create(ctx context.Context, model interface{}) error {
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *baseRepository) First(ctx context.Context, out interface{}, conds ...interface{}) error {
	return r.db.WithContext(ctx).First(out, conds...).Error
}

func (r *baseRepository) Find(ctx context.Context, out interface{}, conds ...interface{}) error {
	return r.db.WithContext(ctx).Find(out, conds...).Error
}

func (r *baseRepository) Save(ctx context.Context, model interface{}) error {
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *baseRepository) Updates(ctx context.Context, model interface{}, values interface{}) error {
	return r.db.WithContext(ctx).Model(model).Updates(values).Error
}

func (r *baseRepository) Delete(ctx context.Context, model interface{}, conds ...interface{}) error {
	return r.db.WithContext(ctx).Delete(model, conds...).Error
}

func (r *baseRepository) DeleteByID(ctx context.Context, model interface{}, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(model, id).Error
}

func (r *baseRepository) Count(ctx context.Context, model interface{}, count *int64) error {
	return r.db.WithContext(ctx).Model(model).Count(count).Error
}

func (r *baseRepository) WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return r.db.WithContext(ctx).Transaction(fn)
}
