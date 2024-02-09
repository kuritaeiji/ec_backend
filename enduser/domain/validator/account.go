package validator

import (
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
	"github.com/kuritaeiji/ec_backend/enduser/domain/enum"
	"github.com/kuritaeiji/ec_backend/util"
)

// バリデーター登録
func init() {
	util.Validate.RegisterStructValidation(passwordValidator, ValidationAccountForCreation{})
}

// アカウント登録時のバリデーション用アカウント構造体
// アカウントエンティティーのパスワードはハッシュ化されたパスワードのためバリデーション用のオブジェクトを作成する
type ValidationAccountForCreation struct {
	Email                string `validate:"required,lte=255,email"`
	Password             string
	PasswordConfirmation string
	AuthType             enum.AuthType
}

// レビュー投稿者名のバリデーションアカウント構造体
type ValidationAccountForReviewNickname struct {
	ReviewNickname string `validate:"required,lte=20"`
}

// パスワードに使用可能な記号
const (
	availableSymbolForPassword = `!"#$%&'()`
	passwordFieldName          = "Password"
)

// パスワードのバリデーター
// 認証タイプがメールアドレスの場合のみバリデーションを行う
// 必須・8文字以上50文字以下・アルファベット、数字、「!"#$%&'()」のみ使用可能
// パスワードとパスワード（確認用）が一致する
func passwordValidator(sl validator.StructLevel) {
	validationAccount := sl.Current().Interface().(ValidationAccountForCreation)

	if validationAccount.AuthType != enum.AuthTypeEmail {
		return
	}

	password := validationAccount.Password
	if password == "" {
		sl.ReportError(validationAccount.Password, passwordFieldName, passwordFieldName, "required", "")
		return
	}

	if len(password) < 8 {
		sl.ReportError(validationAccount.Password, passwordFieldName, passwordFieldName, "gte", "8")
		return
	}

	if len(password) > 50 {
		sl.ReportError(validationAccount.Password, passwordFieldName, passwordFieldName, "lte", "50")
		return
	}

	for _, char := range password {
		if !isAlphabet(char) && !unicode.IsNumber(char) && !strings.ContainsRune(availableSymbolForPassword, char) {
			sl.ReportError(validationAccount.Password, passwordFieldName, passwordFieldName, "available_symbol_password", availableSymbolForPassword)
			return
		}
	}

	passwordConfirmation := validationAccount.PasswordConfirmation
	if password != passwordConfirmation {
		sl.ReportError(validationAccount.Password, passwordFieldName, passwordFieldName, "password_confirmation", "")
		return
	}
}

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func isAlphabet(r rune) bool {
	return strings.ContainsRune(alphabet, r)
}
