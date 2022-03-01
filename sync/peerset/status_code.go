package peerset

import "encoding/json"

type StatusCode int

const (
	StatusCodeBanned  = StatusCode(-1)
	StatusCodeUnknown = StatusCode(0)
	StatusCodeKnown   = StatusCode(1)
	StatusCodeTrusty  = StatusCode(2)
)

func (code StatusCode) String() string {
	switch code {
	case StatusCodeBanned:
		return "banned"
	case StatusCodeUnknown:
		return "unknown"
	case StatusCodeKnown:
		return "known"
	case StatusCodeTrusty:
		return "trusty"
	}
	return "invalid"
}

func (code StatusCode) MarshalJSON() ([]byte, error) {
	return json.Marshal(code.String())
}
