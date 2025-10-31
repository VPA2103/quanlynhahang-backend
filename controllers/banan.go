package controllers

import (
	"net/http"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/config"
	"github.com/vpa/quanlynhahang-backend/models"
)

func CreateBanAn(c *gin.Context) {
	var ban models.BanAn

	// ✅ Bind form data vào struct BanAn
	if err := c.ShouldBind(&ban); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu form không hợp lệ: " + err.Error()})
		return
	}

	// ✅ Nếu trạng thái trống thì mặc định là "Trống"
	if ban.TrangThai == "" {
		ban.TrangThai = "Trống"
	}

	// ✅ Tạo bàn ăn trong DB trước để có MaBan
	if err := config.DB.Create(&ban).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo bàn ăn: " + err.Error()})
		return
	}

	// ✅ Upload ảnh bàn (nếu có)
	file, err := c.FormFile("image")
	if err == nil && file != nil {
		src, err := file.Open()
		if err == nil {
			defer src.Close()

			uploadResult, err := config.CLD.Upload.Upload(c, src, uploader.UploadParams{
				Folder: "banan",
			})
			if err == nil {
				img := models.Images{
					OwnerID:   ban.MaBan,
					OwnerType: "ban_an",
					ImageURL:  uploadResult.SecureURL,
				}
				config.DB.Create(&img)
			}
		}
	}

	// ✅ Preload ảnh bàn
	config.DB.Preload("AnhBan").First(&ban, ban.MaBan)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Tạo bàn ăn thành công",
		"data":    ban,
	})
}

// Lấy tất cả bàn ăn kèm ảnh
func GetAllBanAn(c *gin.Context) {
	var dsBanAn []models.BanAn

	// ✅ Preload ảnh bàn (quan hệ polymorphic)
	if err := config.DB.Preload("AnhBan").Find(&dsBanAn).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể lấy danh sách bàn ăn: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Lấy danh sách bàn ăn thành công",
		"data":    dsBanAn,
	})
}
