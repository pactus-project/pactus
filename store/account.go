package store

import (
	"cmp"
	"errors"
	"slices"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/pactus-project/gopkg/logger"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/account"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type accountStore struct {
	db       *leveldb.DB
	accCache *lru.Cache[crypto.Address, *account.Account]
	total    int32

	// anchors lists the accounts that hold an anchor, sorted by account number.
	// It is not consensus data; it only serves anchor listings.
	anchors []anchorEntry
}

type anchorEntry struct {
	number int32
	addr   crypto.Address
}

func accountKey(addr crypto.Address) []byte { return append(accountPrefix, addr.Bytes()...) }

func newAccountStore(db *leveldb.DB, cacheSize int) *accountStore {
	total := int32(0)
	addrLruCache, err := lru.New[crypto.Address, *account.Account](cacheSize)
	if err != nil {
		logger.Panic("unable to create new instance of lru cache", "error", err)
	}

	// A plain account record has no anchor suffix, so only longer records need decoding.
	plainSize := account.NewAccount(0).SerializeSize()
	anchors := make([]anchorEntry, 0)

	r := util.BytesPrefix(accountPrefix)
	iter := db.NewIterator(r, nil)
	for iter.Next() {
		total++

		if len(iter.Value()) <= plainSize {
			continue
		}
		acc, err := account.FromBytes(iter.Value())
		if err != nil {
			logger.Panic("unable to decode account", "error", err)
		}
		if acc.HasAnchor() {
			var addr crypto.Address
			copy(addr[:], iter.Key()[1:])
			anchors = append(anchors, anchorEntry{number: acc.Number(), addr: addr})
		}
	}
	iter.Release()

	slices.SortFunc(anchors, func(left, right anchorEntry) int {
		return cmp.Compare(left.number, right.number)
	})

	return &accountStore{
		db:       db,
		total:    total,
		accCache: addrLruCache,
		anchors:  anchors,
	}
}

func (as *accountStore) hasAccount(addr crypto.Address) bool {
	ok := as.accCache.Contains(addr)
	if !ok {
		ok = tryHas(as.db, accountKey(addr))
	}

	return ok
}

func (as *accountStore) account(addr crypto.Address) (*account.Account, error) {
	acc, ok := as.accCache.Get(addr)
	if ok {
		return acc.Clone(), nil
	}

	rawData, err := tryGet(as.db, accountKey(addr))
	if err != nil {
		if errors.Is(err, leveldb.ErrNotFound) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	acc, err = account.FromBytes(rawData)
	if err != nil {
		return nil, err
	}

	as.accCache.Add(addr, acc)

	return acc.Clone(), nil
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

		stopped := consumer(addr, acc)
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
	as.accCache.Add(addr, acc)
	as.updateAnchorIndex(addr, acc)

	batch.Put(accountKey(addr), data)
}

func (as *accountStore) updateAnchorIndex(addr crypto.Address, acc *account.Account) {
	idx, found := slices.BinarySearchFunc(as.anchors, acc.Number(),
		func(entry anchorEntry, number int32) int {
			return cmp.Compare(entry.number, number)
		})

	switch {
	case acc.HasAnchor() && !found:
		as.anchors = slices.Insert(as.anchors, idx, anchorEntry{number: acc.Number(), addr: addr})
	case !acc.HasAnchor() && found:
		as.anchors = slices.Delete(as.anchors, idx, idx+1)
	}
}

// anchorAddresses returns a page of anchor holders, ordered by account number,
// and the total number of anchor holders.
func (as *accountStore) anchorAddresses(skip, count uint32) ([]crypto.Address, uint32) {
	total := uint32(len(as.anchors))
	if skip >= total || count == 0 {
		return []crypto.Address{}, total
	}

	end := min(uint64(skip)+uint64(count), uint64(total))
	addrs := make([]crypto.Address, 0, end-uint64(skip))
	for _, entry := range as.anchors[skip:end] {
		addrs = append(addrs, entry.addr)
	}

	return addrs, total
}
