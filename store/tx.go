package store

import (
	"github.com/fxamacker/cbor/v2"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/tx"
)

var (
	txPrefix = []byte{0x01}
)

func txKey(hash crypto.Hash) []byte { return append(txPrefix, hash.RawBytes()...) }

type txStore struct {
	db *leveldb.DB
}

func newTxStore(path string) (*txStore, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &txStore{
		db: db,
	}, nil
}

func (ts *txStore) close() error {
	return ts.db.Close()
}

func (ts *txStore) saveTx(ctrs tx.CommittedTx) error {
	if err := ctrs.SanityCheck(); err != nil {
		return err
	}
	data, err := cbor.Marshal(ctrs)
	if err != nil {
		return err
	}
	txKey := txKey(ctrs.Tx.Hash())
	err = tryPut(ts.db, txKey, data)
	if err != nil {
		return err
	}
	return nil
}

func (ts *txStore) tx(hash crypto.Hash) (*tx.CommittedTx, error) {
	txKey := txKey(hash)
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
