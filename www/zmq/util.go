package zmq

import (
	"encoding/binary"

	"github.com/pactus-project/pactus/crypto"
)

func makeTopicMsg(parts ...any) []byte {
	result := make([]byte, 0, 64)

	for _, part := range parts {
		switch castedVal := part.(type) {
		case crypto.Address:
			result = append(result, castedVal.Bytes()...)
		case Topic:
			result = append(result, castedVal.Bytes()...)
		case []byte:
			result = append(result, castedVal...)
		case uint32:
			result = binary.BigEndian.AppendUint32(result, castedVal)
		case uint16:
			result = binary.BigEndian.AppendUint16(result, castedVal)
		default:
			panic("implement me!!")
		}
	}

	return result
}
