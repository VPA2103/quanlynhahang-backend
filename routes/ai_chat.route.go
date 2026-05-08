package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/controllers"
)

func AIChatRoute(r *gin.Engine) {
	api := r.Group("/api/ai")
	{
		api.POST("/chat", controllers.AIChat)
	}
}