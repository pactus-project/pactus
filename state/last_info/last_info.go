package last_info

import (
	"fmt"
	"time"

	"github.com/fxamacker/cbor/v2"
	"github.com/sasha-s/go-deadlock"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/committee"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/sortition"
	"github.com/zarbchain/zarb-go/store"
	"github.com/zarbchain/zarb-go/tx/payload"
	"github.com/zarbchain/zarb-go/validator"
)

type lastInfoData struct {
	LastHeight      int
	LastCertificate *block.Certificate
}

type LastInfo struct {
	lk deadlock.RWMutex

	store             store.Store
	lastBlockHeight   int
	lastSortitionSeed sortition.Seed
	lastBlockHash     crypto.Hash
	lastCertificate   *block.Certificate
	lastBlockTime     time.Time
}

func NewLastInfo(store store.Store) *LastInfo {
	return &LastInfo{store: store}
}

func (li *LastInfo) SortitionSeed() sortition.Seed {
	li.lk.RLock()
	defer li.lk.RUnlock()

	return li.lastSortitionSeed
}

func (li *LastInfo) BlockHeight() int {
	li.lk.RLock()
	defer li.lk.RUnlock()

	return li.lastBlockHeight
}

func (li *LastInfo) BlockHash() crypto.Hash {
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

func (li *LastInfo) SetSortitionSeed(lastSortitionSeed sortition.Seed) {
	li.lk.Lock()
	defer li.lk.Unlock()

	li.lastSortitionSeed = lastSortitionSeed
}

func (li *LastInfo) SetBlockHeight(lastBlockHeight int) {
	li.lk.Lock()
	defer li.lk.Unlock()

	li.lastBlockHeight = lastBlockHeight
}

func (li *LastInfo) SetBlockHash(lastBlockHash crypto.Hash) {
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

func (li *LastInfo) SaveLastInfo() {
	lid := lastInfoData{
		LastHeight:      li.lastBlockHeight,
		LastCertificate: li.lastCertificate,
	}

	bs, _ := cbor.Marshal(&lid)
	li.store.SaveLastInfo(bs)
}

func (li *LastInfo) RestoreLastInfo() (*committee.Committee, error) {
	bs := li.store.RestoreLastInfo()
	lid := new(lastInfoData)
	err := cbor.Unmarshal(bs, lid)
	if err != nil {
		return nil, err
	}
	logger.Debug("Try to restore last state info", "height", lid.LastHeight)

	b, err := li.store.Block(lid.LastHeight)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve block %v: %v", lid.LastHeight, err)
	}

	joinedVals := make([]*validator.Validator, 0)
	for _, id := range b.TxIDs().IDs() {
		ctx, err := li.store.Transaction(id)
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve transaction %s: %v", id, err)
		}
		if ctx.Tx.IsSortitionTx() {
			pld := ctx.Tx.Payload().(*payload.SortitionPayload)
			val, err := li.store.Validator(pld.Address)
			if err != nil {
				return nil, fmt.Errorf("unable to retrieve validator %s: %v", pld.Address, err)
			}
			joinedVals = append(joinedVals, val)
		}
	}

	li.lastBlockHeight = lid.LastHeight
	li.lastCertificate = lid.LastCertificate
	li.lastSortitionSeed = b.Header().SortitionSeed()
	li.lastBlockHash = b.Hash()
	li.lastBlockTime = b.Header().Time()

	vals := make([]*validator.Validator, len(b.LastCertificate().Committers()))
	for i, num := range b.LastCertificate().Committers() {
		val, err := li.store.ValidatorByNumber(num)
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve committee member %v: %v", num, err)
		}
		vals[i] = val
	}
	committee, err := committee.NewCommittee(vals, len(b.LastCertificate().Committers()), b.Header().ProposerAddress())
	if err != nil {
		return nil, fmt.Errorf("unable to create last committee: %v", err)
	}

	err = committee.Update(0, joinedVals)
	if err != nil {
		return nil, fmt.Errorf("unable to update last committee: %v", err)
	}

	return committee, nil
}
