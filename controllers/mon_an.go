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

	// Lấy dữ liệu form-data
	if err := c.ShouldBind(&monan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Dữ liệu không hợp lệ: " + err.Error(),
		})
		return
	}

	if monan.TenMonAn == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Tên món ăn không được để trống",
		})
		return
	}

	// Upload ảnh
	file, err := c.FormFile("image")
	if err == nil && file != nil {
		src, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Không thể mở file ảnh",
			})
			return
		}
		defer src.Close()

		uploadResult, err := config.CLD.Upload.Upload(c, src, uploader.UploadParams{
			Folder: "monan",
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Upload ảnh thất bại: " + err.Error(),
			})
			return
		}

		monan.AnhMonAn = uploadResult.SecureURL
	}

	// Lưu database
	if err := config.DB.Create(&monan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Không thể tạo món ăn: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Tạo món ăn thành công",
		"data":    monan,
	})
}

// ======================= GET ALL =======================
func GetAllMonAn(c *gin.Context) {
	var list []models.MonAn

	if err := config.DB.Find(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Không thể lấy danh sách món ăn: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": list,
	})
}

// ======================= GET ONE =======================
func GetMonAnByID(c *gin.Context) {
	id := c.Param("id")
	var monan models.MonAn

	if err := config.DB.First(&monan, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Không tìm thấy món ăn",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": monan,
	})
}

// ======================= UPDATE =======================
func UpdateMonAn(c *gin.Context) {
	id := c.Param("id")
	var monan models.MonAn

	// Kiểm tra tồn tại
	if err := config.DB.First(&monan, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Không tìm thấy món ăn",
		})
		return
	}

	// Bind dữ liệu mới
	if err := c.ShouldBind(&monan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Dữ liệu không hợp lệ: " + err.Error(),
		})
		return
	}

	// Upload ảnh mới nếu có
	file, err := c.FormFile("image")
	if err == nil && file != nil {
		src, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Không thể mở file ảnh",
			})
			return
		}
		defer src.Close()

		uploadResult, err := config.CLD.Upload.Upload(c, src, uploader.UploadParams{
			Folder: "monan",
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Upload ảnh thất bại: " + err.Error(),
			})
			return
		}

		monan.AnhMonAn = uploadResult.SecureURL
	}

	// Cập nhật DB
	if err := config.DB.Save(&monan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Không thể cập nhật món ăn: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cập nhật món ăn thành công",
		"data":    monan,
	})
}

// ======================= DELETE =======================
func DeleteMonAn(c *gin.Context) {
	id := c.Param("id")
	var monan models.MonAn

	// Kiểm tra tồn tại
	if err := config.DB.First(&monan, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Không tìm thấy món ăn",
		})
		return
	}

	// Xóa DB
	if err := config.DB.Delete(&monan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Không thể xóa món ăn: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Xóa món ăn thành công",
	})
}
