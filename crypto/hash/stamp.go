package hash

import (
	"encoding/hex"
	"fmt"
)

const StampSize = 4

type Stamp [StampSize]byte

func StampFromString(str string) (Stamp, error) {
	data, err := hex.DecodeString(str)
	if err != nil {
		return Stamp{}, err
	}
	if len(data) != StampSize {
		return Stamp{}, fmt.Errorf("Stamp should be %d bytes, but it is %v bytes", StampSize, len(data))
	}
	var s Stamp
	copy(s[:], data[:StampSize])
	return s, nil
}

func (s Stamp) Bytes() []byte {
	return s[:]
}

func (s Stamp) String() string {
	return hex.EncodeToString(s[:])
}

func (s Stamp) EqualsTo(r Stamp) bool {
	return s == r
}
