package utils

import "os"

var jwtSecret = []byte(getSecret())

func getSecret() string {
	if key := os.Getenv("SECRET_KEY"); key != "" {
		return key
	}
	return "MY_DEFAULT_SECRET_KEY" // fallback nếu chưa có env
}

func JWTSecret() []byte {
	return jwtSecret
}
