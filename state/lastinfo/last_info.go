package lastinfo

import (
	"fmt"
	"sync"
	"time"

	"github.com/pactus-project/pactus/committee"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/sortition"
	"github.com/pactus-project/pactus/store"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util/logger"
)

type LastInfo struct {
	lk sync.RWMutex // TODO: this lock looks unnecessary

	lastSortitionSeed sortition.VerifiableSeed
	lastBlockHash     hash.Hash
	lastCert          *certificate.BlockCertificate
	lastBlockTime     time.Time
	lastValidators    []*validator.Validator
}

func NewLastInfo() *LastInfo {
	return &LastInfo{}
}

func (li *LastInfo) SortitionSeed() sortition.VerifiableSeed {
	li.lk.RLock()
	defer li.lk.RUnlock()

	return li.lastSortitionSeed
}

func (li *LastInfo) BlockHeight() uint32 {
	li.lk.RLock()
	defer li.lk.RUnlock()

	if li.lastCert == nil {
		return 0
	}

	return li.lastCert.Height()
}

func (li *LastInfo) BlockHash() hash.Hash {
	li.lk.RLock()
	defer li.lk.RUnlock()

	return li.lastBlockHash
}

func (li *LastInfo) Certificate() *certificate.BlockCertificate {
	li.lk.RLock()
	defer li.lk.RUnlock()

	return li.lastCert
}

func (li *LastInfo) BlockTime() time.Time {
	li.lk.RLock()
	defer li.lk.RUnlock()

	return li.lastBlockTime
}

func (li *LastInfo) Validators() []*validator.Validator {
	li.lk.RLock()
	defer li.lk.RUnlock()

	return li.lastValidators
}

func (li *LastInfo) UpdateSortitionSeed(lastSortitionSeed sortition.VerifiableSeed) {
	li.lk.Lock()
	defer li.lk.Unlock()

	li.lastSortitionSeed = lastSortitionSeed
}

func (li *LastInfo) UpdateBlockHash(lastBlockHash hash.Hash) {
	li.lk.Lock()
	defer li.lk.Unlock()

	li.lastBlockHash = lastBlockHash
}

func (li *LastInfo) UpdateCertificate(lastCertificate *certificate.BlockCertificate) {
	li.lk.Lock()
	defer li.lk.Unlock()

	li.lastCert = lastCertificate
}

func (li *LastInfo) UpdateBlockTime(lastBlockTime time.Time) {
	li.lk.Lock()
	defer li.lk.Unlock()

	li.lastBlockTime = lastBlockTime
}

func (li *LastInfo) UpdateValidators(vals []*validator.Validator) {
	li.lk.Lock()
	defer li.lk.Unlock()

	li.lastValidators = vals
}

func (li *LastInfo) RestoreLastInfo(store store.Store, committeeSize int) (committee.Committee, error) {
	lastCert := store.LastCertificate()
	lastHeight := lastCert.Height()
	logger.Debug("try to restore last state info", "height", lastHeight)
	sb, err := store.Block(lastHeight)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve block %v: %w", lastHeight, err)
	}

	lastBlock, err := sb.ToBlock()
	if err != nil {
		return nil, err
	}

	li.lastCert = lastCert
	li.lastSortitionSeed = lastBlock.Header().SortitionSeed()
	li.lastBlockHash = lastBlock.Hash()
	li.lastBlockTime = lastBlock.Header().Time()

	cmt, err := li.restoreCommittee(store, lastBlock, committeeSize)
	if err != nil {
		return nil, err
	}

	return cmt, nil
}

func (li *LastInfo) restoreCommittee(store store.Store, lastBlock *block.Block,
	committeeSize int,
) (committee.Committee, error) {
	joinedVals := make([]*validator.Validator, 0)
	for _, trx := range lastBlock.Transactions() {
		// If there is any sortition transaction in the last block,
		// we should update the last committee.
		if trx.IsSortitionTx() {
			pld := trx.Payload().(*payload.SortitionPayload)
			val, err := store.Validator(pld.Validator)
			if err != nil {
				return nil, fmt.Errorf("unable to retrieve validator %s: %w", pld.Validator, err)
			}
			joinedVals = append(joinedVals, val)
		}
	}

	proposerIndex := -1
	curCommitteeSize := len(li.lastCert.Committers())
	vals := make([]*validator.Validator, len(li.lastCert.Committers()))
	for i, num := range li.lastCert.Committers() {
		val, err := store.ValidatorByNumber(num)
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve committee member %v: %w", num, err)
		}
		if lastBlock.Header().ProposerAddress() == val.Address() {
			proposerIndex = i
		}
		vals[i] = val
	}
	li.lastValidators = vals

	// First, we restore the previous committee; then, we update it to get the latest committee.
	proposerIndex = (proposerIndex + curCommitteeSize -
		(int(li.lastCert.Round()) % curCommitteeSize)) % curCommitteeSize
	cmt, err := committee.NewCommittee(vals, committeeSize, vals[proposerIndex].Address())
	if err != nil {
		return nil, fmt.Errorf("unable to create last committee: %w", err)
	}
	cmt.Update(li.lastCert.Round(), joinedVals)

	return cmt, nil
}
