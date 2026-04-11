package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/vpa/quanlynhahang-backend/config"
	"github.com/vpa/quanlynhahang-backend/models"
	"github.com/vpa/quanlynhahang-backend/routes"
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
		// 1. Bảng nền (không phụ thuộc)
		&models.KhachHang{},
		&models.NhanVien{},
		&models.LoaiMonAn{},
		&models.BanAn{},

		// 2. Bảng phụ thuộc mức 1
		&models.MonAn{},  // phụ thuộc LoaiMonAn
		&models.DatBan{}, // phụ thuộc NhanVien, BanAn
		&models.HoaDon{}, // phụ thuộc KhachHang, NhanVien

		// 3. Bảng phụ thuộc mức 2
		&models.ChiTietHoaDon{}, // phụ thuộc HoaDon, MonAn
		&models.ThanhToan{},     // phụ thuộc HoaDon

		// 4. Polymorphic / phụ
		&models.Images{},
		&models.Notification{},
		&models.LienHe{},
	)

	if err != nil {
		log.Fatalf("❌ Lỗi khi migrate DB: %v", err)
	}

	// 🚏 Đăng ký route
	routes.SetupRoutes(r)

	routes.UploadRoutes(r)

	routes.WebSocketRoutes(r)

	// 🚀 Chạy server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // chạy local
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("❌ Không thể khởi chạy server: %v", err)
	}
}
