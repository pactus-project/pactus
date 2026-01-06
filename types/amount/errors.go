package amount

import (
	"errors"
)

// ErrInvalidSQLType is returned when the type of the data
// is not supported for SQL database operations.
var ErrInvalidSQLType = errors.New("invalid SQL type")

// ErrInvalidYAMLType is returned when the type of the data
// is not supported for YAML unmarshaling.
var ErrInvalidYAMLType = errors.New("invalid YAML type")
