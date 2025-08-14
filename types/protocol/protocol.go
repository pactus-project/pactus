package protocol

import (
	"fmt"
	"strconv"
)

type Version int8

const (
	ProtocolVersion1 Version = 1
	ProtocolVersion2 Version = 2

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
