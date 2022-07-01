package routes

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/laluardian/gin-ecommerce-api/handlers"
	"github.com/laluardian/gin-ecommerce-api/libs"
	"github.com/laluardian/gin-ecommerce-api/middlewares"
)

func RunApi() error {
	dsn := os.Getenv("DATA_SOURCE_NAME")
	db := libs.InitDB(dsn)
	userHandler := handlers.NewUserHandler(db)
	productHandler := handlers.NewProductHandler(db)
	addressHandler := handlers.NewAddressHandler(db)
	categoryHandler := handlers.NewCategoryHandler(db)

	r := gin.Default()
	api := r.Group("/api")

	userRoutes := api.Group("/users")
	{
		userRoutes.POST("/signup", userHandler.SignUp)
		userRoutes.POST("/signin", userHandler.SignIn)
	}

	userProtectedRoutes := api.Group("/users", middlewares.JwtAuthorization())
	{
		userProtectedRoutes.GET("/", userHandler.GetMultipleUsers)
		userProtectedRoutes.GET("/:userId", userHandler.GetUser)
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

	categoryRoutes := api.Group("/categories")
	{
		categoryRoutes.GET("/", categoryHandler.GetMultipleCategories)
		categoryRoutes.GET("/:slug", categoryHandler.GetCategory)
	}

	categoryProtectedRoutes := api.Group("/categories", middlewares.JwtAuthorization())
	{
		categoryProtectedRoutes.POST("/", categoryHandler.AddCategory)
		categoryProtectedRoutes.PATCH("/:slug", categoryHandler.UpdateCategory)
		categoryProtectedRoutes.DELETE("/:slug", categoryHandler.DeleteCategory)
	}

	port := os.Getenv("PORT")
	return r.Run(":" + port)
}
