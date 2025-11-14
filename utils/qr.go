package utils

import (
	"encoding/base64"
	"fmt"

	"github.com/skip2/go-qrcode"
	"github.com/subiz/vietqr"
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

//amount := 12000.0
//bankBIN := "970423"
//account := "00005897596"
//note := "Ủng hộ lũ lụt"

func GenerateQRPayment(amount float64, bankBIN string, accountnumber, note string) (string, error) {
	code := vietqr.Generate(amount, bankBIN, accountnumber, note)

	png, err := qrcode.Encode(code, qrcode.Medium, 256)
	if err != nil {
		return "", err
	}

	// Chuyển []byte → base64 string
	base64Img := base64.StdEncoding.EncodeToString(png)

	return base64Img, nil
}
