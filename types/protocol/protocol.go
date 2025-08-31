package protocol

import (
	"fmt"
	"strconv"
)

type Version uint8

const (
	ProtocolVersionUnknown Version = 0
	ProtocolVersion1       Version = 1 // Initial version
	ProtocolVersion2       Version = 2 // Split Reward Fork (PIP-43)

	ProtocolVersionLatest = ProtocolVersion2
)

func ParseVersion(s string) (Version, error) {
	v, err := strconv.ParseInt(s, 10, 8)
	if err != nil {
		return 0, err
	}

	return Version(v), nil
}

func (v Version) String() string {
	return fmt.Sprintf("%d", v)
}
