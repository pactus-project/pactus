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

func (ts *txStore) saveTx(batch *leveldb.Batch, ctrs *tx.CommittedTx) error {
	if err := ctrs.SanityCheck(); err != nil {
		return err
	}
	data, err := cbor.Marshal(ctrs)
	if err != nil {
		return err
	}
	txKey := txKey(ctrs.Tx.ID())
	batch.Put(txKey, data)

	if err != nil {
		return err
	}

	return nil
}

func (ts *txStore) tx(id tx.ID) (*tx.CommittedTx, error) {
	txKey := txKey(id)
	data, err := tryGet(ts.db, txKey)
	if err != nil {
		return nil, err
	}
	ctrs := new(tx.CommittedTx)
	err = cbor.Unmarshal(data, ctrs)
	if err != nil {
		return nil, err
	}
	if err := ctrs.SanityCheck(); err != nil {
		return nil, err
	}
	return ctrs, nil
}
