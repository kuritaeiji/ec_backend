package util

import (
	"github.com/google/uuid"
)

var IDutils = idUtils{}

type idUtils struct {
}

func (iu idUtils) GenerateID() string {
	return uuid.NewString()
}
