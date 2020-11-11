package store

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/util"
)

var (
	blockPrefix     = []byte{0x01}
	blockHashPrefix = []byte{0x02}
)

func blockKey(height int) []byte           { return append(blockPrefix, util.IntToSlice(height)...) }
func blockHashKey(hash crypto.Hash) []byte { return append(blockHashPrefix, hash.RawBytes()...) }

type blockStore struct {
	db *leveldb.DB
}

func newBlockStore(path string) (*blockStore, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &blockStore{
		db: db,
	}, nil
}

func (bs *blockStore) SaveBlock(block block.Block, height int) error {
	blockData, err := block.Encode()
	if err != nil {
		return err
	}
	blockKey := blockKey(height)
	blockHashKey := blockHashKey(block.Hash())
	has, err := bs.db.Has(blockKey, nil)
	if has {
		// TODO: uncomment it later
		//logger.Warn("The blockkey exists in database, rewrite it.", "hash", block.Hash())
	}
	err = bs.db.Put(blockKey, blockData, nil)
	if err != nil {
		return err
	}
	err = bs.db.Put(blockHashKey, util.IntToSlice(height), nil)
	if err != nil {
		return err
	}
	return nil
}

func (bs *blockStore) RetrieveBlock(height int) (*block.Block, error) {
	blockKey := blockKey(height)
	blockData, err := bs.db.Get(blockKey, nil)
	if err != nil {
		return nil, err
	}
	block := new(block.Block)
	err = block.Decode(blockData)
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (bs *blockStore) blockHeight(hash crypto.Hash) (int, error) {
	blockHashKey := blockHashKey(hash)
	heightData, err := bs.db.Get(blockHashKey, nil)
	if err != nil {
		return -1, err
	}
	return util.SliceToInt(heightData), nil

}
