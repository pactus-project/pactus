package store

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/crypto"
)

type accountStore struct {
	db *leveldb.DB
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
	return &accountStore{
		db: db,
	}, nil
}

func (as *accountStore) HasAccount(addr crypto.Address) bool {
	has, err := as.db.Has(accountKey(addr), nil)
	if err != nil {
		panic(err)
	}
	return has
}

func (as *accountStore) RetrieveAccount(addr crypto.Address) *account.Account {
	bs, _ := as.db.Get(accountKey(addr), nil)
	if bs == nil {
		return nil
	}

	acc := new(account.Account)
	if err := acc.Decode(bs); err != nil {
		panic(err)
	}

	return acc
}

func (as *accountStore) AccountCount() int {
	count := 0
	as.IterateAccounts(func(acc *account.Account) (stop bool) {
		count++
		return false
	})
	return count
}

func (as *accountStore) IterateAccounts(consumer func(*account.Account) (stop bool)) {
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

func (as *accountStore) UpdateAccount(acc *account.Account) {
	bs, err := acc.Encode()
	if err != nil {
		panic(err)
	}

	as.db.Put(accountKey(acc.Address()), bs, nil)
}
