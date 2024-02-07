package enum

type AuthType int

const (
	AuthTypeEmail AuthType = iota + 1
	AuthTypeGoogle
	AuthTypeApple
)
