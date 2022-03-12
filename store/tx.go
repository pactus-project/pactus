package store

import (
	"github.com/fxamacker/cbor/v2"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/tx"
)

type txPos struct {
	Height int `cbor:"1,keyasint"`
	Index  int `cbor:"2,keyasint"`
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
	data, err := cbor.Marshal(pos)
	if err != nil {
		logger.Panic("unable to encode transaction: %v", err)
	}
	txKey := txKey(id)
	batch.Put(txKey, data)
}

func (ts *txStore) tx(id tx.ID) (*txPos, error) {
	data, err := tryGet(ts.db, txKey(id))
	if err != nil {
		return nil, err
	}
	pos := new(txPos)
	err = cbor.Unmarshal(data, pos)
	if err != nil {
		return nil, err
	}
	return pos, nil
}
