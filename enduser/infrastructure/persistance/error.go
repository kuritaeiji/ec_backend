package persistance

import "github.com/cockroachdb/errors"

var (
	ErrOptimisticLocking = errors.New("楽観ロックエラー")
)
