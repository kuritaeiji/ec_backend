package persistance

import "github.com/cockroachdb/errors"

var (
	ErrSessionNotFound   = errors.New("セッションIDが見つかりません")
	ErrSessionExpired    = errors.New("有効期限が切れています")
	ErrOptimisticLocking = errors.New("楽観ロックエラー")
)
