package genesis

// InvalidGenesisError is returned when a genesis document is not valid.
type InvalidGenesisError struct {
	Reason string
}

func (e InvalidGenesisError) Error() string {
	return "invalid genesis: " + e.Reason
}
