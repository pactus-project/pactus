package validator

import (
	"github.com/sasha-s/go-deadlock"
	"github.com/zarbchain/zarb-go/crypto"
)

type ValidatorInfo struct {
	Validator

	BoundingHeight int
	ProposedBlocks int
	SignedBlocks   int
	MissedBlocks   int
}

type ValidatorPool struct {
	lk deadlock.RWMutex

	validators map[crypto.Address]*ValidatorInfo
}

func NewValidatorPool() *ValidatorPool {
	return &ValidatorPool{
		validators: make(map[crypto.Address]*ValidatorInfo),
	}
}

func (vp *ValidatorPool) AddValidator(val *Validator, height int) {
	vi := &ValidatorInfo{
		Validator:      *val,
		BoundingHeight: height,
	}

	vp.validators[val.Address()] = vi
}

func (vp *ValidatorPool) ValidatorInfo(addr crypto.Address) *ValidatorInfo {
	return vp.validators[addr]
}

func (vp *ValidatorPool) TotalStake() int64 {
	total := int64(0)
	for _, vi := range vp.validators {
		total += vi.Stake()
	}

	return total
}
