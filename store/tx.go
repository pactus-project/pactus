package store

import (
	"bytes"

	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util/encoding"
	"github.com/syndtr/goleveldb/leveldb"
)

type blockRegion struct {
	height uint32
	offset uint32
	length uint32
}

func txKey(id tx.ID) []byte { return append(txPrefix, id.Bytes()...) }

type txStore struct {
	db *leveldb.DB
}

func newTxStore(db *leveldb.DB) *txStore {
	return &txStore{
		db: db,
	}
}

func (ts *txStore) saveTx(batch *leveldb.Batch, id tx.ID, reg *blockRegion) {
	w := bytes.NewBuffer(make([]byte, 0, 32+4))

	err := encoding.WriteElements(w, &reg.height, &reg.offset, &reg.length)
	if err != nil {
		panic(err)
	}

	txKey := txKey(id)
	batch.Put(txKey, w.Bytes())
}

func (ts *txStore) tx(id tx.ID) (*blockRegion, error) {
	data, err := tryGet(ts.db, txKey(id))
	if err != nil {
		return nil, err
	}

	r := bytes.NewReader(data)
	reg := new(blockRegion)
	err = encoding.ReadElements(r, &reg.height, &reg.offset, &reg.length)

	if err != nil {
		return nil, err
	}

	return reg, nil
}
