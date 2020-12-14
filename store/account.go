package store

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/crypto"
)

type accountStore struct {
	db    *leveldb.DB
	total int
}

var (
	accountPrefix = []byte{0x01}
)

func accountKey(addr crypto.Address) []byte { return append(accountPrefix, addr.RawBytes()...) }

func newAccountStore(path string) (*accountStore, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	as := &accountStore{
		db: db,
	}
	as.total = as.countAccounts()

	return as, nil
}

func (as *accountStore) close() error {
	return as.db.Close()
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

	acc := new(account.Account)
	if err := acc.Decode(bs); err != nil {
		return nil, err
	}

	return acc, nil
}

func (as *accountStore) iterateAccounts(consumer func(*account.Account) (stop bool)) {
	r := util.BytesPrefix(accountPrefix)
	iter := as.db.NewIterator(r, nil)
	for iter.Next() {
		// key := iter.Key()
		value := iter.Value()

		acc := new(account.Account)
		if err := acc.Decode(value); err != nil {
			panic(err)
		}

		stopped := consumer(acc)
		if stopped {
			return
		}

	}
	iter.Release()
}

func (as *accountStore) updateAccount(acc *account.Account) error {
	data, err := acc.Encode()
	if err != nil {
		panic(err)
	}
	if !as.hasAccount(acc.Address()) {
		as.total++
	}

	return tryPut(as.db, accountKey(acc.Address()), data)
}

func (as *accountStore) countAccounts() int {
	count := 0
	as.iterateAccounts(func(acc *account.Account) bool {
		count++
		return false
	})
	return count
}
