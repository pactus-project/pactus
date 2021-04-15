package store

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/zarbchain/zarb-go/logger"
)

func tryGet(db *leveldb.DB, key []byte) ([]byte, error) {
	data, err := db.Get(key, nil)
	if err != nil {
		// Probably key doesn't exist in database
		logger.Trace("DB error on get", "err", err, "key", key)
		return nil, err
	}
	return data, nil
}
