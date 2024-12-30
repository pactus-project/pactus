package zmq

type Topic int16

const (
	BlockInfo       Topic = 0x0001
	TransactionInfo Topic = 0x0002
	RawBlock        Topic = 0x0003
	RawTransaction  Topic = 0x0004
)

func (t Topic) String() string {
	switch t {
	case BlockInfo:
		return "block_info"

	case TransactionInfo:
		return "transaction_info"

	case RawBlock:
		return "raw_block"

	case RawTransaction:
		return "raw_transaction"

	default:
		return ""
	}
}
