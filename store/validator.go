package store

import (
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/validator"
)

type validatorStore struct {
	db     *leveldb.DB
	valMap map[int]*validator.Validator
	total  int
}

func validatorKey(addr crypto.Address) []byte { return append(validatorPrefix, addr.RawBytes()...) }

func newValidatorStore(db *leveldb.DB) *validatorStore {
	vs := &validatorStore{
		db: db,
	}

	total := 0
	valMap := make(map[int]*validator.Validator)
	vs.iterateValidators(func(val *validator.Validator) bool {
		valMap[val.Number()] = val
		total++
		return false
	})

	vs.total = total
	vs.valMap = valMap

	return vs
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

func (vs *validatorStore) validatorByNumber(num int) (*validator.Validator, error) {
	val, ok := vs.valMap[num]
	if ok {
		return val, nil
	}

	return nil, fmt.Errorf("not found")
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

func (vs *validatorStore) updateValidator(batch *leveldb.Batch, val *validator.Validator) error {
	data, err := val.Encode()
	if err != nil {
		return err
	}
	if !vs.hasValidator(val.Address()) {
		vs.total++
	}
	vs.valMap[val.Number()] = val

	batch.Put(validatorKey(val.Address()), data)

	return nil
}
