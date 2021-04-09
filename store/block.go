package store

import (
	"github.com/syndtr/goleveldb/leveldb"
	dbutil "github.com/syndtr/goleveldb/leveldb/util"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/util"
)

var (
	blockPrefix     = []byte{0x01}
	blockHashPrefix = []byte{0x03}
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

func (bs *blockStore) close() error {
	return bs.db.Close()
}

func (bs *blockStore) saveBlock(block *block.Block, height int) error {
	blockData, err := block.Encode()
	if err != nil {
		return err
	}
	blockKey := blockKey(height)
	blockHashKey := blockHashKey(block.Hash())
	err = tryPut(bs.db, blockKey, blockData)
	if err != nil {
		return err
	}
	err = tryPut(bs.db, blockHashKey, util.IntToSlice(height))
	if err != nil {
		return err
	}
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

func (bs *blockStore) blockHeight(hash crypto.Hash) (int, error) {
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
