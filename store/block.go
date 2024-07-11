package store

import (
	"bytes"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/sortition"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/encoding"
	"github.com/pactus-project/pactus/util/pairslice"
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
	pubKeyCache        *lru.Cache[crypto.Address, *bls.PublicKey]
	sortitionSeedCache *pairslice.PairSlice[uint32, *sortition.VerifiableSeed]
	sortitionCacheSize uint32
}

func newBlockStore(db *leveldb.DB, sortitionCacheSize uint32, publicKeyCacheSize int) *blockStore {
	pubKeyCache, err := lru.New[crypto.Address, *bls.PublicKey](publicKeyCacheSize)
	if err != nil {
		return nil
	}

	return &blockStore{
		db:                 db,
		sortitionSeedCache: pairslice.New[uint32, *sortition.VerifiableSeed](int(sortitionCacheSize)),
		pubKeyCache:        pubKeyCache,
		sortitionCacheSize: sortitionCacheSize,
	}
}

func (bs *blockStore) saveBlock(batch *leveldb.Batch, height uint32, blk *block.Block) []blockRegion {
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
	bs.saveToCache(height, sortitionSeed)

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

func (bs *blockStore) sortitionSeed(blockHeight uint32) *sortition.VerifiableSeed {
	startHeight, _, _ := bs.sortitionSeedCache.First()

	if blockHeight < startHeight {
		return nil
	}

	index := blockHeight - startHeight
	_, sortitionSeed, ok := bs.sortitionSeedCache.Get(int(index))
	if !ok {
		return nil
	}

	return sortitionSeed
}

func (bs *blockStore) hasBlock(height uint32) bool {
	return tryHas(bs.db, blockKey(height))
}

func (bs *blockStore) publicKey(addr crypto.Address) (*bls.PublicKey, error) {
	if pubKey, ok := bs.pubKeyCache.Get(addr); ok {
		return pubKey, nil
	}

	data, err := tryGet(bs.db, publicKeyKey(addr))
	if err != nil {
		return nil, err
	}
	pubKey, err := bls.PublicKeyFromBytes(data)
	if err != nil {
		return nil, err
	}

	bs.pubKeyCache.Add(addr, pubKey)

	return pubKey, err
}

func (bs *blockStore) hasPublicKey(addr crypto.Address) bool {
	ok := bs.pubKeyCache.Contains(addr)
	if !ok {
		ok = tryHas(bs.db, publicKeyKey(addr))
	}

	return ok
}

func (bs *blockStore) saveToCache(blockHeight uint32, sortitionSeed sortition.VerifiableSeed) {
	bs.sortitionSeedCache.Append(blockHeight, &sortitionSeed)
	if bs.sortitionSeedCache.Len() > int(bs.sortitionCacheSize) {
		bs.sortitionSeedCache.RemoveFirst()
	}
}
