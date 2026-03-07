package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/services"
)

type AddMonRequest struct {
	MaBan   uint `json:"ma_ban"`
	MaMonAn uint `json:"ma_mon_an"`
	SoLuong int  `json:"so_luong"`
}

func AddMon(c *gin.Context) {

	var req AddMonRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := services.AddMon(req.MaBan, req.MaMonAn, req.SoLuong)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Gọi món thành công",
	})

}
