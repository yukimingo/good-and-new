package repositories

import (
	"errors"
	"good-and-new/models"

	"gorm.io/gorm"
)

type INewsRepository interface {
	FindAll() (*[]models.News, error)
	FindNewsById(id uint64) (*models.News, error)
	Create(news models.News) (*models.News, error)
	Delete(id uint64) error
}

type NewsRepository struct {
	db *gorm.DB
}

func NewNewsRepository(db *gorm.DB) INewsRepository {
	return &NewsRepository{db: db}
}

func (r *NewsRepository) FindAll() (*[]models.News, error) {
	var news []models.News
	if err := r.db.Find(&news).Error; err != nil {
		return nil, err
	}

	return &news, nil
}

func (r *NewsRepository) FindNewsById(id uint64) (*models.News, error) {
	var news models.News
	result := r.db.First(&news, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("news not found")
		}
		return nil, result.Error
	}

	return &news, nil
}

func (r *NewsRepository) Create(news models.News) (*models.News, error) {
	if err := r.db.Create(&news).Error; err != nil {
		return nil, err
	}

	return &news, nil
}

func (r *NewsRepository) Delete(id uint64) error {
	news, err := r.FindNewsById(id)
	if err != nil {
		return err
	}

	if err = r.db.Delete(news).Error; err != nil {
		return err
	}

	return nil
}
