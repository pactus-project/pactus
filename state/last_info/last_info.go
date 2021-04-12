package last_info

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/sasha-s/go-deadlock"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/committee"
	"github.com/zarbchain/zarb-go/crypto"
	merkle "github.com/zarbchain/zarb-go/libs/merkle"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/sortition"
	"github.com/zarbchain/zarb-go/store"
	"github.com/zarbchain/zarb-go/tx/payload"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/validator"
)

type lastInfoData struct {
	LastHeight      int
	LastCertificate *block.Certificate
}

type LastInfo struct {
	lk deadlock.RWMutex

	path string // temproray

	lastBlockHeight   int
	lastSortitionSeed sortition.Seed
	lastBlockHash     crypto.Hash
	lastReceiptsHash  crypto.Hash
	lastCertificate   *block.Certificate
	lastBlockTime     time.Time
}

func NewLastInfo(path string) *LastInfo {
	return &LastInfo{path: path}
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

func (li *LastInfo) ReceiptsHash() crypto.Hash {
	li.lk.RLock()
	defer li.lk.RUnlock()

	return li.lastReceiptsHash
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

func (li *LastInfo) SetReceiptsHash(lastReceiptsHash crypto.Hash) {
	li.lk.Lock()
	defer li.lk.Unlock()

	li.lastReceiptsHash = lastReceiptsHash
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

func (li *LastInfo) SaveLastInfo() error {
	path := li.path + "/last_info.json"
	lid := lastInfoData{
		LastHeight:      li.lastBlockHeight,
		LastCertificate: li.lastCertificate,
	}

	bs, _ := json.Marshal(&lid)
	if err := util.WriteFile(path, bs); err != nil {
		return fmt.Errorf("Unable to write last sate info: %v", err)
	}
	return nil
}

func (li *LastInfo) RestoreLastInfo(store store.StoreReader) (*committee.Committee, error) {
	path := li.path + "/last_info.json"
	if !util.PathExists(path) {
		return nil, fmt.Errorf("Unable to load %v", path)
	}
	bs, err := util.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lid := new(lastInfoData)
	err = json.Unmarshal(bs, lid)
	if err != nil {
		return nil, err
	}
	logger.Debug("Try to restore last state info", "height", lid.LastHeight)

	b, err := store.Block(lid.LastHeight)
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve block %v: %v", lid.LastHeight, err)
	}

	joinedVals := make([]*validator.Validator, 0)
	receiptsHashes := make([]crypto.Hash, len(b.TxIDs().IDs()))
	for i, id := range b.TxIDs().IDs() {
		ctx, err := store.Transaction(id)
		if err != nil {
			return nil, fmt.Errorf("Unable to retrieve transaction %s: %v", id, err)
		}
		receiptsHashes[i] = ctx.Receipt.Hash()

		if ctx.Tx.IsSortitionTx() {
			pld := ctx.Tx.Payload().(*payload.SortitionPayload)
			val, err := store.Validator(pld.Address)
			if err != nil {
				return nil, fmt.Errorf("Unable to retrieve validator %s: %v", pld.Address, err)
			}
			joinedVals = append(joinedVals, val)
		}
	}
	receiptsMerkle := merkle.NewTreeFromHashes(receiptsHashes)

	li.lastBlockHeight = lid.LastHeight
	li.lastCertificate = lid.LastCertificate
	li.lastSortitionSeed = b.Header().SortitionSeed()
	li.lastBlockHash = b.Hash()
	li.lastReceiptsHash = receiptsMerkle.Root()
	li.lastBlockTime = b.Header().Time()

	vals := make([]*validator.Validator, len(b.LastCertificate().Committers()))
	for i, num := range b.LastCertificate().Committers() {
		val, err := store.ValidatorByNumber(num)
		if err != nil {
			return nil, fmt.Errorf("Unable to retrieve committee member %v: %v", num, err)
		}
		vals[i] = val
	}
	committee, err := committee.NewCommittee(vals, len(b.LastCertificate().Committers()), b.Header().ProposerAddress())
	if err != nil {
		return nil, fmt.Errorf("Unable to create last committee: %v", err)
	}

	err = committee.Update(0, joinedVals)
	if err != nil {
		return nil, fmt.Errorf("Unable to update last committee: %v", err)
	}

	return committee, nil
}
