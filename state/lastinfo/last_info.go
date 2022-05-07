package lastinfo

import (
	"fmt"
	"sync"
	"time"

	"github.com/zarbchain/zarb-go/committee"
	"github.com/zarbchain/zarb-go/sortition"
	"github.com/zarbchain/zarb-go/store"
	"github.com/zarbchain/zarb-go/types/block"
	"github.com/zarbchain/zarb-go/types/crypto/hash"
	"github.com/zarbchain/zarb-go/types/tx/payload"
	"github.com/zarbchain/zarb-go/types/validator"
	"github.com/zarbchain/zarb-go/util/logger"
)

type LastInfo struct {
	lk sync.RWMutex

	store             store.Store
	lastBlockHeight   int32
	lastSortitionSeed sortition.VerifiableSeed
	lastBlockHash     hash.Hash
	lastCertificate   *block.Certificate
	lastBlockTime     time.Time
}

func NewLastInfo(store store.Store) *LastInfo {
	return &LastInfo{store: store}
}

func (li *LastInfo) SortitionSeed() sortition.VerifiableSeed {
	li.lk.RLock()
	defer li.lk.RUnlock()

	return li.lastSortitionSeed
}

func (li *LastInfo) BlockHeight() int32 {
	li.lk.RLock()
	defer li.lk.RUnlock()

	return li.lastBlockHeight
}

func (li *LastInfo) BlockHash() hash.Hash {
	li.lk.RLock()
	defer li.lk.RUnlock()

	return li.lastBlockHash
}

func (li *LastInfo) Certificate() *block.Certificate {
	li.lk.RLock()
	defer li.lk.RUnlock()

	return li.lastCertificate
}

func (li *LastInfo) BlockTime() time.Time {
	li.lk.RLock()
	defer li.lk.RUnlock()

	return li.lastBlockTime
}

func (li *LastInfo) SetSortitionSeed(lastSortitionSeed sortition.VerifiableSeed) {
	li.lk.Lock()
	defer li.lk.Unlock()

	li.lastSortitionSeed = lastSortitionSeed
}

func (li *LastInfo) SetBlockHeight(lastBlockHeight int32) {
	li.lk.Lock()
	defer li.lk.Unlock()

	li.lastBlockHeight = lastBlockHeight
}

func (li *LastInfo) SetBlockHash(lastBlockHash hash.Hash) {
	li.lk.Lock()
	defer li.lk.Unlock()

	li.lastBlockHash = lastBlockHash
}

func (li *LastInfo) SetCertificate(lastCertificate *block.Certificate) {
	li.lk.Lock()
	defer li.lk.Unlock()

	li.lastCertificate = lastCertificate
}

func (li *LastInfo) SetBlockTime(lastBlockTime time.Time) {
	li.lk.Lock()
	defer li.lk.Unlock()

	li.lastBlockTime = lastBlockTime
}

func (li *LastInfo) RestoreLastInfo(committeeSize int) (committee.Committee, error) {
	height, cert := li.store.LastCertificate()

	logger.Debug("try to restore last state info", "height", height)

	h := li.store.BlockHash(height)
	bi, err := li.store.Block(h)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve block %v: %v", height, err)
	}

	b, _ := bi.ToFullBlock()

	li.lastBlockHeight = height
	li.lastCertificate = cert
	li.lastSortitionSeed = b.Header().SortitionSeed()
	li.lastBlockHash = b.Hash()
	li.lastBlockTime = b.Header().Time()

	cmt, err := li.restoreCommittee(b, committeeSize)
	if err != nil {
		return nil, err
	}

	return cmt, nil
}

func (li *LastInfo) restoreCommittee(b *block.Block, committeeSize int) (committee.Committee, error) {
	joinedVals := make([]*validator.Validator, 0)
	for _, trx := range b.Transactions() {
		// If there is any sortition transaction in last block,
		// we should update last committee
		if trx.IsSortitionTx() {
			pld := trx.Payload().(*payload.SortitionPayload)
			val, err := li.store.Validator(pld.Address)
			if err != nil {
				return nil, fmt.Errorf("unable to retrieve validator %s: %v", pld.Address, err)
			}
			joinedVals = append(joinedVals, val)
		}
	}

	proposerIndex := -1
	curCommitteeSize := len(li.lastCertificate.Committers())
	vals := make([]*validator.Validator, len(li.lastCertificate.Committers()))
	for i, num := range li.lastCertificate.Committers() {
		val, err := li.store.ValidatorByNumber(num)
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve committee member %v: %v", num, err)
		}
		if b.Header().ProposerAddress().EqualsTo(val.Address()) {
			proposerIndex = i
		}
		vals[i] = val
	}

	// First we restore previous committee, then we update it to get the latest committee.
	proposerIndex = (proposerIndex + curCommitteeSize -
		(int(li.lastCertificate.Round()) % curCommitteeSize)) % curCommitteeSize
	committee, err := committee.NewCommittee(vals, committeeSize, vals[proposerIndex].Address())
	if err != nil {
		return nil, fmt.Errorf("unable to create last committee: %v", err)
	}
	committee.Update(li.lastCertificate.Round(), joinedVals)

	return committee, nil
}
