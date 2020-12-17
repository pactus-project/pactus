package sortition

import (
	"github.com/sasha-s/go-deadlock"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/validator"
)

type Sortition struct {
	lk deadlock.RWMutex

	signer crypto.Signer
	vrf    *VRF
}

func NewSortition(signer crypto.Signer) *Sortition {
	return &Sortition{
		signer: signer,
		vrf:    NewVRF(signer),
	}
}

func (s *Sortition) SetTotalStake(totalStake int64) {
	s.lk.RLock()
	defer s.lk.RUnlock()

	s.vrf.SetMax(totalStake)
}

// AddToTotalStake adds new stakes to total stake. stake can be negative
func (s *Sortition) AddToTotalStake(stake int64) {
	s.lk.RLock()
	defer s.lk.RUnlock()

	s.vrf.AddToMax(stake)
}

func (s *Sortition) TotalStake() int64 {
	s.lk.RLock()
	defer s.lk.RUnlock()

	return s.vrf.Max()
}

func (s *Sortition) Evaluate(hash crypto.Hash, val *validator.Validator) *tx.Tx {
	s.lk.RLock()
	defer s.lk.RUnlock()

	index, proof := s.vrf.Evaluate(hash)
	if index >= val.Stake() {
		return nil
	}

	trx := tx.NewSortitionTx(hash, val.Sequence()+1, val.Address(), index, proof, "", nil, nil)
	s.signer.SignMsg(trx)
	return trx
}

func (s *Sortition) VerifySortition(blockHash crypto.Hash, index int64, proof []byte, val *validator.Validator) bool {
	s.lk.RLock()
	defer s.lk.RUnlock()

	index, result := s.vrf.Verify(blockHash, val.PublicKey(), proof)
	if !result {
		return false
	}
	if index >= val.Stake() {
		return false
	}

	return true
}
