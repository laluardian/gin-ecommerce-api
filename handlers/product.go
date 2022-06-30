package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/laluardian/gin-ecommerce-api/libs"
	"github.com/laluardian/gin-ecommerce-api/models"
	"github.com/laluardian/gin-ecommerce-api/repositories"
	"github.com/rs/xid"
	"gorm.io/gorm"
)

type ProductHandler interface {
	AddProduct(c *gin.Context)
	GetMultipleProducts(c *gin.Context)
	GetProduct(c *gin.Context)
	UpdateProduct(c *gin.Context)
	DeleteProduct(c *gin.Context)
	AddOrRemoveWishlistProduct(c *gin.Context)
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

	payload := libs.CheckUserRole(c)
	if payload == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	// TODO check if the user exist

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
	payload := libs.CheckUserRole(c)
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
	payload := libs.CheckUserRole(c)
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

func (ph *productHandler) AddOrRemoveWishlistProduct(c *gin.Context) {
	var product models.Product
	productId, _ := xid.FromString(c.Param("productId"))
	product, err := ph.repo.FindById(productId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	var user models.User
	payload := c.MustGet(libs.JwtPayloadKey).(*libs.JwtPayload)
	user.ID = payload.Sub

	// if the user already wishlisted the product, remove the product from wishlist
	for _, dbUser := range product.WishlistedBy {
		if dbUser.ID == user.ID {
			if err := ph.repo.RemoveFromWishlist(&product, &user); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "Product successfully removed from wishlist",
			})
			return
		}
	}

	// otherwise, add the product to wishlist
	product.WishlistedBy = append(product.WishlistedBy, &user)
	if err := ph.repo.AddToWishlist(&product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product successfully added to wishlist",
	})
}
