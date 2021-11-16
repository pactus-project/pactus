package store

import (
	"github.com/syndtr/goleveldb/leveldb"
	dbutil "github.com/syndtr/goleveldb/leveldb/util"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/util"
)

func blockKey(height int) []byte         { return append(blockPrefix, util.IntToSlice(height)...) }
func blockHashKey(hash hash.Hash) []byte { return append(blockHashPrefix, hash.RawBytes()...) }

type blockStore struct {
	db *leveldb.DB
}

func newBlockStore(db *leveldb.DB) *blockStore {
	return &blockStore{
		db: db,
	}
}

func (bs *blockStore) saveBlock(batch *leveldb.Batch, height int, block *block.Block) error {
	blockData, err := block.Encode()
	if err != nil {
		return err
	}
	blockKey := blockKey(height)
	blockHashKey := blockHashKey(block.Hash())

	batch.Put(blockKey, blockData)
	batch.Put(blockHashKey, util.IntToSlice(height))

	return nil
}

func (bs *blockStore) block(height int) (*block.Block, error) {
	blockKey := blockKey(height)
	data, err := tryGet(bs.db, blockKey)
	if err != nil {
		return nil, err
	}
	block := new(block.Block)
	err = block.Decode(data)
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (bs *blockStore) blockHeight(hash hash.Hash) (int, error) {
	blockHashKey := blockHashKey(hash)
	heightData, err := tryGet(bs.db, blockHashKey)
	if err != nil {
		return -1, err
	}
	return util.SliceToInt(heightData), nil
}

func (bs *blockStore) hasAnyBlock() bool {
	iter := bs.db.NewIterator(dbutil.BytesPrefix(blockHashPrefix), nil)
	return iter.First()
}
