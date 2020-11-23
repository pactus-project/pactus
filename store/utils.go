package store

import (
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/zarbchain/zarb-go/logger"
)

func tryGet(db *leveldb.DB, key []byte) ([]byte, error) {
	data, err := db.Get(key, nil)
	if err != nil {
		logger.Error("DB error", "err", err, "key", key)
		return nil, err
	}
	if data == nil {
		return nil, fmt.Errorf("No value is associated with this key : %x", key)
	}
	return data, nil
}

func tryPut(db *leveldb.DB, key, value []byte) error {
	has, err := db.Has(key, nil)
	if has {
		logger.Trace("The key exists in database, update it.", "key", key)
	}
	err = db.Put(key, value, nil)
	if err != nil {
		logger.Error("DB error", "err", err)
		return err
	}

	return nil
}
