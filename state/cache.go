package state

import (
	"bytes"

	"github.com/sasha-s/go-deadlock"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/libs/orderedmap"
	"github.com/zarbchain/zarb-go/store"
	"github.com/zarbchain/zarb-go/validator"
)

type Cache struct {
	lk deadlock.Mutex

	name       string
	store      *store.Store
	valChanges *orderedmap.OrderedMap
	accChanges *orderedmap.OrderedMap
}

type validatorInfo struct {
	validator *validator.Validator
}

type accountInfo struct {
	account  *account.Account
	storages *orderedmap.OrderedMap
}

type CacheOption func(*Cache)

func lessFn(l, r interface{}) bool {
	return bytes.Compare(l.(crypto.Address).RawBytes(), r.(crypto.Address).RawBytes()) < 0
}

func newCache(store *store.Store) *Cache {
	ch := &Cache{
		store:      store,
		valChanges: orderedmap.NewMap(lessFn),
		accChanges: orderedmap.NewMap(lessFn),
	}
	return ch
}

func (c *Cache) reset() {
	c.lk.Lock()
	defer c.lk.Unlock()

	c.accChanges = orderedmap.NewMap(lessFn)
	c.valChanges = orderedmap.NewMap(lessFn)
}

func (c *Cache) commit(set *validator.ValidatorSet) error {
	c.lk.Lock()
	defer c.lk.Unlock()

	c.accChanges.Iter(func(key, value interface{}) (more bool) {
		i := value.(*accountInfo)

		c.store.UpdateAccount(i.account)

		// if i.storages != nil {
		// 	i.storages.Iter(func(k, v interface{}) (more bool) {
		// 		if err := c.store.setStorage(i.account.Address(), k.(Word256), v.(Word256)); err != nil {
		// 			panic(err)
		// 		}
		// 		return true
		// 	})
		// }

		return true
	})

	c.valChanges.Iter(func(key, value interface{}) (more bool) {
		i := value.(*validatorInfo)
		c.store.UpdateValidator(i.validator)
		return true
	})

	/// reset cache
	c.accChanges = orderedmap.NewMap(lessFn)
	c.valChanges = orderedmap.NewMap(lessFn)

	return nil
}

func (c *Cache) HasAccount(addr crypto.Address) bool {
	c.lk.Lock()
	defer c.lk.Unlock()

	_, ok := c.accChanges.GetOk(addr)
	if ok {
		return true
	}

	return c.store.HasAccount(addr)
}

func (c *Cache) Account(addr crypto.Address) *account.Account {
	c.lk.Lock()
	defer c.lk.Unlock()

	i, ok := c.accChanges.GetOk(addr)
	if ok {
		return i.(*accountInfo).account
	}

	return c.store.RetrieveAccount(addr)
}

func (c *Cache) UpdateAccount(acc *account.Account) {
	c.lk.Lock()
	defer c.lk.Unlock()

	addr := acc.Address()
	i, ok := c.accChanges.GetOk(addr)
	if ok {
		i.(*accountInfo).account = acc
	} else {
		c.accChanges.Set(addr, &accountInfo{account: acc})
	}
}

func (c *Cache) HasValidator(addr crypto.Address) bool {
	c.lk.Lock()
	defer c.lk.Unlock()

	_, ok := c.valChanges.GetOk(addr)
	if ok {
		return true
	}

	return c.store.HasValidator(addr)
}

func (c *Cache) Validator(addr crypto.Address) *validator.Validator {
	c.lk.Lock()
	defer c.lk.Unlock()

	i, ok := c.valChanges.GetOk(addr)
	if ok {
		return i.(*validatorInfo).validator
	}

	return c.store.RetrieveValidator(addr)
}

func (c *Cache) UpdateValidator(val *validator.Validator) {
	c.lk.Lock()
	defer c.lk.Unlock()

	addr := val.Address()
	i, ok := c.valChanges.GetOk(addr)
	if ok {
		i.(*validatorInfo).validator = val
	} else {
		c.valChanges.Set(addr, &validatorInfo{validator: val})
	}
}

// func (c *Cache) AddToSet(val *validator.Validator) error {
// 	c.lk.Lock()
// 	defer c.lk.Unlock()

// 	addr := val.Address()
// 	_, ok := c.valChanges.GetOk(addr)
// 	if ok {
// 		return errValidatorChanged
// 	}

// 	c.valChanges.Set(addr, &validatorInfo{addToSet, val})
// 	return nil
// }

// func (c *Cache) GetStorage(addr crypto.Address, key Word256) (Word256, error) {
// 	c.lk.Lock()
// 	defer c.lk.Unlock()

// 	i, ok := c.accChanges.GetOk(addr)
// 	if ok {
// 		if i.(*accountInfo).storages != nil {
// 			s, ok := i.(*accountInfo).storages.GetOk(key)
// 			if ok {
// 				return s.(Word256), nil
// 			}
// 		}
// 	}

// 	return c.store.GetStorage(addr, key)
// }

// func (c *Cache) SetStorage(addr crypto.Address, key, value Word256) error {
// 	c.lk.Lock()
// 	defer c.lk.Unlock()

// 	i, ok := c.accChanges.GetOk(addr)
// 	if !ok {
// 		acc, _ := c.store.GetAccount(addr)
// 		if acc == nil {
// 			acc, _ = account.NewAccount(addr)
// 		}

// 		i = &accountInfo{account: acc}
// 		c.accChanges.Set(addr, i)
// 	}

// 	if i.(*accountInfo).storages == nil {
// 		i.(*accountInfo).storages = orderedmap.NewMap(lessFn2)
// 	}

// 	i.(*accountInfo).storages.Set(key, value)
// 	return nil
// }
