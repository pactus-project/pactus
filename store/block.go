package store

import (
	"bytes"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/encoding"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/syndtr/goleveldb/leveldb"
)

func blockKey(hash hash.Hash) []byte { return append(blockPrefix, hash.Bytes()...) }
func blockHeightKey(height uint32) []byte {
	return append(blockHeightPrefix, util.Uint32ToSlice(height)...)
}

type blockStore struct {
	db *leveldb.DB
}

func newBlockStore(db *leveldb.DB) *blockStore {
	return &blockStore{
		db: db,
	}
}

func (bs *blockStore) saveBlock(batch *leveldb.Batch, height uint32, block *block.Block) []txPos {
	if height > 1 {
		if !bs.hasBlock(height - 1) {
			logger.Panic("previous block not found: %v", height)
		}
	}
	if bs.hasBlock(height) {
		logger.Panic("duplicated block: %v", height)
	}

	txsPos := make([]txPos, block.Transactions().Len())

	w := bytes.NewBuffer(make([]byte, 0, block.SerializeSize()+4))
	err := encoding.WriteElement(w, height)
	if err != nil {
		panic(err) // Should we panic?
	}
	err = block.Header().Encode(w)
	if err != nil {
		panic(err) // Should we panic?
	}
	if block.PrevCertificate() != nil {
		err = block.PrevCertificate().Encode(w)
		if err != nil {
			panic(err) // Should we panic?
		}
	}
	err = encoding.WriteVarInt(w, uint64(block.Transactions().Len()))
	if err != nil {
		panic(err) // Should we panic?
	}
	h := block.Hash()
	for i, trx := range block.Transactions() {
		txsPos[i].Hash = h
		txsPos[i].Offset = int32(w.Len())

		err := trx.Encode(w)
		if err != nil {
			panic(err) // Should we panic?
		}
	}
	blockKey := blockKey(block.Hash())
	blockHeightKey := blockHeightKey(height)

	batch.Put(blockKey, w.Bytes())
	batch.Put(blockHeightKey, block.Hash().Bytes())

	return txsPos
}

func (bs *blockStore) block(h hash.Hash) ([]byte, error) {
	data, err := tryGet(bs.db, blockKey(h))
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (bs *blockStore) BlockHash(height uint32) hash.Hash {
	// TODO: we can use flat file (height to hash) to reduce the size of level_db
	data, err := tryGet(bs.db, blockHeightKey(height))
	if err != nil {
		return hash.UndefHash
	}
	h, _ := hash.FromBytes(data)
	return h
}

func (bs *blockStore) hasBlock(height uint32) bool {
	has, err := bs.db.Has(blockHeightKey(height), nil)
	if err != nil {
		return false
	}
	return has
}
