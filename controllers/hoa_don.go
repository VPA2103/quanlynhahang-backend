package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/services"
)

func ThanhToanHoaDon(c *gin.Context) {

	var req struct {
		MaBan uint `json:"ma_ban"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	if err := services.CloseHoaDon(req.MaBan); err != nil {
		c.JSON(500, gin.H{"error": "Không thể thanh toán"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Thanh toán thành công",
	})
}
