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

func blockKey(height uint32) []byte { return append(blockPrefix, util.Uint32ToSlice(height)...) }
func blockHashKey(hash hash.Hash) []byte {
	return append(blockHeightPrefix, hash.Bytes()...)
}

type blockStore struct {
	db *leveldb.DB
}

func newBlockStore(db *leveldb.DB) *blockStore {
	return &blockStore{
		db: db,
	}
}

func (bs *blockStore) saveBlock(batch *leveldb.Batch, height uint32, block *block.Block) []blockRegion {
	if height > 1 {
		if !bs.hasBlock(height - 1) {
			logger.Panic("previous block not found: %v", height)
		}
	}
	if bs.hasBlock(height) {
		logger.Panic("duplicated block: %v", height)
	}

	blockHash := block.Hash()
	regs := make([]blockRegion, block.Transactions().Len())
	w := bytes.NewBuffer(make([]byte, 0, block.SerializeSize()+hash.HashSize))
	err := encoding.WriteElement(w, &blockHash)
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
	for i, trx := range block.Transactions() {
		offset := w.Len()
		regs[i].height = height
		regs[i].offset = uint32(offset)

		err := trx.Encode(w)
		if err != nil {
			panic(err) // Should we panic?
		}
		regs[i].length = uint32(w.Len() - offset)
	}
	blockKey := blockKey(height)
	blockHashKey := blockHashKey(blockHash)

	batch.Put(blockKey, w.Bytes())
	batch.Put(blockHashKey, util.Uint32ToSlice(height))

	return regs
}

func (bs *blockStore) block(height uint32) ([]byte, error) {
	data, err := tryGet(bs.db, blockKey(height))
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (bs *blockStore) BlockHeight(hash hash.Hash) uint32 {
	data, err := tryGet(bs.db, blockHashKey(hash))
	if err != nil {
		return 0
	}
	return util.SliceToUint32(data)
}

func (bs *blockStore) hasBlock(height uint32) bool {
	has, err := bs.db.Has(blockKey(height), nil)
	if err != nil {
		return false
	}
	return has
}
