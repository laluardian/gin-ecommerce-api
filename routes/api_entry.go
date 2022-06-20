package routes

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/laluardian/gin-ecommerce-api/handlers"
	"github.com/laluardian/gin-ecommerce-api/middlewares"
	"github.com/laluardian/gin-ecommerce-api/utils"
)

func RunApi() error {
	dsn := os.Getenv("DATA_SOURCE_NAME")
	db := utils.InitDB(dsn)
	userHandler := handlers.NewUserHandler(db)

	r := gin.Default()
	api := r.Group("/api")

	userRoutes := api.Group("/users")
	{
		userRoutes.POST("/signup", userHandler.SignUp)
		userRoutes.POST("/signin", userHandler.SignIn)
		userRoutes.GET("/", userHandler.GetMultipleUsers)
		userRoutes.GET("/:userId", userHandler.GetUser)
	}

	userProtectedRoutes := api.Group("/users", middlewares.JwtAuthorization())
	{
		userProtectedRoutes.PATCH("/:userId", userHandler.UpdateUser)
		userProtectedRoutes.PATCH("/:userId/password", userHandler.UpdatePassword)
		userProtectedRoutes.DELETE("/:userId", userHandler.DeleteUser)
	}

	port := os.Getenv("PORT")
	return r.Run(":" + port)
}
