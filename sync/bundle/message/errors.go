package message

import "fmt"

// BasicCheckError is returned when the basic check on the message fails.
type BasicCheckError struct {
	Reason string
}

func (e BasicCheckError) Error() string {
	return e.Reason
}

// InvalidMessageTypeError is returned when the message type is not valid.
type InvalidMessageTypeError struct {
	Type int
}

func (e InvalidMessageTypeError) Error() string {
	return fmt.Sprintf("invalid message type: %d", e.Type)
}
