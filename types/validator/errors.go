package validator

import "errors"

// ErrInvalidData is returned when the data of a validator is not canonical.
var ErrInvalidData = errors.New("invalid validator data")
