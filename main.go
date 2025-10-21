package main

import (
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/config"
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

	r.Run(":8080") // chạy ở port 8080
}
