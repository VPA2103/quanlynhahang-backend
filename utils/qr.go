package utils

import (
	"fmt"
	"log"

	"github.com/skip2/go-qrcode"
)

func GenerateQR(maBan int, tenBan string, soChoNgoi int, trangThai string) (string, error) {
	content := fmt.Sprintf("Tên bàn: %s\nMã bàn: %d\nSố chỗ ngồi: %d\nTrạng thái: %s",
		tenBan, maBan, soChoNgoi, trangThai)

	filename := fmt.Sprintf("qr_ban%d.png", maBan)

	err := qrcode.WriteFile(content, qrcode.Medium, 256, filename)
	if err != nil {
		log.Println("❌ Lỗi tạo QR:", err)
		return "", err
	}

	// Có thể trả về đường dẫn local hoặc upload lên Cloudinary sau này
	return filename, nil
}
