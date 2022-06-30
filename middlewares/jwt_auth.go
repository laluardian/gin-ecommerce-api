package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/laluardian/gin-ecommerce-api/libs"
)

func JwtAuthorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if len(authHeader) == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header not found",
			})
		}

		// the Authorization header value looks more or less like this: "Bearer TheToken"
		// in this case we want to get only the "TheToken" part
		const bearerSchema = "Bearer "
		getToken := authHeader[len(bearerSchema):]

		payload, err := libs.VerifyToken(getToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.Set(libs.JwtPayloadKey, payload)
		c.Next()
	}
}
