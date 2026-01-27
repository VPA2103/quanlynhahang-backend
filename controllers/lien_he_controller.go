package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/config"
	"github.com/vpa/quanlynhahang-backend/models"
)

func GuiLienHe(c *gin.Context) {
	var lienHe models.LienHe

	if err := c.ShouldBind(&lienHe); err != nil {
		fmt.Println("❌ Bind error:", err.Error())

		c.JSON(400, gin.H{
			"message": "Dữ liệu không hợp lệ",
			"error":   err.Error(),
		})
		return
	}

	if lienHe.HoTen == "" || lienHe.Email == "" || lienHe.TieuDe == "" || lienHe.NoiDung == "" {
		c.JSON(400, gin.H{"message": "Vui lòng nhập đầy đủ thông tin"})
		return
	}

	if err := config.DB.Create(&lienHe).Error; err != nil {
		c.JSON(500, gin.H{"message": "Lưu liên hệ thất bại"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Gửi thành công",
		"data":    lienHe,
	})
}

func AdminGetAllLienHe(c *gin.Context) {
	// 👉 Nếu bạn đã có middleware check admin
	// thì KHÔNG cần đoạn check quyền ở đây

	var danhSachLienHe []models.LienHe

	if err := config.DB.
		Order("created_at DESC").
		Find(&danhSachLienHe).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Không thể lấy danh sách liên hệ",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Lấy danh sách liên hệ thành công",
		"data":    danhSachLienHe,
	})
}
