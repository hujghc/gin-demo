package services

import (
	"gin-demo/models"
	"gin-demo/repositories"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type UserService struct {
	userRepo *repositories.UserRepository
	db       *gorm.DB
}

func NewUserService(userRepo *repositories.UserRepository, db *gorm.DB) *UserService {
	return &UserService{
		userRepo: userRepo,
		db:       db,
	}
}

func (s *UserService) CreateUser(user *models.User) error {
	return s.userRepo.Create(user)
}

func (s *UserService) GetUser(id uint) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *UserService) GetUsers(page, limit int) ([]models.User, int64, error) {
	return s.userRepo.FindAll(page, limit)
}

func (s *UserService) UpdateUser(id uint, user *models.User) error {
	_, err := s.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("用户不存在")
		}
		return err
	}
	user.ID = id
	return s.userRepo.Update(user)
}

func (s *UserService) DeleteUser(id uint) error {
	err := s.userRepo.Delete(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("用户不存在")
	}
	return err
}

func (s *UserService) TransferBalance(fromUserID, toUserID uint, amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("转账金额必须大于0")
	}
	if fromUserID == toUserID {
		return fmt.Errorf("不能向自己转账")
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		userRepo := repositories.NewUserRepository(tx)

		_, err := userRepo.FindByIDWithLock(fromUserID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("转出账户不存在")
			}
			return err
		}

		_, err = userRepo.FindByIDWithLock(toUserID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("转入账户不存在")
			}
			return err
		}

		if err := userRepo.UpdateBalanceWithCondition(fromUserID, -amount, amount); err != nil {
			return fmt.Errorf("扣款失败")
		}

		if err := userRepo.UpdateBalance(toUserID, amount); err != nil {
			return fmt.Errorf("入账失败")
		}

		return nil
	})

	return err
}
