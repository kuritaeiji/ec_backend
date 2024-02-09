package share

type (
	ResultCode int

	Result struct {
		Code     ResultCode `json:"code"`
		Messages []string   `json:"messages"`
	}
)

const (
	ResultCodeSuccess ResultCode = iota + 1
	ResultCodeValidation
)

func OriginalErrorToResult(err OriginalError) Result {
	return Result{
		Code:     ResultCode(err.Code),
		Messages: err.Messages,
	}
}

func SuccessResult() Result {
	return Result{
		Code:     ResultCodeSuccess,
		Messages: []string{},
	}
}
