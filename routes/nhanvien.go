package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/handlers"
)

func NhanVienRoutes(r *gin.Engine) {
	nhanvien := r.Group("/nhanvien")
	{
		nhanvien.POST("/", handlers.CreateNhanVien)
		nhanvien.GET("/", handlers.GetAllNhanVien)
		nhanvien.GET("/:id", handlers.GetNhanVienByID)
		nhanvien.PUT("/:id", handlers.UpdateNhanVien)
		nhanvien.DELETE("/:id", handlers.DeleteNhanVien)
		nhanvien.POST("/change-password", handlers.ChangePasswordNhanVien)
	}
}
