package repositories

import (
	"gin-demo/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(user *models.User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := r.DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindAll(page, limit int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	offset := (page - 1) * limit

	result := r.DB.Model(&models.User{}).Count(&total)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	err := r.DB.Offset(offset).Limit(limit).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *UserRepository) Update(user *models.User) error {
	return r.DB.Save(user).Error
}

func (r *UserRepository) Delete(id uint) error {
	result := r.DB.Delete(&models.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *UserRepository) FindByIDWithLock(id uint) (*models.User, error) {
	var user models.User
	err := r.DB.Raw("SELECT * FROM user WHERE id = ? FOR UPDATE", id).Scan(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateBalance(id uint, amount float64) error {
	return r.DB.Model(&models.User{}).
		Where("id = ?", id).
		Update("balance", gorm.Expr("balance + ?", amount)).Error
}

func (r *UserRepository) UpdateBalanceWithCondition(id uint, amount float64, minBalance float64) error {
	return r.DB.Model(&models.User{}).
		Where("id = ? AND balance >= ?", id, minBalance).
		Update("balance", gorm.Expr("balance + ?", amount)).Error
}
