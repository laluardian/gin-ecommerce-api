package repositories

import (
	"github.com/laluardian/gin-ecommerce-api/models"
	"github.com/rs/xid"
	"gorm.io/gorm"
)

type AddressRepository interface {
	Create(*models.Address) error
	FindByUser(userId xid.ID) ([]models.Address, error)
	FindByIds(userId, addressId xid.ID) (models.Address, error)
	Update(models.Address) error
	Delete(xid.ID) error
}

type addressRepository struct {
	db *gorm.DB
}

func NewAddressRepository(db *gorm.DB) AddressRepository {
	return &addressRepository{db}
}

func (ar *addressRepository) Create(address *models.Address) error {
	return ar.db.Create(&address).Error
}

func (ar *addressRepository) FindByUser(userId xid.ID) (addresses []models.Address, err error) {
	err = ar.db.Find(&addresses, "user_id = ?", userId).Error
	return addresses, err
}

func (ar *addressRepository) FindByIds(userId, addressId xid.ID) (address models.Address, err error) {
	err = ar.db.First(&address, "id = ? AND user_id = ?", addressId, userId).Error
	return address, err
}

func (ar *addressRepository) Update(address models.Address) error {
	return ar.db.Save(&address).Error
}

func (ar *addressRepository) Delete(addressId xid.ID) error {
	var address models.Address
	return ar.db.Delete(&address, "id = ?", addressId).Error
}
