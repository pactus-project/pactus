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

func (vs *validatorStore) hasValidator(addr crypto.Address) bool {
	has, err := vs.db.Has(validatorKey(addr), nil)
	if err != nil {
		return false
	}
	return has
}

func (vs *validatorStore) validator(addr crypto.Address) (*validator.Validator, error) {
	data, err := tryGet(vs.db, validatorKey(addr))
	if err != nil {
		return nil, err
	}

	val := new(validator.Validator)
	if err := val.Decode(data); err != nil {
		return nil, err
	}

	return val, nil
}

func (vs *validatorStore) len() int {
	len := 0
	vs.iterateValidators(func(validator *validator.Validator) (stop bool) {
		len++
		return false
	})
	return len
}

func (vs *validatorStore) iterateValidators(consumer func(*validator.Validator) (stop bool)) {
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

func (vs *validatorStore) updateValidator(val *validator.Validator) error {
	data, err := val.Encode()
	if err != nil {
		return err
	}

	return tryPut(vs.db, validatorKey(val.Address()), data)
}
