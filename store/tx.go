package store

import (
	"bytes"

	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util/encoding"
	"github.com/pactus-project/pactus/util/linkedmap"
	"github.com/syndtr/goleveldb/leveldb"
)

type blockRegion struct {
	height uint32
	offset uint32
	length uint32
}

func txKey(txID tx.ID) []byte { return append(txPrefix, txID.Bytes()...) }

type txStore struct {
	db            *leveldb.DB
	txCache       *linkedmap.LinkedMap[tx.ID, uint32]
	txCacheWindow uint32
}

func newTxStore(db *leveldb.DB, txCacheWindow uint32) *txStore {
	return &txStore{
		db:            db,
		txCache:       linkedmap.New[tx.ID, uint32](0),
		txCacheWindow: txCacheWindow,
	}
}

func (ts *txStore) saveTxs(batch *leveldb.Batch, txs block.Txs, regs []blockRegion) {
	for i, trx := range txs {
		buf := bytes.NewBuffer(make([]byte, 0, 32+4))

		reg := regs[i]
		err := encoding.WriteElements(buf, &reg.height, &reg.offset, &reg.length)
		if err != nil {
			panic(err)
		}

		txID := trx.ID()
		key := txKey(txID)
		batch.Put(key, buf.Bytes())
		ts.addToCache(txID, reg.height)
	}
}

func (ts *txStore) pruneCache(currentHeight uint32) {
	for {
		head := ts.txCache.HeadNode()
		txHeight := head.Data.Value

		if currentHeight-txHeight <= ts.txCacheWindow {
			break
		}
		ts.txCache.RemoveHead()
	}
}

func (ts *txStore) recentTransaction(txID tx.ID) bool {
	return ts.txCache.Has(txID)
}

func (ts *txStore) tx(txID tx.ID) (*blockRegion, error) {
	data, err := tryGet(ts.db, txKey(txID))
	if err != nil {
		return nil, err
	}
	r := bytes.NewReader(data)
	reg := new(blockRegion)
	if err := encoding.ReadElements(r, &reg.height, &reg.offset, &reg.length); err != nil {
		return nil, err
	}

	return reg, nil
}

func (ts *txStore) addToCache(txID tx.ID, height uint32) {
	ts.txCache.PushBack(txID, height)
}
