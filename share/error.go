package share

import "strings"

type (
	ErrorCode int

	OriginalError struct {
		Code     ErrorCode `json:"code"`
		Messages []string  `json:"messages"`
	}
)

const (
	ErrorCodeValidation ErrorCode = iota + 2 // ResultCodeにはSuccessが存在しSuccessが1なので2から始めている
	ErrorCodeNoLogin
	ErrorCodeOther
)

func CreateOriginalError(code ErrorCode, messages []string) OriginalError {
	return OriginalError{
		Code: code,
		Messages: messages,
	}
}

func (oe OriginalError) Error() string {
	return strings.Join(oe.Messages, "\n")
}
