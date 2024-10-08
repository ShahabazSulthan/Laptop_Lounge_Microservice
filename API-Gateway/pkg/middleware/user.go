package middleware

import (
	"api-gateway/pkg/helper"
	"api-gateway/pkg/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserAuthMiddleware ensures that a request is authenticated as a user
func UserAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			authHeader, _ = c.Cookie("Authorization")
		}

		if authHeader == "" {
			response := response.ClientResponse(http.StatusUnauthorized, "No auth header provided", nil, nil)
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}

		tokenString := helper.GetTokenFromHeader(authHeader)
		userID, userEmail, err := helper.ExtractUserIDFromToken(tokenString)
		if err != nil {
			response := response.ClientResponse(http.StatusUnauthorized, "Invalid Token", nil, err.Error())
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Set("user_email", userEmail)
		c.Next()
	}
}
