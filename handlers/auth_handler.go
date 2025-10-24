package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/config"
	"github.com/vpa/quanlynhahang-backend/models"
	"github.com/vpa/quanlynhahang-backend/utils"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var user models.NhanVien
	if err := config.DB.Where("ten_dang_nhap = ?", input.Username).First(&user).Error; err != nil {

		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// So sánh password (hash)
	if err := bcrypt.CompareHashAndPassword([]byte(user.MatKhau), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect password"})
		return
	}

	// Tạo JWT token
	token, err := utils.GenerateToken(user.MaNV, user.HoTen, user.LoaiNhanVien)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}

func Register(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Mã hoá mật khẩu
	hashed, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 10)

	user := models.NhanVien{
		TenDangNhap:  input.Username,
		MatKhau:      string(hashed),
		LoaiNhanVien: input.Role,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered"})
}

// Handler admin
func AdminDashboard(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Welcome Admin Dashboard"})
}

func GetProfile(c *gin.Context) {
	c.JSON(200, gin.H{"message": "User profile"})
}
