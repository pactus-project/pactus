package sortition

import (
	"github.com/sasha-s/go-deadlock"
	"github.com/zarbchain/zarb-go/crypto"
)

type Sortition struct {
	lk deadlock.RWMutex

	vrf *VRF
}

func NewSortition() *Sortition {
	return &Sortition{
		vrf: NewVRF(),
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

func (s *Sortition) EvaluateSortition(seed Seed, signer crypto.Signer, threshold int64) (bool, Proof) {
	s.lk.RLock()
	defer s.lk.RUnlock()

	index, proof := s.vrf.Evaluate(seed, signer)
	if index > threshold {
		return false, proof
	}

	return true, proof
}

func (s *Sortition) VerifyProof(seed Seed, proof Proof, public crypto.PublicKey, threshold int64) bool {
	s.lk.RLock()
	defer s.lk.RUnlock()

	index, result := s.vrf.Verify(seed, public, proof)
	if !result {
		return false
	}
	if index > threshold {
		return false
	}

	return true
}
