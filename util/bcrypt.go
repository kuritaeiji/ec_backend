package util

import "golang.org/x/crypto/bcrypt"

type bcryptUtils struct{}

var BcryptUtils = bcryptUtils{}

// bcryptでハッシュ化したパスワードを返却する
func (bu bcryptUtils) GeneratePasswordDigest(password string) (string, error) {
	passwordDigestByte, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(passwordDigestByte), nil
}
