package util

import (
	"os"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/golang-jwt/jwt/v5"
)

type jwtUtils struct{}

var JwtUtils = jwtUtils{}

func (ju jwtUtils) CreateJwt(subject string, expireDuration time.Duration) (string, error) {
	// シークレットキー
	secretKey := os.Getenv("JWT_SECRET_KEY")

	// トークンの有効期限を設定
	expirationTime := time.Now().Add(expireDuration)

	// トークンのペイロードを作成
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		Subject: subject,
	}

	// トークンを作成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// トークンを署名
	tokenString, err := token.SignedString([]byte(secretKey))
	return tokenString, errors.WithStack(err)
}
