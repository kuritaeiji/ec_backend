package util

import (
	"os"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/golang-jwt/jwt/v5"
)

type jwtUtils struct {
	secretKey []byte
}

var JwtUtils = jwtUtils{
	secretKey: []byte(os.Getenv("JWT_SECRET_KEY")),
}

func (ju jwtUtils) CreateJwt(subject string, expireDuration time.Duration) (string, error) {
	// トークンの有効期限を設定
	expirationTime := time.Now().Add(expireDuration)

	// トークンのペイロードを作成
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		Subject:   subject,
	}

	// トークンを作成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// トークンを署名
	tokenString, err := token.SignedString([]byte(ju.secretKey))
	return tokenString, errors.WithStack(err)
}

type (
	ErrTokenSignatureInvalid struct{}
	ErrTokenExpired          struct{}
)

func (err ErrTokenSignatureInvalid) Error() string { return "署名エラー" }
func (err ErrTokenExpired) Error() string          { return "有効期限切れエラー" }

// JWTを検証し、subjectを返却する
func (ju jwtUtils) ParseJwt(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return ju.secretKey, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			return "", errors.WithStack(ErrTokenSignatureInvalid{})
		}
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", errors.WithStack(ErrTokenExpired{})
		}
		return "", errors.WithStack(err)
	}

	sub, err := token.Claims.GetSubject()
	if err != nil {
		return "", errors.WithStack(err)
	}

	return sub, nil
}
