package repositories

import (
	"github.com/laluardian/gin-ecommerce-api/models"
	"github.com/rs/xid"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByEmail(email string) (models.User, error)
	FindById(userId xid.ID) (models.User, error)
	FindMany() ([]models.User, error)
	UpdateUser(user *models.User) error
	UpdatePassword(user *models.User) error
	Delete(user *models.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (ur *userRepository) Create(user *models.User) error {
	return ur.db.Create(&user).Error
}

func (ur *userRepository) FindByEmail(email string) (user models.User, err error) {
	err = ur.db.First(&user, "email = ?", email).Error
	return user, err
}

func (ur *userRepository) FindById(userId xid.ID) (user models.User, err error) {
	err = ur.db.First(&user, "id = ?", userId).Error
	return user, err
}

func (ur *userRepository) FindMany() (users []models.User, err error) {
	err = ur.db.Find(&users).Error
	return users, err
}

func (ur *userRepository) UpdateUser(user *models.User) error {
	return ur.db.Save(&user).Error
}

func (ur *userRepository) UpdatePassword(user *models.User) error {
	return ur.db.Model(&user).Update("password", user.Password).Error
}

func (ur *userRepository) Delete(user *models.User) error {
	return ur.db.Delete(&user).Error
}
