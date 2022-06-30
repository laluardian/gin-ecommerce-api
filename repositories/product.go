package repositories

import (
	"github.com/laluardian/gin-ecommerce-api/models"
	"github.com/rs/xid"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(product *models.Product) error
	FindMany(keyword string) ([]models.Product, error)
	FindById(userId xid.ID) (models.Product, error)
	Update(product *models.Product) error
	Delete(productId xid.ID) error
	AddToWishlist(product *models.Product) error
	RemoveFromWishlist(product *models.Product, user *models.User) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db}
}

func (pr *productRepository) Create(product *models.Product) error {
	return pr.db.Create(&product).Error
}

func (pr *productRepository) FindMany(keyword string) (products []models.Product, err error) {
	err = pr.db.Find(&products, "LOWER(name) LIKE LOWER(?)", "%"+keyword+"%").Error
	return products, err
}

func (pr *productRepository) FindById(productId xid.ID) (product models.Product, err error) {
	err = pr.db.Preload("WishlistedBy").First(&product, "id = ?", productId).Error
	return product, err
}

func (pr *productRepository) Update(product *models.Product) error {
	return pr.db.Save(&product).Error
}

func (pr *productRepository) Delete(productId xid.ID) error {
	var product models.Product
	return pr.db.Delete(&product, "id = ?", productId).Error
}

func (pr *productRepository) AddToWishlist(product *models.Product) error {
	err := pr.db.
		Omit("WishlistedBy.*").
		Session(&gorm.Session{FullSaveAssociations: true}).
		Updates(&product).Error
	return err
}

func (pr *productRepository) RemoveFromWishlist(product *models.Product, user *models.User) error {
	err := pr.db.Model(&product).Association("WishlistedBy").Delete(user)
	return err
}
