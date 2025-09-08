package status

type Status int

const (
	// The Status here, tells us the status os the connection.
	StatusBanned       = Status(-1)
	StatusUnknown      = Status(0)
	StatusDisconnected = Status(1)
	StatusConnected    = Status(2)
	StatusKnown        = Status(3)
)

func (s Status) String() string {
	switch s {
	case StatusBanned:
		return "banned"
	case StatusUnknown:
		return "unknown"
	case StatusDisconnected:
		return "disconnected"
	case StatusConnected:
		return "connected"
	case StatusKnown:
		return "known"
	}

	return "invalid"
}

func (s Status) IsUnknown() bool {
	return s == StatusUnknown
}

func (s Status) IsKnown() bool {
	return s == StatusKnown
}

func (s Status) IsDisconnected() bool {
	return s == StatusDisconnected
}

func (s Status) IsConnected() bool {
	return s == StatusConnected
}

func (s Status) IsConnectedOrKnown() bool {
	return s == StatusConnected || s == StatusKnown
}

func (s Status) IsBanned() bool {
	return s == StatusBanned
}
