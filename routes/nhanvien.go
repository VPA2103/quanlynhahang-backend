package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/handlers"
)

func NhanVienRoutes(r *gin.Engine) {
	nhanvien := r.Group("/nhanvien")
	{
		nhanvien.POST("/create", handlers.CreateNhanVien)
		nhanvien.GET("/layTatCaNhanVien", handlers.GetAllNhanVien)
		nhanvien.GET("/:id", handlers.GetNhanVienByID)
		nhanvien.PUT("/update/:id", handlers.UpdateNhanVien)
		nhanvien.DELETE("/delete/:id", handlers.DeleteNhanVien)

	}
}
