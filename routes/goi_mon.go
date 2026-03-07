package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/controllers"
)

func GoiMonRoutes(r *gin.Engine) {
	goimon := r.Group("/goi-mon")
	{
		goimon.POST("/create", controllers.AddMon)
	}
}
