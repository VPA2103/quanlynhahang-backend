package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RoleMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleValue, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Không tìm thấy quyền truy cập"})
			c.Abort()
			return
		}

		userRole := roleValue.(string)

		for _, allowed := range roles {
			if userRole == allowed {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Không có quyền truy cập"})
		c.Abort()
	}
}
