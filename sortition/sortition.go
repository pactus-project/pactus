package sortition

import (
	"sync"

	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/libs/linkedmap"
)

const changeCoefficient = 100000000

type blockParams struct {
	seed  VerifiableSeed
	stake int64
}

type Sortition struct {
	lk sync.RWMutex

	params *linkedmap.LinkedMap
}

func NewSortition() *Sortition {
	return &Sortition{
		params: linkedmap.NewLinkedMap(30), // Sortitions are valid for 30 height
	}
}

func (s *Sortition) SetParams(blockHash hash.Hash, seed VerifiableSeed, stake int64) {
	s.lk.Lock()
	defer s.lk.Unlock()

	p := &blockParams{
		seed:  seed,
		stake: stake,
	}
	s.params.PushBack(blockHash, p)
}

func (s *Sortition) GetParams(blockHash hash.Hash) (seed VerifiableSeed, stake int64) {
	s.lk.RLock()
	defer s.lk.RUnlock()

	p := s.getParam(blockHash)
	if p == nil {
		return
	}

	return p.seed, p.stake
}

func (s *Sortition) VerifyProof(blockHash hash.Hash, proof Proof, public crypto.PublicKey, stake int64) bool {
	s.lk.RLock()
	defer s.lk.RUnlock()

	if proof.Coin <= 0 || proof.Coin > stakeToCoin(stake) {
		return false
	}

	p := s.getParam(blockHash)
	if p == nil {
		return false
	}

	return verifyProof(p.seed, public, proof, stakeToCoin(p.stake))
}

func (s *Sortition) getParam(hash hash.Hash) *blockParams {
	p, ok := s.params.Get(hash)
	if !ok {
		return nil
	}

	return p.(*blockParams)
}

func EvaluateSortition(blockSeed VerifiableSeed, signer crypto.Signer, validatorStake, totalStake int64) (bool, Proof) {
	return evaluate(blockSeed, signer, stakeToCoin(validatorStake), stakeToCoin(totalStake))
}

func stakeToCoin(stake int64) int {
	return int(stake / changeCoefficient)
}
