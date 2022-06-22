package repositories

import (
	"github.com/laluardian/gin-ecommerce-api/models"
	"github.com/rs/xid"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(*models.Product) error
	FindMany(string) ([]models.Product, error)
	FindById(xid.ID) (models.Product, error)
	Update(models.Product) error
	Delete(xid.ID) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db}
}

func (pr *productRepository) Create(product *models.Product) error {
	return pr.db.Create(product).Error
}

func (pr *productRepository) FindMany(keyword string) (products []models.Product, err error) {
	err = pr.db.Find(&products, "LOWER(name) LIKE LOWER(?)", "%"+keyword+"%").Error
	return products, err
}

func (pr *productRepository) FindById(productId xid.ID) (product models.Product, err error) {
	err = pr.db.First(&product, "id = ?", productId).Error
	return product, err
}

func (pr *productRepository) Update(product models.Product) error {
	return pr.db.Save(&product).Error
}

func (pr *productRepository) Delete(productId xid.ID) error {
	var product models.Product
	return pr.db.Delete(&product, "id = ?", productId).Error
}
