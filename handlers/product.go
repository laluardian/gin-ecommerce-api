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

type ProductHandler interface {
	AddProduct(c *gin.Context)
	GetMultipleProducts(c *gin.Context)
	GetProduct(c *gin.Context)
	UpdateProduct(c *gin.Context)
	DeleteProduct(c *gin.Context)
}

type productHandler struct {
	repo repositories.ProductRepository
}

func NewProductHandler(db *gorm.DB) ProductHandler {
	return &productHandler{
		repositories.NewProductRepository(db),
	}
}

func (ph *productHandler) AddProduct(c *gin.Context) {
	var productInput models.Product
	if err := c.ShouldBindJSON(&productInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	payload := utils.CheckUserRole(c)
	if payload == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	if err := ph.repo.Create(&productInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "A new product successfully added",
	})
}

func (ph *productHandler) GetMultipleProducts(c *gin.Context) {
	keyword := c.Query("search")
	// if the keyword is empty all products will be returned
	products, err := ph.repo.FindMany(keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
	})
}

func (ph *productHandler) GetProduct(c *gin.Context) {
	productId, _ := xid.FromString(c.Param("productId"))
	product, err := ph.repo.FindById(productId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"product": product,
	})
}

func (ph *productHandler) UpdateProduct(c *gin.Context) {
	payload := utils.CheckUserRole(c)
	if payload == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	var productInput models.Product
	if err := c.ShouldBindJSON(&productInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	productId, _ := xid.FromString(c.Param("productId"))
	productInput.ID = productId
	if err := ph.repo.Update(&productInput); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product successfully updated",
	})
}

func (ph *productHandler) DeleteProduct(c *gin.Context) {
	payload := utils.CheckUserRole(c)
	if payload == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	productId, _ := xid.FromString(c.Param("productId"))
	if err := ph.repo.Delete(productId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product successfully deleted",
	})
}
