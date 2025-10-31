package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("MY_SECRET_KEY")

type JWTClaims struct {
	UserID   uint   `json:"ma_nv"`
	Username string `json:"username"`
	Role     string `json:"role"` // 🔥 thêm dòng này
	jwt.RegisteredClaims
}

// ✅ Sinh token có cả quyền (role)
func GenerateToken(userID uint, username string, role string) (string, error) {
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		Role:     role, // 🔥 gán role vào token
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // hết hạn sau 1 ngày
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// ✅ Hàm dùng trong middleware để xác thực token
func SecretKey() []byte {
	return secretKey
}
