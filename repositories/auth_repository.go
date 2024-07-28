package repositories

import (
	"good-and-new/models"

	"gorm.io/gorm"
)

type IAuthRepository interface {
	Create(user models.User) error
}

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) IAuthRepository {
	return AuthRepository{db: db}
}

func (r AuthRepository) Create(user models.User) error {
	if err := r.db.Create(&user).Error; err != nil {
		return err
	}
	return nil
}
