package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/laluardian/gin-ecommerce-api/libs"
	"github.com/laluardian/gin-ecommerce-api/models"
	"github.com/laluardian/gin-ecommerce-api/repositories"
	"gorm.io/gorm"
)

type CategoryHandler interface {
	AddCategory(c *gin.Context)
	GetCategory(c *gin.Context)
	GetMultipleCategories(c *gin.Context)
	UpdateCategory(c *gin.Context)
	DeleteCategory(c *gin.Context)
}

type categoryHandler struct {
	repo repositories.CategoryRepository
}

func NewCategoryHandler(db *gorm.DB) CategoryHandler {
	return &categoryHandler{
		repositories.NewCategoryRepository(db),
	}
}

func (ch *categoryHandler) AddCategory(c *gin.Context) {
	payload := libs.CheckUserRole(c)
	if payload == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	var categoryInput models.Category
	if err := c.ShouldBindJSON(&categoryInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	if err := ch.repo.Create(&categoryInput); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "A new category successfully added",
	})
}

func (ch *categoryHandler) GetMultipleCategories(c *gin.Context) {
	categories, err := ch.repo.FindMany()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"categories": categories,
	})
}

func (ch *categoryHandler) GetCategory(c *gin.Context) {
	slug := c.Param("slug")
	category, err := ch.repo.FindBySlug(slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"category":                 category,
		"category_products_length": len(category.Products),
	})
}

func (ch *categoryHandler) UpdateCategory(c *gin.Context) {
	payload := libs.CheckUserRole(c)
	fmt.Println(payload)
	if payload == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	var categoryInput models.Category
	if err := c.ShouldBindJSON(&categoryInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	// in this case getting the category record from db is needed in order to get the category id
	// which is needed for updating the category record, this step can be skipped if we instead set
	// the category primary key with slug OR we can also include the id in the URL param and get the
	// id from there... and THERE MUST BE MANY OTHER WAYS to achieve this though and surely this is
	// not the best way, but for this 'example' project I think doing it this way is enough...
	slug := c.Param("slug")
	dbCategory, err := ch.repo.FindBySlug(slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	categoryInput.ID = dbCategory.ID
	if err := ch.repo.Update(&categoryInput); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Category successfully updated",
	})
}

func (ch *categoryHandler) DeleteCategory(c *gin.Context) {
	payload := libs.CheckUserRole(c)
	if payload == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	slug := c.Param("slug")
	if err := ch.repo.Delete(slug); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Category successfully deleted",
	})
}
