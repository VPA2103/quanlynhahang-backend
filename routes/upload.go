package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/handlers"
)

func UploadRoutes(route *gin.Engine) {
	route.POST("/upload", handlers.UploadHandler)
	route.GET("/images", handlers.GetProductImage)
	//route.GET("/images", controllers.GetProductImage)

}
