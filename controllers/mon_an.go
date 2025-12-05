package controllers

import (
	"net/http"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/config"
	"github.com/vpa/quanlynhahang-backend/models"
)

// ======================= CREATE =======================
func CreateMonAn(c *gin.Context) {
	var monan models.MonAn

	// Bind dữ liệu form
	if err := c.ShouldBind(&monan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ: " + err.Error()})
		return
	}

	// Validate
	if monan.TenMonAn == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tên món ăn không được để trống"})
		return
	}

	// Tạo trước để lấy ID
	if err := config.DB.Create(&monan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo món ăn: " + err.Error()})
		return
	}

	// Upload ảnh món ăn nếu có
	file, err := c.FormFile("image")
	if err == nil && file != nil {
		src, err := file.Open()
		if err == nil {
			defer src.Close()

			uploadResult, err := config.CLD.Upload.Upload(c, src, uploader.UploadParams{
				Folder: "monan",
			})

			if err == nil {
				img := models.Images{
					OwnerID:   monan.MaMonAn,
					OwnerType: "mon_an",
					ImageURL:  uploadResult.SecureURL,
				}
				config.DB.Create(&img)
			}
		}
	}

	// Lấy món ăn kèm ảnh trả về client
	config.DB.Preload("AnhMonAn").First(&monan, monan.MaMonAn)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Tạo món ăn thành công",
		"data":    monan,
	})
}

// ======================= GET ALL =======================
func GetAllMonAn(c *gin.Context) {
	var list []models.MonAn
	config.DB.Preload("AnhMonAn").Find(&list)

	c.JSON(http.StatusOK, gin.H{"data": list})
}

func GetMonAnByID(c *gin.Context) {
	id := c.Param("id")
	var monan models.MonAn

	if err := config.DB.Preload("AnhMonAn").First(&monan, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy món ăn"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": monan})
}

// ======================= UPDATE =======================
func UpdateMonAn(c *gin.Context) {
	id := c.Param("id")
	var monan models.MonAn

	if err := config.DB.First(&monan, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy món ăn"})
		return
	}

	// Cập nhật thông tin text
	c.ShouldBind(&monan)
	config.DB.Save(&monan)

	// Nếu có upload ảnh mới → tạo bản ghi mới vào bảng Images (không ghi đè)
	file, err := c.FormFile("image")
	if err == nil {
		src, _ := file.Open()
		defer src.Close()

		upload, err := config.CLD.Upload.Upload(c, src, uploader.UploadParams{Folder: "monan"})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Upload ảnh lỗi"})
			return
		}

		config.DB.Create(&models.Images{
			ImageURL:  upload.SecureURL,
			OwnerID:   monan.MaMonAn,
			OwnerType: "mon_an",
		})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cập nhật món ăn thành công", "data": monan})
}

// ======================= DELETE =======================
func DeleteMonAn(c *gin.Context) {
	id := c.Param("id")
	var monan models.MonAn

	if err := config.DB.First(&monan, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy món ăn"})
		return
	}

	// Xóa ảnh thuộc món ăn
	config.DB.Where("owner_id = ? AND owner_type = ?", id, "mon_an").Delete(&models.Images{})

	// Xóa món ăn
	config.DB.Delete(&monan)

	c.JSON(http.StatusOK, gin.H{"message": "Xóa món ăn thành công"})
}
