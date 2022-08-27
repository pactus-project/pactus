package event

import (
	"bytes"

	"github.com/zarbchain/zarb-go/types/crypto/hash"
	"github.com/zarbchain/zarb-go/util/encoding"
	"github.com/zarbchain/zarb-go/util/logger"
)

type Event struct {
	Topic string
	Body  []byte
}

func CreateBlockEvent(blockHash hash.Hash, height uint32) Event {
	buf := make([]byte, 0)
	w := bytes.NewBuffer(buf)
	err := encoding.WriteElement(w, blockHash)
	if err != nil {
		logger.Error("error on block_hash event", "err", err)
	}
	err = encoding.WriteElement(w, height)
	if err != nil {
		logger.Error("error on height event", "err", err)
	}
	return Event{
		Topic: "block",
		Body:  w.Bytes(),
	}
}