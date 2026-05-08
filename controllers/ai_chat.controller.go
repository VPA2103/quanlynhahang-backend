package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/services"
)

type ChatRequest struct {
	Message string `json:"message"`
	UserID  string `json:"user_id"`
}

func AIChat(c *gin.Context) {
	var req ChatRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Dữ liệu không hợp lệ",
		})
		return
	}

	// ✅ trả về AIChatResponse
	res := services.HandleFoodChat(req.UserID, req.Message)

	c.JSON(http.StatusOK, res)
}