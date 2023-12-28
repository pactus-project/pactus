package store

import (
	"bytes"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/sortition"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/encoding"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/syndtr/goleveldb/leveldb"
)

func blockKey(height uint32) []byte { return append(blockPrefix, util.Uint32ToSlice(height)...) }
func publicKeyKey(addr crypto.Address) []byte {
	return append(publicKeyPrefix, addr.Bytes()...)
}

func blockHashKey(h hash.Hash) []byte {
	return append(blockHeightPrefix, h.Bytes()...)
}

type blockStore struct {
	db                 *leveldb.DB
	sortitionSeedCache []*sortition.VerifiableSeed
	sortitionInterval  uint32
}

func newBlockStore(db *leveldb.DB, sortitionCacheSize uint32) *blockStore {
	return &blockStore{
		db:                 db,
		sortitionSeedCache: make([]*sortition.VerifiableSeed, 0, sortitionCacheSize),
		sortitionInterval:  sortitionCacheSize,
	}
}

func (bs *blockStore) saveBlock(batch *leveldb.Batch, height uint32, blk *block.Block) []blockRegion {
	if height > 1 {
		if !bs.hasBlock(height - 1) {
			logger.Panic("previous block not found", "height", height)
		}
	}
	if bs.hasBlock(height) {
		logger.Panic("duplicated block", "height", height)
	}

	blockHash := blk.Hash()
	regs := make([]blockRegion, blk.Transactions().Len())
	w := bytes.NewBuffer(make([]byte, 0, blk.SerializeSize()+hash.HashSize))
	err := encoding.WriteElement(w, &blockHash)
	if err != nil {
		panic(err) // Should we panic?
	}
	err = blk.Header().Encode(w)
	if err != nil {
		panic(err) // Should we panic?
	}
	if blk.PrevCertificate() != nil {
		err = blk.PrevCertificate().Encode(w)
		if err != nil {
			panic(err) // Should we panic?
		}
	}
	err = encoding.WriteVarInt(w, uint64(blk.Transactions().Len()))
	if err != nil {
		panic(err) // Should we panic?
	}
	for i, trx := range blk.Transactions() {
		offset := w.Len()
		regs[i].height = height
		regs[i].offset = uint32(offset)

		pubKey := trx.PublicKey()
		if pubKey != nil {
			// TODO: improve my performance by caching public keys
			if !bs.hasPublicKey(trx.Payload().Signer()) {
				publicKeyKey := publicKeyKey(trx.Payload().Signer())
				batch.Put(publicKeyKey, pubKey.Bytes())
			} else {
				// we have indexed this public key, se we can remove it
				trx.SetPublicKey(nil)
			}
		}

		err := trx.Encode(w)
		if err != nil {
			panic(err) // Should we panic?
		}
		regs[i].length = uint32(w.Len() - offset)

		trx.SetPublicKey(pubKey)
	}
	blockKey := blockKey(height)
	blockHashKey := blockHashKey(blockHash)

	batch.Put(blockKey, w.Bytes())
	batch.Put(blockHashKey, util.Uint32ToSlice(height))

	sortitionSeed := blk.Header().SortitionSeed()
	bs.saveToCache(sortitionSeed)

	return regs
}

func (bs *blockStore) block(height uint32) ([]byte, error) {
	data, err := tryGet(bs.db, blockKey(height))
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (bs *blockStore) blockHeight(h hash.Hash) uint32 {
	data, err := tryGet(bs.db, blockHashKey(h))
	if err != nil {
		return 0
	}
	return util.SliceToUint32(data)
}

func (bs *blockStore) sortitionSeed(blockHeight, currentHeight uint32) *sortition.VerifiableSeed {
	index := currentHeight - blockHeight
	if index > bs.sortitionInterval {
		return nil
	}

	if index != 0 {
		index = uint32(len(bs.sortitionSeedCache)) - index
	} else {
		index = uint32(len(bs.sortitionSeedCache)) - 1
	}

	return bs.sortitionSeedCache[index]
}

func (bs *blockStore) hasBlock(height uint32) bool {
	return tryHas(bs.db, blockKey(height))
}

func (bs *blockStore) hasPublicKey(addr crypto.Address) bool {
	return tryHas(bs.db, publicKeyKey(addr))
}

func (bs *blockStore) saveToCache(sortitionSeed sortition.VerifiableSeed) {
	bs.sortitionSeedCache = append(bs.sortitionSeedCache, &sortitionSeed)
	if len(bs.sortitionSeedCache) > int(bs.sortitionInterval) {
		bs.sortitionSeedCache = bs.sortitionSeedCache[1:]
	}
}
