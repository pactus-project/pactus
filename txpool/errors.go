package txpool

// ConfigError is returned when the txpool configuration is invalid.
type ConfigError struct {
	Reason string
}

func (e ConfigError) Error() string {
	return e.Reason
}
