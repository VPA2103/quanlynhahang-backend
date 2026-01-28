package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/controllers"
)

func LienHeRoutes(r *gin.Engine) {
	lienhe := r.Group("/lien-he")
	{
		lienhe.POST("/create", controllers.GuiLienHe)
		lienhe.GET("", controllers.AdminGetAllLienHe)
		lienhe.DELETE("/:id", controllers.DeleteLienHe)
	}
}
