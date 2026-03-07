package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/services"
)

type MonOrder struct {
	MaMonAn uint `json:"ma_mon_an"`
	SoLuong int  `json:"so_luong"`
}

type AddMonRequest struct {
	MaBan  uint       `json:"ma_ban"`
	MonAns []MonOrder `json:"mon_ans"`
}

func AddMon(c *gin.Context) {

	var req AddMonRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	for _, mon := range req.MonAns {

		err := services.AddMon(req.MaBan, mon.MaMonAn, mon.SoLuong)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Gọi món thành công",
	})
}
