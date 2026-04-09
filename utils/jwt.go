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

func GenerateToken(id uint, email string, role string) (string, error) {
	claims := jwt.MapClaims{
		"id":    id,
		"email": email,
		"role":  role,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecret())
}

// ✅ Hàm dùng trong middleware để xác thực token
func SecretKey() []byte {
	return secretKey
}
