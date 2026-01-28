package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/vpa/quanlynhahang-backend/config"
	"github.com/vpa/quanlynhahang-backend/models"
	"github.com/vpa/quanlynhahang-backend/routes"
	"github.com/vpa/quanlynhahang-backend/routes/realtime"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("⚠Không tìm thấy file .env, dùng SECRET_KEY mặc định")
	}
	config.LoadPaymentConfig()
	// 💾 Kết nối Cloudinary
	config.InitCloudinary()
	// 🔧 Khởi tạo Gin
	r := gin.Default()

	// ⚙️ Cấu hình CORS
	config.SetupCORS(r)

	// 💾 Kết nối DB
	config.ConnectDB()

	// 🧱 Tự động migrate
	err := config.DB.AutoMigrate(
		&models.KhachHang{},
		&models.BanAn{},
		&models.MonAn{},
		&models.LoaiMonAn{},
		&models.DatBan{},
		&models.NhanVien{},
		&models.Images{},
		&models.LienHe{},
		&models.HoaDon{},
		&models.Notification{},
		&models.ChiTietHoaDon{},
		&models.ThanhToan{},
	)
	if err != nil {
		log.Fatalf("❌ Lỗi khi migrate DB: %v", err)
	}

	// 🚏 Đăng ký route
	routes.SetupRoutes(r)

	routes.UploadRoutes(r)

	realtime.WebSocketRoutes(r)

	// 🚀 Chạy server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // chạy local
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("❌ Không thể khởi chạy server: %v", err)
	}
}
