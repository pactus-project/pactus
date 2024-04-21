package config

// NodeConfigError is returned when the config configuration is invalid.
type NodeConfigError struct {
	Reason string
}

func (e NodeConfigError) Error() string {
	return e.Reason
}
