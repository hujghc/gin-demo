package repositories

import (
	"gin-demo/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	*BaseRepository[models.User]
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		BaseRepository: NewBaseRepository[models.User](db),
		DB:             db,
	}
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
