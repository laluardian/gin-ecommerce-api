package repositories

import (
	"github.com/laluardian/gin-ecommerce-api/models"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(category *models.Category) error
	FindMany() ([]models.Category, error)
	FindBySlug(slug string) (models.Category, error)
	Update(category *models.Category) error
	Delete(slug string) error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db}
}

func (cr *categoryRepository) Create(category *models.Category) error {
	err := cr.db.Create(&category).Error
	return err
}

func (cr *categoryRepository) FindMany() (categories []models.Category, err error) {
	err = cr.db.Find(&categories).Error
	return categories, err
}

func (cr *categoryRepository) FindBySlug(slug string) (category models.Category, err error) {
	err = cr.db.Preload("Products").First(&category, "slug = ?", slug).Error
	return category, err
}

func (cr *categoryRepository) Update(category *models.Category) error {
	err := cr.db.Updates(&category).Error
	return err
}

func (cr *categoryRepository) Delete(slug string) error {
	var category models.Category
	err := cr.db.Delete(&category, "slug = ?", slug).Error
	return err
}
