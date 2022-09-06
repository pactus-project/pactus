package util

import (
	"encoding/hex"
)

// CopyBytes returns an exact copy of the provided bytes.
func CopyBytes(b []byte) (copiedBytes []byte) {
	if b == nil {
		return nil
	}
	copiedBytes = make([]byte, len(b))
	copy(copiedBytes, b)

	return
}

// has0xPrefix validates str begins with '0x' or '0X'.
func has0xPrefix(str string) bool {
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}

// Bytes2Hex returns the hexadecimal encoding of d.
func Bytes2Hex(d []byte) string {
	return hex.EncodeToString(d)
}

// FromHex returns the bytes represented by the hexadecimal string str.
// str may be prefixed with "0x".
func Hex2Bytes(str string) []byte {
	if has0xPrefix(str) {
		str = str[2:]
	}
	if len(str)%2 == 1 {
		str = "0" + str
	}
	h, _ := hex.DecodeString(str)
	return h
}
