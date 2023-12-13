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

func txKey(id tx.ID) []byte { return append(txPrefix, id.Bytes()...) }

type txStore struct {
	db      *leveldb.DB
	txCache *linkedmap.LinkedMap[tx.ID, uint32]
	ttl     uint32
}

func newTxStore(db *leveldb.DB, ttl uint32) *txStore {
	return &txStore{
		db:      db,
		txCache: linkedmap.NewLinkedMap[tx.ID, uint32](0),
		ttl:     ttl,
	}
}

func (ts *txStore) saveTxs(batch *leveldb.Batch, txs block.Txs, regs []blockRegion) {
	w := bytes.NewBuffer(make([]byte, 0, 32+4))

	for i, trx := range txs {
		reg := regs[i]
		err := encoding.WriteElements(w, &reg.height, &reg.offset, &reg.length)
		if err != nil {
			panic(err)
		}

		id := trx.ID()
		txKey := txKey(id)
		batch.Put(txKey, w.Bytes())
		ts.saveToCache(id, reg.height)
	}
}

func (ts *txStore) pruneCache(currentHeight uint32) {
	for {
		head := ts.txCache.HeadNode()
		txHeight := head.Data.Value

		if currentHeight-txHeight <= ts.ttl {
			break
		}
		ts.txCache.RemoveHead()
	}
}

func (ts *txStore) hasTX(id tx.ID) bool {
	return ts.txCache.Has(id)
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

func (ts *txStore) saveToCache(id tx.ID, height uint32) {
	ts.txCache.PushBack(id, height)
}
