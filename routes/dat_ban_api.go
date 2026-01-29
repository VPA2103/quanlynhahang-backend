package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/controllers"
	"github.com/vpa/quanlynhahang-backend/middleware"
)

func DatBanRoutes(r *gin.Engine) {
	datban := r.Group("/dat-ban")
	{
		// Khách
		datban.POST("", controllers.CreateDatBan)                                                                 // tạo đặt bàn //ok
		datban.GET("", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"), controllers.GetAllDatBan) // danh sách //ok
		datban.GET("/:id", controllers.GetDatBanByID)                                                             // chi tiết
		datban.PUT("/:id", controllers.UpdateDatBan)                                                              // sửa thông tin
		datban.DELETE("/:id", controllers.DeleteDatBan)                                                           //ok

		// Nhân viên
		datban.PUT("/:id/xac-nhan", middleware.AuthMiddleware(), controllers.XacNhanDatBan)
	}
}
