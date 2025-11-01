package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/controllers"
)

func NhanVienRoutes(r *gin.Engine) {
	nhanvien := r.Group("/nhanvien")
	{
		nhanvien.POST("/create", controllers.CreateNhanVien)
		nhanvien.GET("/layTatCa", controllers.GetAllNhanVien)
		nhanvien.GET("/:id", controllers.GetNhanVienByID)
		nhanvien.PUT("/update/:id", controllers.UpdateNhanVien)
		nhanvien.DELETE("/delete/:id", controllers.DeleteNhanVien)

	}
}
