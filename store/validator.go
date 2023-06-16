package store

import (
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type validatorStore struct {
	db         *leveldb.DB
	numberMap  map[int32]*validator.Validator
	addressMap map[crypto.Address]*validator.Validator
	total      int32
}

func validatorKey(addr crypto.Address) []byte { return append(validatorPrefix, addr.Bytes()...) }

func newValidatorStore(db *leveldb.DB) *validatorStore {
	total := int32(0)
	numberMap := make(map[int32]*validator.Validator)
	addressMap := make(map[crypto.Address]*validator.Validator)
	r := util.BytesPrefix(validatorPrefix)
	iter := db.NewIterator(r, nil)
	for iter.Next() {
		// key := iter.Key()
		value := iter.Value()

		val, err := validator.FromBytes(value)
		if err != nil {
			logger.Panic("unable to decode validator", "err", err)
		}

		numberMap[val.Number()] = val
		addressMap[val.Address()] = val
		total++
	}
	iter.Release()

	return &validatorStore{
		db:         db,
		total:      total,
		numberMap:  numberMap,
		addressMap: addressMap,
	}
}

func (vs *validatorStore) hasValidator(addr crypto.Address) bool {
	_, ok := vs.addressMap[addr]
	return ok
}

func (vs *validatorStore) validator(addr crypto.Address) (*validator.Validator, error) {
	val, ok := vs.addressMap[addr]
	if ok {
		return val.Clone(), nil
	}

	return nil, ErrNotFound
}

func (vs *validatorStore) validatorByNumber(num int32) (*validator.Validator, error) {
	val, ok := vs.numberMap[num]
	if ok {
		return val.Clone(), nil
	}

	return nil, ErrNotFound
}

func (vs *validatorStore) iterateValidators(consumer func(*validator.Validator) (stop bool)) {
	for _, val := range vs.addressMap {
		stopped := consumer(val.Clone())
		if stopped {
			return
		}
	}
}

// This function takes ownership of the validator pointer.
// It is important that the caller should not modify the validator data and
// keep it immutable.
func (vs *validatorStore) updateValidator(batch *leveldb.Batch, val *validator.Validator) {
	data, err := val.Bytes()
	if err != nil {
		logger.Panic("unable to encode validator", "err", err)
	}
	if !vs.hasValidator(val.Address()) {
		vs.total++
	}
	vs.numberMap[val.Number()] = val
	vs.addressMap[val.Address()] = val

	batch.Put(validatorKey(val.Address()), data)
}
