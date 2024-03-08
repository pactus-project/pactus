package config

// ConfigError is returned when the store configuration is invalid.
type ConfigError struct {
	Reason string
}

func (e ConfigError) Error() string {
	return e.Reason
}
