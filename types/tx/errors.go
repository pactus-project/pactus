package tx

import "github.com/pactus-project/pactus/types/tx/payload"

// BasicCheckError is returned when the basic check on the certificate fails.
type BasicCheckError struct {
	Reason string
}

func (e BasicCheckError) Error() string {
	return "transaction basic check failed: " + e.Reason
}

// InvalidPayloadTypeError is returned when the payload type is not valid.
type InvalidPayloadTypeError struct {
	PayloadType payload.Type
}

func (e InvalidPayloadTypeError) Error() string {
	return "invalid payload type: " + e.PayloadType.String()
}
