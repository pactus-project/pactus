package validator

import (
	"github.com/sasha-s/go-deadlock"
	"github.com/zarbchain/zarb-go/crypto"
)

type ValidatorPool struct {
	lk deadlock.RWMutex

	validators map[crypto.Address]*Validator
}

func NewValidatorPool() *ValidatorPool {
	return &ValidatorPool{
		validators: make(map[crypto.Address]*Validator),
	}
}

func (vp *ValidatorPool) AddValidator(val *Validator, height int) {
	vp.validators[val.Address()] = val
}

func (vp *ValidatorPool) ValidatorInfo(addr crypto.Address) *Validator {
	return vp.validators[addr]
}

func (vp *ValidatorPool) TotalStake() int64 {
	total := int64(0)
	for _, vi := range vp.validators {
		total += vi.Stake()
	}

	return total
}
