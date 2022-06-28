package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/laluardian/gin-ecommerce-api/models"
	"github.com/laluardian/gin-ecommerce-api/repositories"
	"github.com/laluardian/gin-ecommerce-api/utils"
	"github.com/rs/xid"
	"gorm.io/gorm"
)

type AddressHandler interface {
	AddAddress(c *gin.Context)
	GetUserAddresses(c *gin.Context)
	GetAddress(c *gin.Context)
	UpdateAddress(c *gin.Context)
	DeleteAddress(c *gin.Context)
}

type addressHandler struct {
	repo repositories.AddressRepository
}

func NewAddressHandler(db *gorm.DB) AddressHandler {
	return &addressHandler{
		repositories.NewAddressRepository(db),
	}
}

func (ah *addressHandler) AddAddress(c *gin.Context) {
	userId, _ := xid.FromString(c.Param("userId"))
	payload := utils.CheckUserId(c, userId)
	if payload == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	var addressInput models.Address
	if err := c.ShouldBindJSON(&addressInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	addressInput.UserID = userId
	if user := ah.repo.Create(&addressInput); user != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"user": user,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "A new address successfully added",
	})
}

func (ah *addressHandler) GetUserAddresses(c *gin.Context) {
	userId, _ := xid.FromString(c.Param("userId"))
	payload := utils.CheckUserId(c, userId)
	if payload == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	addresses, err := ah.repo.FindByUser(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"addresses": addresses,
	})
}

func (ah *addressHandler) GetAddress(c *gin.Context) {
	userId, _ := xid.FromString(c.Param("userId"))
	payload := utils.CheckUserId(c, userId)
	if payload == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	addressId, _ := xid.FromString(c.Param("addressId"))
	address, err := ah.repo.FindByIds(userId, addressId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"address": address,
	})
}

func (ah *addressHandler) UpdateAddress(c *gin.Context) {
	userId, _ := xid.FromString(c.Param("userId"))
	payload := utils.CheckUserId(c, userId)
	if payload == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	var addressInput models.Address
	if err := c.ShouldBindJSON(&addressInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	addressId, _ := xid.FromString(c.Param("addressId"))
	addressInput.ID = addressId
	addressInput.UserID = userId
	if err := ah.repo.Update(&addressInput); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Address successfully updated",
	})
}

func (ah *addressHandler) DeleteAddress(c *gin.Context) {
	userId, _ := xid.FromString(c.Param("userId"))
	payload := utils.CheckUserId(c, userId)
	if payload == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	addressId, _ := xid.FromString(c.Param("addressId"))
	if err := ah.repo.Delete(addressId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Address successfully deleted",
	})
}
