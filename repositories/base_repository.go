package repositories

import (
	"gorm.io/gorm"
)

type BaseRepository[T any] struct {
	DB *gorm.DB
}

func NewBaseRepository[T any](db *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{DB: db}
}

func (r *BaseRepository[T]) Create(entity *T) error {
	return r.DB.Create(entity).Error
}

func (r *BaseRepository[T]) FindByID(id uint) (*T, error) {
	var entity T
	err := r.DB.First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *BaseRepository[T]) FindAll(page, limit int) ([]T, int64, error) {
	var entities []T
	var total int64

	offset := (page - 1) * limit

	result := r.DB.Model(new(T)).Count(&total)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	err := r.DB.Offset(offset).Limit(limit).Find(&entities).Error
	if err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

func (r *BaseRepository[T]) Update(entity *T) error {
	return r.DB.Save(entity).Error
}

func (r *BaseRepository[T]) Delete(id uint) error {
	result := r.DB.Delete(new(T), id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *BaseRepository[T]) Count() (int64, error) {
	var count int64
	err := r.DB.Model(new(T)).Count(&count).Error
	return count, err
}
