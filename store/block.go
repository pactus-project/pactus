package store

import (
	"github.com/fxamacker/cbor/v2"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/libs/linkedmap"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/util"
)

func blockKey(height int) []byte         { return append(blockPrefix, util.IntToSlice(height)...) }
func blockHashKey(hash hash.Hash) []byte { return append(blockHashPrefix, hash.RawBytes()...) }

type blockStore struct {
	db            *leveldb.DB
	stampToHeight *linkedmap.LinkedMap
}

func newBlockStore(db *leveldb.DB, stampToHeightCapacity int) *blockStore {
	bs := &blockStore{
		db:            db,
		stampToHeight: linkedmap.NewLinkedMap(stampToHeightCapacity),
	}

	height, _ := bs.lastCertificate()

	// Add genesis block stamp
	bs.stampToHeight.PushFront(hash.UndefHash.Stamp(), 0)

	for i := 0; i < stampToHeightCapacity; i++ {
		if height-i > 0 {
			b, _ := bs.block(height - i)
			bs.stampToHeight.PushFront(b.Stamp(), height-i)
		}
	}

	return bs
}

func (bs *blockStore) saveBlock(batch *leveldb.Batch, height int, block *block.Block, cert *block.Certificate) {
	if height > 1 {
		_, err := tryGet(bs.db, blockKey(height-1))
		if err != nil {
			logger.Panic("previous block not found: %v", height)
		}
	}
	blockData, err := block.Encode()
	if err != nil {
		logger.Panic("unable to encode block: %v", err)
	}
	blockKey := blockKey(height)
	blockHashKey := blockHashKey(block.Hash())

	batch.Put(blockKey, blockData)
	batch.Put(blockHashKey, util.IntToSlice(height))

	// Save last certificate
	lc := lastCertificate{
		Height: height,
		Cert:   cert,
	}
	data, err := cbor.Marshal(lc)
	if err != nil {
		logger.Panic("unable to encode last certificate: %v", err)
	}
	batch.Put(lastCertKey, data)

	// Update stamp to height lookup
	bs.stampToHeight.PushBack(block.Stamp(), height)
}

func (bs *blockStore) block(height int) (*block.Block, error) {
	data, err := tryGet(bs.db, blockKey(height))
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

func (bs *blockStore) blockHeightByStamp(stamp hash.Stamp) int {
	v, ok := bs.stampToHeight.Get(stamp)
	if ok {
		return v.(int)
	}
	return -1
}

func (bs *blockStore) lastCertificate() (int, *block.Certificate) {
	data, _ := tryGet(bs.db, lastCertKey)
	lc := new(lastCertificate)
	err := cbor.Unmarshal(data, lc)
	if err != nil {
		// Genesis block
		return 0, nil
	}
	return lc.Height, lc.Cert
}
