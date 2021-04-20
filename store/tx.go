package store

import (
	"github.com/fxamacker/cbor/v2"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/zarbchain/zarb-go/tx"
)

func txKey(id tx.ID) []byte { return append(txPrefix, id.RawBytes()...) }

type txStore struct {
	db *leveldb.DB
}

func newTxStore(db *leveldb.DB) (*txStore, error) {
	return &txStore{
		db: db,
	}, nil
}

func (ts *txStore) saveTx(batch *leveldb.Batch, trx *tx.Tx) error {
	data, err := cbor.Marshal(trx)
	if err != nil {
		return err
	}
	txKey := txKey(trx.ID())
	batch.Put(txKey, data)

	if err != nil {
		return err
	}

	return nil
}

func (ts *txStore) tx(id tx.ID) (*tx.Tx, error) {
	txKey := txKey(id)
	data, err := tryGet(ts.db, txKey)
	if err != nil {
		return nil, err
	}
	trx := new(tx.Tx)
	err = cbor.Unmarshal(data, trx)
	if err != nil {
		return nil, err
	}
	if err := trx.SanityCheck(); err != nil {
		return nil, err
	}
	return trx, nil
}
