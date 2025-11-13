package utils

import (
	"fmt"

	"github.com/skip2/go-qrcode"
)

func GenerateQRBytes(maBan int, tenBan string, soChoNgoi int, trangThai string) ([]byte, error) {
	content := fmt.Sprintf(
		"Tên bàn: %s\nMã bàn: %d\nSố chỗ ngồi: %d\nTrạng thái: %s",
		tenBan, maBan, soChoNgoi, trangThai,
	)

	// ✅ Tạo QR code trong bộ nhớ (RAM)
	png, err := qrcode.Encode(content, qrcode.Medium, 256)
	if err != nil {
		return nil, err
	}
	return png, nil
}
