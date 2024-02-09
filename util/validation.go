package util

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/kuritaeiji/ec_backend/share"
	"github.com/labstack/echo/v4"
)

var Validate = validator.New()

const (
	RequriedMsg               = "%sは必須です"
	LteMsg                    = "%sは%s文字以下にしてください"
	GteMsg                    = "%sは%s文字以上にしてください"
	AvailableSymbolPaswordMsg = "%sはアルファベット・数字・%sのみ使用できます"
	EmailMsg                  = "%sはメールアドレスとして正しい形式ではありません"
	PasswordConfirmationMsg   = "%sが一致しません"
)

type ValidationUtils interface {
	Struct(s interface{}) error
	CreateValidationMessages(err error) error
}

type validationUtils struct {
	*validator.Validate
	logger echo.Logger
}

func NewValidationUtils(logger echo.Logger) validationUtils {
	return validationUtils{
		Validate: Validate,
		logger:   logger,
	}
}

func (vu validationUtils) CreateValidationMessages(err error) error {
	if err == nil {
		return nil
	}

	vErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return nil
	}

	msgs := []string{}

	for _, vErr := range vErrs {
		fieldName, ok := FieldNames[vErr.Namespace()]
		if !ok {
			vu.logger.Debug(fmt.Sprintf("%sの日本語フィールド名がfieldNamesに登録されていません", vErr.Namespace()))
			continue
		}

		tag := vErr.Tag()
		param := vErr.Param()

		switch tag {
		case "required":
			msgs = append(msgs, fmt.Sprintf(RequriedMsg, fieldName))
		case "lte":
			msgs = append(msgs, fmt.Sprintf(LteMsg, fieldName, param))
		case "gte":
			msgs = append(msgs, fmt.Sprintf(GteMsg, fieldName, param))
		case "available_symbol_password":
			msgs = append(msgs, fmt.Sprintf(AvailableSymbolPaswordMsg, fieldName, param))
		case "email":
			msgs = append(msgs, fmt.Sprintf(EmailMsg, fieldName))
		case "password_confirmation":
			msgs = append(msgs, fmt.Sprintf(PasswordConfirmationMsg, fieldName))
		}
	}

	return share.OriginalError{
		Code:     share.ErrorCodeValidation,
		Messages: msgs,
	}
}
