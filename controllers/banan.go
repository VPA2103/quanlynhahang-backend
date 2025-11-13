package controllers

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/config"
	"github.com/vpa/quanlynhahang-backend/models"
	"github.com/vpa/quanlynhahang-backend/utils"
)

func CreateBanAn(c *gin.Context) {
	var ban models.BanAn

	// ✅ Bind form data
	if err := c.ShouldBind(&ban); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu form không hợp lệ: " + err.Error()})
		return
	}

	// ✅ Mặc định trạng thái là "Trống"
	if ban.TrangThai == "" {
		ban.TrangThai = "Trống"
	}

	// ✅ Tạo record trong DB trước để có MaBan
	if err := config.DB.Create(&ban).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo bàn ăn: " + err.Error()})
		return
	}

	// ✅ Tạo QR trong bộ nhớ
	qrBytes, err := utils.GenerateQRBytes(int(ban.MaBan), ban.TenBan, ban.SoChoNgoi, ban.TrangThai)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo mã QR: " + err.Error()})
		return
	}

	// ✅ Upload QR trực tiếp lên Cloudinary
	uploadResult, err := config.CLD.Upload.Upload(c, bytes.NewReader(qrBytes), uploader.UploadParams{
		Folder:   "banan_qr",
		PublicID: fmt.Sprintf("qr_ban_%d", ban.MaBan),
	})
	if err == nil {
		ban.Anh_QR = uploadResult.SecureURL
		config.DB.Save(&ban)
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

// ✅ Cập nhật thông tin bàn ăn
func UpdateBanAn(c *gin.Context) {
	id := c.Param("id")
	var ban models.BanAn

	// 🔹 Tìm bàn ăn theo ID
	if err := config.DB.First(&ban, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy bàn ăn"})
		return
	}

	// 🔹 Bind dữ liệu form
	var input models.BanAn
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu gửi lên không hợp lệ: " + err.Error()})
		return
	}

	// 🔹 Cập nhật thông tin
	ban.TenBan = input.TenBan
	ban.SoChoNgoi = input.SoChoNgoi
	ban.TrangThai = input.TrangThai

	if err := config.DB.Save(&ban).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể cập nhật bàn ăn: " + err.Error()})
		return
	}

	// 🔹 Nếu có upload ảnh mới
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

	config.DB.Preload("AnhBan").First(&ban, ban.MaBan)

	c.JSON(http.StatusOK, gin.H{
		"message": "Cập nhật bàn ăn thành công",
		"data":    ban,
	})
}

// ✅ Xóa bàn ăn
func DeleteBanAn(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	var ban models.BanAn
	if err := config.DB.First(&ban, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy bàn ăn"})
		return
	}

	// 🔹 Xóa ảnh liên quan (nếu có)
	config.DB.Where("owner_id = ? AND owner_type = ?", id, "ban_an").Delete(&models.Images{})

	// 🔹 Xóa bàn ăn
	if err := config.DB.Delete(&ban).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể xóa bàn ăn: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Xóa bàn ăn thành công",
	})
}
