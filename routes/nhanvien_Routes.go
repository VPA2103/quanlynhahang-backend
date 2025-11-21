package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/controllers"
	"github.com/vpa/quanlynhahang-backend/middleware"
)

func NhanVienRoutes(r *gin.Engine) {
	nhanvien := r.Group("/nhanvien")
	{
		// ✅ Chỉ admin được phép
		nhanvien.POST("/create", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"), controllers.CreateNhanVien)
		nhanvien.PATCH("/update/:id", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"), controllers.UpdateNhanVien)
		nhanvien.DELETE("/delete/:id", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"), controllers.DeleteNhanVien)

		nhanvien.GET("/layRaThongTinNhanVien/:id", controllers.GetNhanVienByID)

		// ✅ Chỉ nhân viên được phép
		nhanvien.PATCH("/capNhatThongTinCaNhan/:id", middleware.AuthMiddleware(), middleware.RoleMiddleware("user"), controllers.UpdateThongTinCaNhan)

		// ✅ Cả admin và user đều có thể xem danh sách
		nhanvien.GET("/layTatCa", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin", "user"), controllers.GetAllNhanVien)

	}
}
