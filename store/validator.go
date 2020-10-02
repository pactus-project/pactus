package store

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/validator"
)

type validatorStore struct {
	db *leveldb.DB
}

var (
	validatorPrefix = []byte{0x01}
)

func validatorKey(addr crypto.Address) []byte { return append(validatorPrefix, addr.RawBytes()...) }

func newValidatorStore(path string) (*validatorStore, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &validatorStore{
		db: db,
	}, nil
}

func (vs *validatorStore) HasValidator(addr crypto.Address) bool {
	has, err := vs.db.Has(validatorKey(addr), nil)
	if err != nil {
		panic(err)
	}
	return has
}

func (vs *validatorStore) RetrieveValidator(addr crypto.Address) *validator.Validator {
	bs, _ := vs.db.Get(validatorKey(addr), nil)
	if bs == nil {
		return nil
	}

	val := new(validator.Validator)
	if err := val.Decode(bs); err != nil {
		panic(err)
	}

	return val
}

func (vs *validatorStore) ValidatorCount() int {
	count := 0
	vs.IterateValidators(func(validator *validator.Validator) (stop bool) {
		count++
		return false
	})
	return count
}

func (vs *validatorStore) IterateValidators(consumer func(*validator.Validator) (stop bool)) {
	r := util.BytesPrefix(validatorPrefix)
	iter := vs.db.NewIterator(r, nil)
	for iter.Next() {
		// key := iter.Key()
		value := iter.Value()

		val := new(validator.Validator)
		if err := val.Decode(value); err != nil {
			panic(err)
		}

		stopped := consumer(val)
		if stopped {
			return
		}

	}
	iter.Release()
}

func (vs *validatorStore) UpdateValidator(val *validator.Validator) {
	bs, err := val.Encode()
	if err != nil {
		panic(err)
	}

	vs.db.Put(validatorKey(val.Address()), bs, nil)
}
