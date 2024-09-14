package downloader

import (
	"fmt"
)

type Error struct {
	Message string
	Reason  error
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Message, e.Reason.Error())
}

func (e *Error) Unwrap() error {
	return e.Reason
}
