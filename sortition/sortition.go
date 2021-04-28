package sortition

import (
	"github.com/sasha-s/go-deadlock"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/libs/linkedmap"
)

type param struct {
	seed  Seed
	stake int64
}

type Sortition struct {
	lk deadlock.RWMutex

	params *linkedmap.LinkedMap
	vrf    *VRF
}

func NewSortition() *Sortition {
	return &Sortition{
		vrf:    NewVRF(),
		params: linkedmap.NewLinkedMap(7), // Sortitions are valid for 7 height
	}
}

func (s *Sortition) SetParams(blockHash crypto.Hash, seed Seed, poolStake int64) {
	s.lk.Lock()
	defer s.lk.Unlock()

	p := &param{
		seed:  seed,
		stake: poolStake,
	}
	s.params.PushBack(blockHash, p)
}

func (s *Sortition) GetParams(blockHash crypto.Hash) (seed Seed, poolStake int64) {
	s.lk.RLock()
	defer s.lk.RUnlock()

	p := s.getParam(blockHash)
	if p == nil {
		return
	}

	return p.seed, p.stake
}

func (s *Sortition) EvaluateSortition(blockHash crypto.Hash, signer crypto.Signer, threshold int64) (bool, Proof) {
	s.lk.RLock()
	defer s.lk.RUnlock()

	p := s.getParam(blockHash)
	if p == nil {
		return false, Proof{}
	}

	index, proof := s.vrf.Evaluate(p.seed, signer, p.stake)
	if index < threshold {
		return true, proof
	}

	return false, Proof{}
}

func (s *Sortition) VerifyProof(blockHash crypto.Hash, proof Proof, public crypto.PublicKey, threshold int64) bool {
	s.lk.RLock()
	defer s.lk.RUnlock()

	p := s.getParam(blockHash)
	if p == nil {
		return false
	}

	index, result := s.vrf.Verify(p.seed, public, proof, p.stake)
	if !result {
		return false
	}
	return index < threshold
}

func (s *Sortition) getParam(hash crypto.Hash) *param {
	p, ok := s.params.Get(hash)
	if !ok {
		return nil
	}

	return p.(*param)
}
