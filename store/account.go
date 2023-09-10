package store

import (
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type accountStore struct {
	db         *leveldb.DB
	numberMap  map[int32]*account.Account
	addressMap map[crypto.Address]*account.Account
	total      int32
}

func accountKey(addr crypto.Address) []byte { return append(accountPrefix, addr.Bytes()...) }

func newAccountStore(db *leveldb.DB) *accountStore {
	total := int32(0)
	numberMap := make(map[int32]*account.Account)
	addressMap := make(map[crypto.Address]*account.Account)
	r := util.BytesPrefix(accountPrefix)
	iter := db.NewIterator(r, nil)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()

		acc, err := account.FromBytes(value)
		if err != nil {
			logger.Panic("unable to decode account", "error", err)
		}

		var addr crypto.Address
		copy(addr[:], key[1:])

		numberMap[acc.Number()] = acc
		addressMap[addr] = acc
		total++
	}
	iter.Release()

	return &accountStore{
		db:         db,
		total:      total,
		numberMap:  numberMap,
		addressMap: addressMap,
	}
}

func (as *accountStore) hasAccount(addr crypto.Address) bool {
	_, ok := as.addressMap[addr]
	return ok
}

func (as *accountStore) account(addr crypto.Address) (*account.Account, error) {
	acc, ok := as.addressMap[addr]
	if ok {
		return acc.Clone(), nil
	}

	return nil, ErrNotFound
}

func (as *accountStore) accountByNumber(number int32) (*account.Account, error) {
	acc, ok := as.numberMap[number]
	if ok {
		return acc.Clone(), nil
	}

	return nil, ErrNotFound
}

func (as *accountStore) iterateAccounts(consumer func(crypto.Address, *account.Account) (stop bool)) {
	for addr, acc := range as.addressMap {
		stopped := consumer(addr, acc.Clone())
		if stopped {
			return
		}
	}
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
	as.numberMap[acc.Number()] = acc
	as.addressMap[addr] = acc

	batch.Put(accountKey(addr), data)
}
