package peerset

type StatusCode int

const (
	StatusCodeBanned       = StatusCode(-1)
	StatusCodeUnknown      = StatusCode(0)
	StatusCodeDisconnected = StatusCode(1)
	StatusCodeConnected    = StatusCode(2)
	StatusCodeKnown        = StatusCode(3)
	StatusCodeTrusty       = StatusCode(4)
)

func (code StatusCode) String() string {
	switch code {
	case StatusCodeBanned:
		return "banned"
	case StatusCodeDisconnected:
		return "disconnected"
	case StatusCodeConnected:
		return "connected"
	case StatusCodeUnknown:
		return "unknown"
	case StatusCodeKnown:
		return "known"
	case StatusCodeTrusty:
		return "trusty"
	}
	return "invalid"
}
