package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/controllers"
	"github.com/vpa/quanlynhahang-backend/middleware"
)

func NhanVienRoutes(r *gin.Engine) {
	nhanvien := r.Group("/nhanvien")
	{
		nhanvien.POST("/create", controllers.CreateNhanVien)
		nhanvien.GET("/layTatCa", controllers.GetAllNhanVien)
		nhanvien.GET("/:id", controllers.GetNhanVienByID)
		nhanvien.PUT("/update/:id", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"), controllers.UpdateNhanVien)
		nhanvien.DELETE("/delete/:id", controllers.DeleteNhanVien)
		nhanvien.PUT("/update-profile",
			middleware.AuthMiddleware(),
			controllers.UpdateOwnProfile)

	}
}
