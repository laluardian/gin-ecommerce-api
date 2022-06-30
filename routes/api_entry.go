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
	productHandler := handlers.NewProductHandler(db)
	addressHandler := handlers.NewAddressHandler(db)

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
		userProtectedRoutes.GET("/:userId/wishlist", userHandler.GetUserWishlist)
		userProtectedRoutes.PATCH("/:userId", userHandler.UpdateUser)
		userProtectedRoutes.PATCH("/:userId/password", userHandler.UpdatePassword)
		userProtectedRoutes.DELETE("/:userId", userHandler.DeleteUser)
	}

	addressRoutes := userProtectedRoutes.Group("/:userId/addresses")
	{
		addressRoutes.POST("/", addressHandler.AddAddress)
		addressRoutes.GET("/", addressHandler.GetUserAddresses)
		addressRoutes.GET("/:addressId", addressHandler.GetAddress)
		addressRoutes.PATCH("/:addressId", addressHandler.UpdateAddress)
		addressRoutes.DELETE("/:addressId", addressHandler.DeleteAddress)
	}

	productRoutes := api.Group("/products")
	{
		productRoutes.GET("/", productHandler.GetMultipleProducts)
		productRoutes.GET("/:productId", productHandler.GetProduct)
	}

	productProtectedRoutes := api.Group("/products", middlewares.JwtAuthorization())
	{
		productProtectedRoutes.POST("/", productHandler.AddProduct)
		productProtectedRoutes.POST("/:productId/wishlist", productHandler.AddOrRemoveWishlistProduct)
		productProtectedRoutes.PATCH("/:productId", productHandler.UpdateProduct)
		productProtectedRoutes.DELETE("/:productId", productHandler.DeleteProduct)
	}

	port := os.Getenv("PORT")
	return r.Run(":" + port)
}
