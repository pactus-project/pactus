package store

import (
	"github.com/fxamacker/cbor/v2"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/tx"
)

var (
	txPrefix = []byte{0x01}
)

func txKey(hash crypto.Hash) []byte { return append(txPrefix, hash.RawBytes()...) }

type txWithReceipt struct {
	Tx      tx.Tx      `cbor:"1,keyasint"`
	Receipt tx.Receipt `cbor:"2,keyasint"`
}

type txStore struct {
	db     *leveldb.DB
	logger *logger.Logger
}

func newTxStore(path string, logger *logger.Logger) (*txStore, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &txStore{
		db:     db,
		logger: logger,
	}, nil
}

func (bs *txStore) SaveTx(tx tx.Tx, receipt tx.Receipt) error {
	tnr := txWithReceipt{tx, receipt}
	data, err := cbor.Marshal(tnr)
	if err != nil {
		return err
	}
	txKey := txKey(tx.Hash())
	err = bs.db.Put(txKey, data, nil)
	if err != nil {
		return err
	}
	return nil
}

func (bs *txStore) Tx(hash crypto.Hash) (*tx.Tx, *tx.Receipt, error) {
	txKey := txKey(hash)
	data, err := bs.db.Get(txKey, nil)
	if err != nil {
		return nil, nil, err
	}
	tnr := new(txWithReceipt)
	err = cbor.Unmarshal(data, &tnr)
	if err != nil {
		return nil, nil, err
	}
	return &tnr.Tx, &tnr.Receipt, nil
}
