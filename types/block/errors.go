package block

// BasicCheckError is returned when the basic check on the certificate fails.
type BasicCheckError struct {
	Reason string
}

func (e BasicCheckError) Error() string {
	return "block basic check failed: " + e.Reason
}
