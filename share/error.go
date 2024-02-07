package share

import "strings"

type (
	OriginalErrorCode int

	OriginalError struct {
		Code     OriginalErrorCode `json:"code"`
		Messages []string `json:"messages"`
	}
)

const (
	ErrorCodeValidation OriginalErrorCode = iota + 1
)

func (oe OriginalError) Error() string {
	return strings.Join(oe.Messages, "\n")
}
