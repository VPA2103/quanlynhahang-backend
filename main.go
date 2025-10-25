package main

import (
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/config"
	"github.com/vpa/quanlynhahang-backend/handlers"
	"github.com/vpa/quanlynhahang-backend/middleware"
	"github.com/vpa/quanlynhahang-backend/models"
)

func main() {

	r := gin.Default()
	config.SetupCORS(r)

	//
	config.ConnectDB()
	config.DB.AutoMigrate(
		&models.KhachHang{},
		&models.BanAn{},
		&models.MonAn{},
		&models.LoaiMonAn{},
		&models.DatBan{},
		&models.NhanVien{},
		&models.HoaDon{},
		&models.ChiTietHoaDon{},
		&models.ThanhToan{},
	)

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello Gin!",
		})
	})

	auth := r.Group("/api")
	auth.Use(middleware.AuthMiddleware()) // Kiểm tra token JWT

	auth.GET("/profile", handlers.GetProfile) // Bất kỳ user nào có token đều xem được

	admin := auth.Group("/admin")
	admin.Use(middleware.RoleMiddleware("admin")) // Chỉ admin mới truy cập
	admin.GET("/dashboard", handlers.AdminDashboard)

	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	r.Run(":3000") // chạy ở port 3000
}
