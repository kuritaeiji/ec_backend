package util

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/kuritaeiji/ec_backend/share"
	"github.com/labstack/echo/v4"
)

var Validate = validator.New()

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
		fieldName, ok := fieldNames[vErr.Namespace()]
		if !ok {
			vu.logger.Debug(fmt.Sprintf("%sの日本語フィールド名がfieldNamesに登録されていません", vErr.Namespace()))
			continue
		}

		tag := vErr.Tag()
		param := vErr.Param()

		switch tag {
		case "requried":
			msgs = append(msgs, fmt.Sprintf("%sは必須です", fieldName))
		case "lte":
			msgs = append(msgs, fmt.Sprintf("%sは%s文字以下にしてください", fieldName, param))
		case "gte":
			msgs = append(msgs, fmt.Sprintf("%sは%s文字以上にしてください", fieldName, param))
		case "available_symbol_pasword":
			msgs = append(msgs, fmt.Sprintf(`%sはアルファベット・数字・%sのみ使用できます`, fieldName, param))
		case "email":
			msgs = append(msgs, fmt.Sprintf(`%sはメールアドレスとして正しい形式ではありません`, fieldName))
		case "password_confirmation":
			msgs = append(msgs, fmt.Sprintf("%sが一致しません", fieldName))
		}
	}

	return share.OriginalError{
		Code:     share.ErrorCodeValidation,
		Messages: msgs,
	}
}
