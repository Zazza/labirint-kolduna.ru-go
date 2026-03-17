package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"gamebook-backend/pkg/utils"
	"gorm.io/gorm"
)

const (
	ErrCodeNotFound     = "NOT_FOUND"
	ErrCodeDuplicate    = "DUPLICATE"
	ErrCodeValidation   = "VALIDATION"
	ErrCodeDatabase     = "DATABASE"
	ErrCodeConflict     = "CONFLICT"
	ErrCodeBadRequest   = "BAD_REQUEST"
	ErrCodeUnauthorized = "UNAUTHORIZED"
	ErrCodeForbidden    = "FORBIDDEN"
)

type BaseRepository struct {
	db *gorm.DB
}

func NewBaseRepository(db *gorm.DB) *BaseRepository {
	return &BaseRepository{
		db: db,
	}
}

func (r *BaseRepository) GetDB(db *gorm.DB) *gorm.DB {
	if db != nil {
		return db
	}
	return r.db
}

func (r *BaseRepository) WithContext(ctx context.Context, db *gorm.DB) *gorm.DB {
	return r.GetDB(db).WithContext(ctx)
}

func (r *BaseRepository) HandleNotFound(err error, entity interface{}) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return utils.NewAppError(ErrCodeNotFound, "record not found", err)
	}
	return err
}

func (r *BaseRepository) CreateEntity(ctx context.Context, db *gorm.DB, entity interface{}) error {
	return r.WithContext(ctx, db).Create(entity).Error
}

func (r *BaseRepository) UpdateEntity(ctx context.Context, db *gorm.DB, entity interface{}) error {
	return r.WithContext(ctx, db).Save(entity).Error
}

func (r *BaseRepository) DeleteEntity(ctx context.Context, db *gorm.DB, entity interface{}) error {
	return r.WithContext(ctx, db).Delete(entity).Error
}

func (r *BaseRepository) FindEntity(ctx context.Context, db *gorm.DB, dest interface{}, conds ...interface{}) error {
	result := r.WithContext(ctx, db).First(dest, conds...)
	if result.Error != nil {
		return r.HandleNotFound(result.Error, dest)
	}
	return nil
}

func (r *BaseRepository) FindEntities(ctx context.Context, db *gorm.DB, dest interface{}, conds ...interface{}) error {
	return r.WithContext(ctx, db).Find(dest, conds...).Error
}

func (r *BaseRepository) Transaction(fc func(tx *gorm.DB) error) error {
	return r.db.Transaction(fc)
}

func (r *BaseRepository) TransactionWithContext(ctx context.Context, fc func(tx *gorm.DB) error) error {
	return r.db.WithContext(ctx).Transaction(fc)
}

func (r *BaseRepository) NullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{String: "", Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

func (r *BaseRepository) NullInt64(i int64) sql.NullInt64 {
	if i == 0 {
		return sql.NullInt64{Int64: 0, Valid: false}
	}
	return sql.NullInt64{Int64: i, Valid: true}
}

func (r *BaseRepository) NullTime(t time.Time) sql.NullTime {
	if t.IsZero() {
		return sql.NullTime{Time: t, Valid: false}
	}
	return sql.NullTime{Time: t, Valid: true}
}

func (r *BaseRepository) NullUUID(id interface{}) sql.NullString {
	switch v := id.(type) {
	case string:
		if v == "" {
			return sql.NullString{String: "", Valid: false}
		}
		return sql.NullString{String: v, Valid: true}
	case nil:
		return sql.NullString{String: "", Valid: false}
	default:
		return sql.NullString{String: "", Valid: false}
	}
}
