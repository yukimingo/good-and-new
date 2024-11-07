package usecases

import (
	"errors"
	"fmt"
	"good-and-new/dto"
	"good-and-new/models"
	"good-and-new/repositories"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type IAuthUsecase interface {
	Login(email string, password string) (*string, error)
	FindAll() (*[]models.User, error)
	FindUser(email string) (*models.User, error)
	FindUserById(id uint64) (*models.User, error)
	CreateUser(userInput dto.SignupInput) (*string, error)
	DeleteUser(id uint64) error
	GetUserFromToken(tokenStrings string) (*models.User, error)
}

type AuthUsecase struct {
	repository repositories.IAuthRepository
}

func NewAuthUsecase(r repositories.IAuthRepository) IAuthUsecase {
	return &AuthUsecase{repository: r}
}

func (u *AuthUsecase) Login(email string, password string) (*string, error) {
	user, err := u.repository.FindUser(email)
	if err != nil {
		return nil, err
	}

	if user.Password != password {
		return nil, errors.New("Invalid password")
	}

	token, err := GenerateToken(uint64(user.ID), user.Email)
	if err != nil {
		return nil, err
	}

	return token, nil
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

func (u *AuthUsecase) CreateUser(userInput dto.SignupInput) (*string, error) {
	newUser := models.User{
		Name:     userInput.Name,
		Email:    userInput.Email,
		Password: userInput.Password,
	}

	user, err := u.repository.FindLatestUser()
	if err != nil {
		return nil, err
	}

	token, err := GenerateToken(uint64(user.ID)+1, newUser.Email)
	if err != nil {
		return nil, err
	}

	return token, u.repository.Create(newUser)
}

func (u *AuthUsecase) DeleteUser(id uint64) error {
	return u.repository.DeleteUser(id)
}

func GenerateToken(userId uint64, email string) (*string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   userId,
		"email": email,
		"exp":   time.Now().Add(time.Hour).Unix(),
	})

	tokenStrings, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return nil, err
	}
	return &tokenStrings, nil
}

func (u *AuthUsecase) GetUserFromToken(tokenStrings string) (*models.User, error) {
	token, err := jwt.Parse(tokenStrings, func(tobj *jwt.Token) (interface{}, error) {
		if _, ok := tobj.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", tobj.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return nil, err
	}

	var user *models.User
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return nil, jwt.ErrTokenExpired
		}

		user, err = u.repository.FindUser(claims["email"].(string))
		if err != nil {
			return nil, err
		}
	}
	return user, nil
}
