package usecases

import (
	"good-and-new/dto"
	"good-and-new/models"
	"good-and-new/repositories"
)

type IAuthUsecase interface {
	FindAll() (*[]models.User, error)
	FindUser(email string) (*models.User, error)
	FindUserById(id uint64) (*models.User, error)
	CreateUser(userInput dto.SignupInput) error
	DeleteUser(id uint64) error
}

type AuthUsecase struct {
	repository repositories.IAuthRepository
}

func NewAuthUsecase(r repositories.IAuthRepository) IAuthUsecase {
	return &AuthUsecase{repository: r}
}

func (u *AuthUsecase) FindAll() (*[]models.User, error) {
	return u.repository.FindAll()
}

func (u *AuthUsecase) FindUser(email string) (*models.User, error) {
	return u.repository.FindUser(email)
}

func (u *AuthUsecase) FindUserById(id uint64) (*models.User, error) {
	return u.repository.FindUserById(id)
}

func (u *AuthUsecase) CreateUser(userInput dto.SignupInput) error {
	newUser := models.User{
		Name:     userInput.Name,
		Email:    userInput.Email,
		Password: userInput.Password,
	}

	return u.repository.Create(newUser)
}

func (u *AuthUsecase) DeleteUser(id uint64) error {
	return u.repository.DeleteUser(id)
}
