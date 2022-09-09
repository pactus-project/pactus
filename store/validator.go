package store

import (
	"fmt"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type validatorStore struct {
	db     *leveldb.DB
	valMap map[int32]*validator.Validator
	total  int32
}

func validatorKey(addr crypto.Address) []byte { return append(validatorPrefix, addr.Bytes()...) }

func newValidatorStore(db *leveldb.DB) *validatorStore {
	vs := &validatorStore{
		db: db,
	}

	total := int32(0)
	valMap := make(map[int32]*validator.Validator)
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

	val, err := validator.FromBytes(data)
	if err != nil {
		return nil, err
	}

	return val, nil
}

func (vs *validatorStore) validatorByNumber(num int32) (*validator.Validator, error) {
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

		val, err := validator.FromBytes(value)
		if err != nil {
			logger.Panic("unable to decode validator: %v", err)
		}

		stopped := consumer(val)
		if stopped {
			return
		}
	}
	iter.Release()
}

func (vs *validatorStore) updateValidator(batch *leveldb.Batch, val *validator.Validator) {
	data, err := val.Bytes()
	if err != nil {
		logger.Panic("unable to encode validator: %v", err)
	}
	if !vs.hasValidator(val.Address()) {
		vs.total++
	}
	vs.valMap[val.Number()] = val

	batch.Put(validatorKey(val.Address()), data)
}
