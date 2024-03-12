package config

// Error is returned when the config configuration is invalid.
type Error struct {
	Reason string
}

func (e Error) Error() string {
	return e.Reason
}
