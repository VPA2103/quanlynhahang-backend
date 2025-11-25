package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/controllers"
	"github.com/vpa/quanlynhahang-backend/middleware"
)

func SetupRoutes(r *gin.Engine) {
	// 🌐 Route gốc
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello Gin!"})
	})

	// 🔐 Auth routes (không cần token)
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	// 🔒 Nhóm yêu cầu xác thực
	auth := r.Group("/api")
	auth.Use(middleware.AuthMiddleware())

	auth.GET("/profile", controllers.GetProfile)

	// 👑 Nhóm chỉ cho admin
	admin := auth.Group("/admin")
	admin.Use(middleware.RoleMiddleware("admin"))
	admin.GET("/dashboard", controllers.AdminDashboard)

	// 👨‍💼 Nhân viên routes (có thể để ngoài hoặc trong nhóm admin)
	NhanVienRoutes(r)
	BanAnRoutes(r)
	LoaiMonAnRoutes(r)
	Payment(r)
}
