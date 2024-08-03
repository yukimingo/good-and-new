package repositories

import (
	"errors"
	"good-and-new/models"

	"gorm.io/gorm"
)

type IAuthRepository interface {
	FindAll() (*[]models.User, error)
	FindUser(email string) (*models.User, error)
	FindUserById(id uint64) (*models.User, error)
	Create(user models.User) error
	DeleteUser(id uint64) error
}

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) IAuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) FindAll() (*[]models.User, error) {
	var users []models.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return &users, nil
}

func (r *AuthRepository) FindUser(email string) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, "email = ?", email)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("User not found")
		}
		return nil, result.Error
	}

	return &user, nil
}

func (r *AuthRepository) FindUserById(id uint64) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("User not found")
		}
		return nil, result.Error
	}

	return &user, nil
}

func (r *AuthRepository) Create(user models.User) error {
	if err := r.db.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func (r *AuthRepository) DeleteUser(id uint64) error {
	user, err := r.FindUserById(id)
	if err != nil {
		return err
	}

	if err := r.db.Delete(&user).Error; err != nil {
		return err
	}
	return nil
}
