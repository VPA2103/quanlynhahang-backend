package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/vpa/quanlynhahang-backend/utils"
	//"github.com/vpa/quanlynhahang-backend/utils"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// ✅ Dùng utils.JWTSecret() để lấy key
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return utils.JWTSecret(), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		userID := uint(claims["id"].(float64))
		username := claims["email"].(string)
		role := claims["role"].(string)

		c.Set("user_id", userID)
		c.Set("username", username)
		c.Set("role", role)

		c.Next()
	}
}
