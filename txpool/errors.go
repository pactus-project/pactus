package txpool

import "fmt"

// ConfigError is returned when the txPool configuration is invalid.
type ConfigError struct {
	Reason string
}

func (e ConfigError) Error() string {
	return e.Reason
}

// AppendError is returned when the txPool configuration is invalid.
type AppendError struct {
	Err error
}

func (e AppendError) Error() string {
	return fmt.Sprintf("unable to append transaction to pool: %s", e.Err)
}
