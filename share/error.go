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
	ErrorCodeValidation ErrorCode = iota + 2
)

func (oe OriginalError) Error() string {
	return strings.Join(oe.Messages, "\n")
}
