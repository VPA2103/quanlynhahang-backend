package controllers

import (
	"net/http"
	"time"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/config"
	"github.com/vpa/quanlynhahang-backend/models"
	"golang.org/x/crypto/bcrypt"
)

// 🧱 Thêm nhân viên
func CreateNhanVien(c *gin.Context) {
	var nv models.NhanVien

	// ✅ Lấy dữ liệu từ form-data
	if err := c.ShouldBind(&nv); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu form không hợp lệ: " + err.Error()})
		return
	}

	// ✅ Kiểm tra loại nhân viên chỉ được phép là "user" hoặc "admin"
	if nv.LoaiNhanVien != "user" && nv.LoaiNhanVien != "admin" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Loại nhân viên không hợp lệ. Chỉ chấp nhận 'user' hoặc 'admin'."})
		return
	}

	// ✅ Mặc định ngày vào làm
	if nv.NgayVaoLam == "" {
		nv.NgayVaoLam = time.Now().Format("02-01-2006 15:04:05")
	}

	// ✅ Kiểm tra mật khẩu
	if nv.MatKhau == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mật khẩu không được để trống"})
		return
	}

	// ✅ Hash mật khẩu
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(nv.MatKhau), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể mã hóa mật khẩu"})
		return
	}
	nv.MatKhau = string(hashedPassword)

	// ✅ Lưu nhân viên trước để có MaNV (ID)
	if err := config.DB.Create(&nv).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo nhân viên: " + err.Error()})
		return
	}

	// ✅ Upload ảnh (nếu có)
	file, err := c.FormFile("image")
	if err == nil && file != nil {
		src, err := file.Open()
		if err == nil {
			defer src.Close()

			uploadResult, err := config.CLD.Upload.Upload(c, src, uploader.UploadParams{
				Folder: "nhanvien",
			})
			if err == nil {
				img := models.Images{
					OwnerID:   nv.MaNV,
					OwnerType: "nhan_vien",
					ImageURL:  uploadResult.SecureURL,
				}
				config.DB.Create(&img)
			}
		}
	}

	// ✅ Preload ảnh khi trả về
	config.DB.Preload("AnhNhanVien").First(&nv, nv.MaNV)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Tạo nhân viên thành công",
		"data":    nv,
	})
}

// 📋 Lấy danh sách nhân viên
func GetAllNhanVien(c *gin.Context) {
	var nhanViens []models.NhanVien
	if err := config.DB.Preload("AnhNhanVien").Find(&nhanViens).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nhanViens)
}

// 🔍 Lấy 1 nhân viên theo ID
func GetNhanVienByID(c *gin.Context) {
	id := c.Param("id")
	var nv models.NhanVien
	if err := config.DB.Preload("AnhNhanVien").Find(&nv, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nv)
}

// ✏️ Cập nhật nhân viên
func UpdateNhanVien(c *gin.Context) {
	id := c.Param("id")
	var nv models.NhanVien

	// 🔹 Lấy loại nhân viên đang đăng nhập
	roleValue, _ := c.Get("loai_nhan_vien")
	currentRole, _ := roleValue.(string)

	// 🔹 Tìm nhân viên theo ID
	if err := config.DB.Preload("AnhNhanVien").First(&nv, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy nhân viên"})
		return
	}

	// 🔹 Lấy dữ liệu form
	hoTen := c.PostForm("ho_ten")
	gioiTinh := c.PostForm("gioi_tinh")
	ngaySinh := c.PostForm("ngay_sinh")
	sdt := c.PostForm("sdt")
	diaChi := c.PostForm("dia_chi")
	email := c.PostForm("email")

	oldPassword := c.PostForm("mat_khau_cu")
	newPassword := c.PostForm("mat_khau_moi")
	confirmPassword := c.PostForm("xac_nhan_mat_khau_moi")

	// ✅ Admin có thể chỉnh sửa tất cả loại nhân viên
	// nhưng nhân viên thường thì KHÔNG được thay đổi loại của mình

	// ✅ Cập nhật các thông tin cơ bản (ai cũng có thể)
	if hoTen != "" {
		nv.HoTen = hoTen
	}
	if gioiTinh != "" {
		nv.GioiTinh = gioiTinh
	}
	if ngaySinh != "" {
		nv.NgaySinh = ngaySinh
	}
	if sdt != "" {
		nv.SDT = sdt
	}
	if diaChi != "" {
		nv.DiaChi = diaChi
	}
	if email != "" {
		nv.Email = email
	}

	// ✅ Xử lý đổi mật khẩu
	if currentRole == "admin" {
		// 👑 ADMIN có thể đổi trực tiếp (chỉ cần nhập mật khẩu mới)
		if newPassword != "" {
			if newPassword != confirmPassword {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Xác nhận mật khẩu mới không khớp"})
				return
			}
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
			nv.MatKhau = string(hashedPassword)
		}
	} else {
		// 👤 USER phải nhập đúng mật khẩu cũ
		if oldPassword != "" || newPassword != "" || confirmPassword != "" {
			if oldPassword == "" || newPassword == "" || confirmPassword == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Cần nhập đủ mật khẩu cũ, mật khẩu mới và xác nhận mật khẩu mới"})
				return
			}

			if err := bcrypt.CompareHashAndPassword([]byte(nv.MatKhau), []byte(oldPassword)); err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Mật khẩu cũ không đúng"})
				return
			}

			if newPassword != confirmPassword {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Xác nhận mật khẩu mới không khớp"})
				return
			}

			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
			nv.MatKhau = string(hashedPassword)
		}
	}

	// ✅ Upload ảnh mới (nếu có)
	file, err := c.FormFile("image")
	if err == nil && file != nil {
		src, _ := file.Open()
		defer src.Close()

		uploadResult, err := config.CLD.Upload.Upload(c, src, uploader.UploadParams{Folder: "nhanvien"})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Upload ảnh thất bại: " + err.Error()})
			return
		}

		// Xóa ảnh cũ
		config.DB.Where("owner_id = ? AND owner_type = ?", nv.MaNV, "nhan_vien").Delete(&models.Images{})

		// Lưu ảnh mới
		newImg := models.Images{
			OwnerID:   nv.MaNV,
			OwnerType: "nhan_vien",
			ImageURL:  uploadResult.SecureURL,
		}
		config.DB.Create(&newImg)
	}

	// ✅ Lưu thay đổi
	if err := config.DB.Save(&nv).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể cập nhật nhân viên: " + err.Error()})
		return
	}

	// ✅ Lấy lại thông tin mới
	config.DB.Preload("AnhNhanVien").First(&nv, nv.MaNV)

	c.JSON(http.StatusOK, gin.H{
		"message": "Cập nhật nhân viên thành công",
		"data":    nv,
	})
}

// 🗑️ Xóa nhân viên
func DeleteNhanVien(c *gin.Context) {
	id := c.Param("id")
	var nv models.NhanVien
	if err := config.DB.First(&nv, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy nhân viên"})
		return
	}

	if err := config.DB.Delete(&nv).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Đã xóa nhân viên thành công"})
}

func UpdateThongTinCaNhan(c *gin.Context) {
	id := c.Param("id")
	var nv models.NhanVien

	if err := config.DB.Preload("AnhNhanVien").First(&nv, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy nhân viên"})
		return
	}

	// Lấy dữ liệu form
	hoTen := c.PostForm("ho_ten")
	gioiTinh := c.PostForm("gioi_tinh")
	ngaySinh := c.PostForm("ngay_sinh")
	sdt := c.PostForm("sdt")
	diaChi := c.PostForm("dia_chi")
	email := c.PostForm("email")

	oldPassword := c.PostForm("mat_khau_cu")
	newPassword := c.PostForm("mat_khau_moi")
	confirmPassword := c.PostForm("xac_nhan_mat_khau_moi")

	// Cập nhật thông tin cơ bản
	if hoTen != "" {
		nv.HoTen = hoTen
	}
	if gioiTinh != "" {
		nv.GioiTinh = gioiTinh
	}
	if ngaySinh != "" {
		nv.NgaySinh = ngaySinh
	}
	if sdt != "" {
		nv.SDT = sdt
	}
	if diaChi != "" {
		nv.DiaChi = diaChi
	}
	if email != "" {
		nv.Email = email
	}

	// ✅ Đổi mật khẩu khi có nhập đủ 3 trường
	if oldPassword != "" || newPassword != "" || confirmPassword != "" {
		if oldPassword == "" || newPassword == "" || confirmPassword == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cần nhập đủ mật khẩu cũ, mật khẩu mới và xác nhận mật khẩu mới"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(nv.MatKhau), []byte(oldPassword)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Mật khẩu cũ không đúng"})
			return
		}

		if newPassword != confirmPassword {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Xác nhận mật khẩu mới không khớp"})
			return
		}

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		nv.MatKhau = string(hashedPassword)
	}

	// ✅ Upload ảnh mới (nếu có)
	file, err := c.FormFile("image")
	if err == nil && file != nil {
		src, _ := file.Open()
		defer src.Close()

		uploadResult, err := config.CLD.Upload.Upload(c, src, uploader.UploadParams{Folder: "nhanvien"})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Upload ảnh thất bại: " + err.Error()})
			return
		}

		config.DB.Where("owner_id = ? AND owner_type = ?", nv.MaNV, "nhan_vien").Delete(&models.Images{})

		newImg := models.Images{
			OwnerID:   nv.MaNV,
			OwnerType: "nhan_vien",
			ImageURL:  uploadResult.SecureURL,
		}
		config.DB.Create(&newImg)
	}

	// ✅ Lưu thay đổi
	if err := config.DB.Save(&nv).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể cập nhật thông tin cá nhân: " + err.Error()})
		return
	}

	config.DB.Preload("AnhNhanVien").First(&nv, nv.MaNV)

	c.JSON(http.StatusOK, gin.H{
		"message": "Cập nhật thông tin cá nhân thành công",
		"data":    nv,
	})
}
