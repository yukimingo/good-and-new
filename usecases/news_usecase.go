package usecases

import (
	"good-and-new/dto"
	"good-and-new/models"
	"good-and-new/repositories"
)

type INewsUsecase interface {
	FindAll() (*[]models.News, error)
	FindById(id uint64) (*models.News, error)
	Create(newsInput dto.NewsInput, userId uint) (*models.News, error)
	Delete(id uint64) error
}

type NewsUsecase struct {
	repository repositories.INewsRepository
}

func NewNewsUsecase(r repositories.INewsRepository) INewsUsecase {
	return &NewsUsecase{repository: r}
}

func (u *NewsUsecase) FindAll() (*[]models.News, error) {
	return u.repository.FindAll()
}

func (u *NewsUsecase) FindById(id uint64) (*models.News, error) {
	return u.repository.FindNewsById(id)
}

func (u *NewsUsecase) Create(newsInput dto.NewsInput, userId uint) (*models.News, error) {
	newNews := models.News{
		Title:       newsInput.Title,
		Description: newsInput.Description,
		UserID:      userId,
	}

	return u.repository.Create(newNews)
}

func (u *NewsUsecase) Delete(id uint64) error {
	return u.repository.Delete(id)
}
