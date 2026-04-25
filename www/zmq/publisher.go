package zmq

import (
	"encoding/binary"

	"github.com/go-zeromq/zmq4"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/util/logger"
)

type Publisher interface {
	Address() string
	TopicName() string
	HWM() int

	onNewBlock(blk *block.Block)
}

type basePub struct {
	topic     Topic
	seqNo     uint32
	zmqSocket zmq4.Socket
	logger    *logger.SubLogger
}

func (b *basePub) Address() string {
	return b.zmqSocket.Addr().String()
}

func (b *basePub) TopicName() string {
	return b.topic.String()
}

func (b *basePub) HWM() int {
	hwmOpt, _ := b.zmqSocket.GetOption(zmq4.OptionHWM)

	return hwmOpt.(int)
}

// makeTopicMsg constructs a ZMQ message with a topic ID, message body, and sequence number.
// The message is constructed as a byte slice with the following structure:
// - Topic ID (2 Bytes)
// - Message body (varies based on provided parts)
// - Sequence number (4 Bytes).
func (b *basePub) makeTopicMsg(parts ...any) []byte {
	result := make([]byte, 0)

	// Append Topic ID to the message (2 Bytes)
	result = append(result, b.topic.Bytes()...)

	// Append message body based on the provided parts
	for _, part := range parts {
		switch castedVal := part.(type) {
		case crypto.Address:
			result = append(result, castedVal.Bytes()...)
		case types.Height:
			result = append(result, castedVal.EncodeAsSlice()...)
		case types.Round:
			result = append(result, castedVal.EncodeAsSlice()...)
		case []byte:
			result = append(result, castedVal...)
		case uint16:
			result = binary.LittleEndian.AppendUint16(result, castedVal)
		case uint32:
			result = binary.LittleEndian.AppendUint32(result, castedVal)
		default:
			panic("implement me!!")
		}
	}

	// Append sequence number to the message (4 Bytes, Little Endian encoding)
	result = binary.LittleEndian.AppendUint32(result, b.seqNo)

	return result
}
