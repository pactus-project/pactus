package store

import (
	"bytes"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/encoding"
	"github.com/zarbchain/zarb-go/tx"
)

type txPos struct {
	Hash   hash.Hash
	Offset int32
}

func txKey(id tx.ID) []byte { return append(txPrefix, id.RawBytes()...) }

type txStore struct {
	db *leveldb.DB
}

func newTxStore(db *leveldb.DB) *txStore {
	return &txStore{
		db: db,
	}
}

func (ts *txStore) saveTx(batch *leveldb.Batch, id tx.ID, pos *txPos) {
	w := bytes.NewBuffer(make([]byte, 0, 32+4))
	err := encoding.WriteElements(w, &pos.Hash, &pos.Offset)
	if err != nil {
		panic(err)
	}

	txKey := txKey(id)
	batch.Put(txKey, w.Bytes())
}

func (ts *txStore) tx(id tx.ID) (*txPos, error) {
	data, err := tryGet(ts.db, txKey(id))
	if err != nil {
		return nil, err
	}
	r := bytes.NewReader(data)
	pos := new(txPos)
	err = encoding.ReadElements(r, &pos.Hash, &pos.Offset)
	if err != nil {
		return nil, err
	}
	return pos, nil
}
