package store

import (
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type accountStore struct {
	db    *leveldb.DB
	total int32
}

func accountKey(addr crypto.Address) []byte { return append(accountPrefix, addr.Bytes()...) }

func newAccountStore(db *leveldb.DB) *accountStore {
	as := &accountStore{
		db: db,
	}
	// TODO: better way to get total accout number?
	total := int32(0)
	as.iterateAccounts(func(_ crypto.Address, _ *account.Account) bool {
		total++
		return false
	})
	as.total = total

	return as
}

func (as *accountStore) hasAccount(addr crypto.Address) bool {
	has, err := as.db.Has(accountKey(addr), nil)
	if err != nil {
		return false
	}
	return has
}

func (as *accountStore) account(addr crypto.Address) (*account.Account, error) {
	bs, err := tryGet(as.db, accountKey(addr))
	if err != nil {
		return nil, err
	}

	acc, err := account.FromBytes(bs)
	if err != nil {
		return nil, err
	}

	return acc, nil
}

func (as *accountStore) iterateAccounts(consumer func(crypto.Address, *account.Account) (stop bool)) {
	r := util.BytesPrefix(accountPrefix)
	iter := as.db.NewIterator(r, nil)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()

		var addr crypto.Address
		copy(addr[:], key[1:])

		acc, err := account.FromBytes(value)
		if err != nil {
			logger.Panic("unable to decode account: %v", err)
		}

		stopped := consumer(addr, acc)
		if stopped {
			return
		}
	}
	iter.Release()
}

func (as *accountStore) updateAccount(batch *leveldb.Batch, addr crypto.Address, acc *account.Account) {
	data, err := acc.Bytes()
	if err != nil {
		logger.Panic("unable to encode account: %v", err)
	}
	if !as.hasAccount(addr) {
		as.total++
	}
	batch.Put(accountKey(addr), data)
}
