package errors

import (
	"errors"
)

var (
	ErrRequestLong = errors.New("The request data is too long")
)
