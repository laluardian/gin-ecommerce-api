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

type UserHandler interface {
	SignUp(c *gin.Context)
	SignIn(c *gin.Context)
	GetUser(c *gin.Context)
	GetMultipleUsers(c *gin.Context)
	GetUserWishlist(c *gin.Context)
	UpdateUser(c *gin.Context)
	UpdatePassword(c *gin.Context)
	DeleteUser(c *gin.Context)
}

type userHandler struct {
	repo repositories.UserRepository
}

func NewUserHandler(db *gorm.DB) UserHandler {
	return &userHandler{
		repositories.NewUserRepository(db),
	}
}

func (uh *userHandler) SignUp(c *gin.Context) {
	var userInput models.User
	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := libs.HashPassword(&userInput.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := uh.repo.Create(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	token, _ := libs.GenerateToken(&userInput)

	c.JSON(http.StatusCreated, gin.H{
		"access_token": token,
	})
}

func (uh *userHandler) SignIn(c *gin.Context) {
	var userInput models.User
	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	const signInErrMsg = "Invalid email or password"

	user, err := uh.repo.FindByEmail(userInput.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": signInErrMsg,
		})
		return
	}

	if isTrue := libs.ComparePassword(user.Password, userInput.Password); isTrue {
		token, _ := libs.GenerateToken(&user)
		c.JSON(http.StatusOK, gin.H{
			"access_token": token,
		})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{
		"error": signInErrMsg,
	})
}

func (uh *userHandler) GetUser(c *gin.Context) {
	userId, _ := xid.FromString(c.Param("userId"))
	payload := libs.CheckUserId(c, userId)
	if payload == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	user, err := uh.repo.FindById(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func (uh *userHandler) GetMultipleUsers(c *gin.Context) {
	payload := libs.CheckUserRole(c)
	if payload == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	users, err := uh.repo.FindMany()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

func (uh *userHandler) GetUserWishlist(c *gin.Context) {
	userId, _ := xid.FromString(c.Param("userId"))
	payload := libs.CheckUserId(c, userId)
	if payload == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	var user models.User
	user.ID = userId
	products, err := uh.repo.FindUserWishlist(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"wishlist": products,
	})
}

func (uh *userHandler) UpdateUser(c *gin.Context) {
	// get user id from param
	//
	// in the context when the user is already authenticated the user id can also be retrieved
	// directly from the jwt payload (since it is set to be included in the jwt payload)
	//
	// for this particular case (user handler---other handlers might, too) I prefer to retrieve the user id from both
	// param and jwt payload and then compare them manually before performing any db operation in a protected route's handler
	userId, _ := xid.FromString(c.Param("userId"))
	// check if the user id from param matches the user id in jwt payload
	payload := libs.CheckUserId(c, userId)
	if payload == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	// check if the user with that id exists
	dbUser, err := uh.repo.FindById(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var userInput models.User
	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// in this scenario, these two fields are not going to be updated
	// id cannot and must not be updated
	// password can be updated, but from a different endpoint
	//
	// there are other ways to exclude these fields from being updated
	// imo, this is the simplest way for this particular case
	userInput.ID = dbUser.ID
	userInput.Password = dbUser.Password

	if err := uh.repo.UpdateUser(&userInput); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User successfully updated",
	})
}

func (uh *userHandler) UpdatePassword(c *gin.Context) {
	userId, _ := xid.FromString(c.Param("userId"))
	payload := libs.CheckUserId(c, userId)
	if payload == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	dbUser, err := uh.repo.FindById(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var userInput models.User
	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// check if the old password is the same as the new password
	if isTrue := libs.ComparePassword(dbUser.Password, userInput.Password); isTrue {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "The old password cannot be the same as the new password",
		})
		return
	}

	if err := libs.HashPassword(&userInput.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userInput.ID = dbUser.ID
	if err := uh.repo.UpdatePassword(&userInput); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password successfully updated",
	})
}

func (uh *userHandler) DeleteUser(c *gin.Context) {
	userId, _ := xid.FromString(c.Param("userId"))
	payload := libs.CheckUserId(c, userId)
	if payload == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	var user models.User
	user.ID = userId
	if err := uh.repo.Delete(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User successfully deleted",
	})
}
