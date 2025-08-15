package jwt

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("your_secret_key") // 推荐后续用 .env 管理

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func init() {
	//JWT_SECRET_KEY
	jwtKey := os.Getenv("JWT_SECRET_KEY")
	if jwtKey == "" {
		jwtKey = "your_secret_key"
	}
}

// GenerateToken 生成 JWT Token
func GenerateToken(userID uint) (string, error) {
	expirationTime := time.Now().Add(72 * time.Hour) // 三天有效
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// ParseToken 验证 JWT Token
func ParseToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("token 无效或已过期")
	}
	return claims, nil
}
