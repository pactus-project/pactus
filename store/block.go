package store

import (
	"fmt"

	"github.com/fxamacker/cbor/v2"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/util"
)

type blockInfo struct {
	Height           int      `cbor:"1,keyasint"`
	HeaderData       []byte   `cbor:"2,keyasint"`
	PrevCertData     []byte   `cbor:"3,keyasint"`
	TransactionsData [][]byte `cbor:"4,keyasint"`
}

func blockKey(hash hash.Hash) []byte   { return append(blockPrefix, hash.RawBytes()...) }
func blockHeightKey(height int) []byte { return append(blockHeightPrefix, util.IntToSlice(height)...) }

type blockStore struct {
	db *leveldb.DB
}

func newBlockStore(db *leveldb.DB) *blockStore {
	return &blockStore{
		db: db,
	}
}

func (bs *blockStore) saveBlock(batch *leveldb.Batch, height int, block *block.Block, cert *block.Certificate) {
	if height > 1 {
		if !bs.hasBlock(height - 1) {
			logger.Panic("previous block not found: %v", height)
		}
	}
	if bs.hasBlock(height) {
		logger.Panic("duplicated block: %v", height)
	}

	var headerData []byte
	var prevCertData []byte
	txsData := make([][]byte, block.Transactions().Len())

	headerData, _ = block.Header().Encode()
	if block.PrevCertificate() != nil {
		prevCertData, _ = block.PrevCertificate().Encode()
	}
	for i, trx := range block.Transactions() {
		txsData[i], _ = trx.Encode()
	}

	bi := blockInfo{
		Height:           height,
		HeaderData:       headerData,
		PrevCertData:     prevCertData,
		TransactionsData: txsData,
	}

	data, err := cbor.Marshal(bi)
	if err != nil {
		logger.Panic("unable to encode block: %v", err)
	}
	blockKey := blockKey(block.Hash())
	blockHeightKey := blockHeightKey(height)

	batch.Put(blockKey, data)
	batch.Put(blockHeightKey, block.Hash().RawBytes())
}

func (bs *blockStore) block(h hash.Hash) (*blockInfo, error) {
	data, err := tryGet(bs.db, blockKey(h))
	if err != nil {
		return nil, err
	}
	bi := new(blockInfo)
	err = cbor.Unmarshal(data, bi)
	if err != nil {
		return nil, err
	}

	blockHash := hash.CalcHash(bi.HeaderData)
	if blockHash != h {
		return nil, fmt.Errorf("header hash is not matched with key")
	}
	return bi, nil
}

func (bs *blockStore) BlockHash(height int) hash.Hash {
	// TODO: we can use flat file (height to hash) to reduce the size of level_db
	data, err := tryGet(bs.db, blockHeightKey(height))
	if err != nil {
		return hash.UndefHash
	}
	h, _ := hash.FromRawBytes(data)
	return h
}

func (bs *blockStore) hasBlock(height int) bool {
	has, err := bs.db.Has(blockHeightKey(height), nil)
	if err != nil {
		return false
	}
	return has
}
