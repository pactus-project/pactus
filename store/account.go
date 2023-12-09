package store

import (
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

const (
	accLruCacheSize = 1024
)

type accountStore struct {
	db           *leveldb.DB
	addrLruCache *lru.Cache[crypto.Address, *account.Account]
	total        int32
}

func accountKey(addr crypto.Address) []byte { return append(accountPrefix, addr.Bytes()...) }

func newAccountStore(db *leveldb.DB) *accountStore {
	total := int32(0)
	addrLruCache, err := lru.New[crypto.Address, *account.Account](accLruCacheSize)
	if err != nil {
		logger.Panic("unable to create new instance of lru cache", "error", err)
	}

	r := util.BytesPrefix(accountPrefix)
	iter := db.NewIterator(r, nil)
	for iter.Next() {
		total++
	}
	iter.Release()

	return &accountStore{
		db:           db,
		total:        total,
		addrLruCache: addrLruCache,
	}
}

func (as *accountStore) hasAccount(addr crypto.Address) bool {
	ok := as.addrLruCache.Contains(addr)
	if !ok {
		ok = tryHas(as.db, accountKey(addr))
	}
	return ok
}

func (as *accountStore) account(addr crypto.Address) (*account.Account, error) {
	acc, ok := as.addrLruCache.Get(addr)
	if ok {
		return acc.Clone(), nil
	}

	rawData, err := tryGet(as.db, accountKey(addr))
	if err != nil {
		return nil, err
	}
	return account.FromBytes(rawData)
}

func (as *accountStore) iterateAccounts(consumer func(crypto.Address, *account.Account) (stop bool)) {
	r := util.BytesPrefix(accountPrefix)
	iter := as.db.NewIterator(r, nil)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()

		acc, err := account.FromBytes(value)
		if err != nil {
			logger.Panic("unable to decode account", "error", err)
		}

		var addr crypto.Address
		copy(addr[:], key[1:])

		stopped := consumer(addr, acc.Clone())
		if stopped {
			return
		}
	}
	iter.Release()
}

// This function takes ownership of the account pointer.
// It is important that the caller should not modify the account data and
// keep it immutable.
func (as *accountStore) updateAccount(batch *leveldb.Batch, addr crypto.Address, acc *account.Account) {
	data, err := acc.Bytes()
	if err != nil {
		logger.Panic("unable to encode account", "error", err)
	}
	if !as.hasAccount(addr) {
		as.total++
	}
	as.addrLruCache.Add(addr, acc)

	batch.Put(accountKey(addr), data)
}
