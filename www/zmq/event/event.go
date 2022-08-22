package event

import (
	"bytes"

	"github.com/zarbchain/zarb-go/types/crypto/hash"
	"github.com/zarbchain/zarb-go/util/encoding"
)

type Event struct {
	Topic string
	Body []byte
}

func  CreateBlockEvent(blockHash hash.Hash, height uint32) Event {
	buf := make([]byte,0)
	w := bytes.NewBuffer(buf)
    encoding.WriteElement(w, blockHash)
	encoding.WriteElement(w, height)

	return Event{
		Topic: "block",
		Body: w.Bytes(),
	}

  
}