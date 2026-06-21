package state

import (
	"errors"
	"sync"
	"time"

	"github.com/pactus-project/pactus/committee"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/state/param"
	"github.com/pactus-project/pactus/store"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/types/protocol"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util/testsuite"
	"go.uber.org/mock/gomock"
)

var (
	errFakeNotFound = errors.New("not found")
	errFakeInvalid  = errors.New("invalid")
)

type FakeState struct {
	*MockState
	*testsuite.TestSuite

	lk             sync.RWMutex
	Committee      committee.Committee
	LastHeight     types.Height
	LastTime       time.Time
	Score          float64
	StateParams    *param.Params
	GenDoc         *genesis.Genesis
	ErrCommit      error
	ErrValidator   error
	TestBlocks     map[types.Height]*block.Block
	TestAccounts   map[crypto.Address]*account.Account
	TestValidators map[crypto.Address]*validator.Validator
}

func NewFakeState(ts *testsuite.TestSuite, committee committee.Committee) *FakeState {
	mock := NewMockState(ts.MockController())

	genDoc := genesis.MainnetGenesis()
	genTime := genDoc.GenesisTime()
	stateParams := param.FromGenesis(genesis.MainnetGenesis())
	stateParams.BlockVersion = protocol.ProtocolVersionLatest
	testBlocks := make(map[types.Height]*block.Block)
	testAccounts := make(map[crypto.Address]*account.Account)
	testValidators := make(map[crypto.Address]*validator.Validator)

	fake := &FakeState{
		MockState:      mock,
		TestSuite:      ts,
		Committee:      committee,
		LastHeight:     0,
		Score:          0.987,
		LastTime:       genTime,
		GenDoc:         genDoc,
		StateParams:    stateParams,
		TestBlocks:     testBlocks,
		TestAccounts:   testAccounts,
		TestValidators: testValidators,
	}

	mock.EXPECT().LastBlockHeight().DoAndReturn(func() types.Height {
		fake.lk.RLock()
		defer fake.lk.RUnlock()

		return fake.LastHeight
	}).AnyTimes()

	mock.EXPECT().LastBlockHash().DoAndReturn(func() hash.Hash {
		fake.lk.RLock()
		defer fake.lk.RUnlock()

		if fake.LastHeight == 0 {
			return hash.UndefHash
		}

		return fake.TestBlocks[fake.LastHeight].Hash()
	}).AnyTimes()

	mock.EXPECT().Genesis().DoAndReturn(func() *genesis.Genesis {
		return fake.GenDoc
	}).AnyTimes()

	mock.EXPECT().LastBlockTime().DoAndReturn(func() time.Time {
		fake.lk.RLock()
		defer fake.lk.RUnlock()

		return fake.LastTime
	}).AnyTimes()

	mock.EXPECT().Params().DoAndReturn(func() *param.Params {
		return fake.StateParams
	}).AnyTimes()

	mock.EXPECT().LastCertificate().DoAndReturn(func() *certificate.Certificate {
		return ts.GenerateTestCertificate(fake.LastHeight)
	}).AnyTimes()

	mock.EXPECT().IsProposer(gomock.Any(), gomock.Any()).DoAndReturn(
		func(addr crypto.Address, round types.Round) bool {
			return fake.IsProposer(addr, round)
		},
	).AnyTimes()

	mock.EXPECT().Proposer(gomock.Any()).DoAndReturn(
		func(round types.Round) *validator.Validator {
			return fake.Proposer(round)
		},
	).AnyTimes()

	mock.EXPECT().AvailabilityScore(gomock.Any()).DoAndReturn(
		func(int32) float64 {
			return fake.Score
		},
	).AnyTimes()

	mock.EXPECT().ProposeBlock(gomock.Any(), gomock.Any()).DoAndReturn(
		func(valKey *bls.ValidatorKey, _ crypto.Address) (*block.Block, error) {
			blk, _ := ts.GenerateTestBlock(fake.LastHeight+1,
				testsuite.BlockWithProposer(valKey.Address()))

			return blk, nil
		},
	).AnyTimes()

	mock.EXPECT().CommitBlock(gomock.Any(), gomock.Any()).DoAndReturn(
		func(blk *block.Block, cert *certificate.Certificate) error {
			fake.lk.Lock()
			defer fake.lk.Unlock()

			if cert.Height() == fake.LastHeight+1 {
				fake.TestBlocks[blk.Height()] = blk
				fake.LastHeight++
				fake.LastTime = fake.LastTime.Add(fake.StateParams.BlockInterval())

				return nil
			}

			return errFakeInvalid
		},
	).AnyTimes()

	mock.EXPECT().BlockHash(gomock.Any()).DoAndReturn(
		func(height types.Height) hash.Hash {
			blk, ok := fake.TestBlocks[height]
			if ok {
				return blk.Hash()
			}

			return hash.Hash{}
		},
	).AnyTimes()

	mock.EXPECT().ValidateBlock(gomock.Any(), gomock.Any()).DoAndReturn(
		func(*block.Block, types.Round) error {
			return fake.ErrValidator
		},
	).AnyTimes()

	mock.EXPECT().CommittedBlock(gomock.Any()).DoAndReturn(
		func(height types.Height) (*store.CommittedBlock, error) {
			blk, ok := fake.TestBlocks[height]
			if ok {
				data, _ := blk.Bytes()

				return &store.CommittedBlock{
					Data:      data,
					Height:    height,
					BlockHash: blk.Hash(),
				}, nil
			}

			return nil, errFakeNotFound
		},
	).AnyTimes()

	mock.EXPECT().CommitteeValidators().DoAndReturn(
		func() []*validator.Validator {
			return fake.Committee.Validators()
		},
	).AnyTimes()

	mock.EXPECT().IsInCommittee(gomock.Any()).DoAndReturn(
		func(addr crypto.Address) bool {
			return fake.Committee.Contains(addr)
		},
	).AnyTimes()

	mock.EXPECT().BlockHeight(gomock.Any()).DoAndReturn(
		func(h hash.Hash) types.Height {
			for height, blk := range fake.TestBlocks {
				if blk.Hash() == h {
					return height
				}
			}

			return 0
		},
	).AnyTimes()

	mock.EXPECT().CommittedTx(gomock.Any()).DoAndReturn(
		func(txID tx.ID) (*store.CommittedTx, error) {
			for height, blk := range fake.TestBlocks {
				for _, trx := range blk.Transactions() {
					if trx.ID() == txID {
						data, _ := trx.Bytes()

						return &store.CommittedTx{
							TxID:      txID,
							Height:    height,
							BlockTime: blk.Header().UnixTime(),
							Data:      data,
						}, nil
					}
				}
			}

			return nil, errFakeNotFound
		},
	).AnyTimes()

	mock.EXPECT().AccountByAddress(gomock.Any()).DoAndReturn(
		func(addr crypto.Address) (*account.Account, error) {
			acc, ok := fake.TestAccounts[addr]
			if ok {
				return acc, nil
			}

			return nil, errFakeNotFound
		},
	).AnyTimes()

	mock.EXPECT().ValidatorByAddress(gomock.Any()).DoAndReturn(
		func(addr crypto.Address) (*validator.Validator, error) {
			val, ok := fake.TestValidators[addr]
			if ok {
				return val, nil
			}

			return nil, errFakeNotFound
		},
	).AnyTimes()

	mock.EXPECT().ValidatorByNumber(gomock.Any()).DoAndReturn(
		func(num int32) (*validator.Validator, error) {
			for _, val := range fake.TestValidators {
				if val.Number() == num {
					return val, nil
				}
			}

			return nil, errFakeNotFound
		},
	).AnyTimes()

	mock.EXPECT().ValidatorAddresses().DoAndReturn(
		func() []crypto.Address {
			addrs := make([]crypto.Address, 0, len(fake.TestValidators))
			for _, val := range fake.TestValidators {
				addrs = append(addrs, val.Address())
			}

			return addrs
		},
	).AnyTimes()

	mock.EXPECT().CommitteeInfo().DoAndReturn(
		func() *CommitteeInfo {
			return &CommitteeInfo{
				Validators:       fake.Committee.Validators(),
				ProtocolVersions: fake.Committee.ProtocolVersions(),
				CommitteePower:   fake.Committee.Power(),
			}
		},
	).AnyTimes()

	mock.EXPECT().ChainInfo().DoAndReturn(
		func() *ChainInfo {
			lastBlockHash := hash.UndefHash
			if fake.LastHeight > 0 {
				lastBlockHash = fake.TestBlocks[fake.LastHeight].Hash()
			}

			return &ChainInfo{
				LastBlockHeight: fake.LastHeight,
				LastBlockHash:   lastBlockHash,
				LastBlockTime:   fake.LastTime,
				TotalPower:      fake.Committee.Power(),
				CommitteePower:  fake.Committee.Power(),
				CommitteeSize:   fake.Committee.Size(),
				TotalAccounts:   int32(len(fake.TestAccounts)),
				TotalValidators: int32(len(fake.TestValidators)),
				AverageScore:    fake.Score,
			}
		},
	).AnyTimes()

	return fake
}

func (f *FakeState) ProposerIndex(round types.Round) int {
	len := f.Committee.Size()
	i := int(f.LastHeight)%len + int(round)%len

	return i % len
}

func (f *FakeState) Proposer(round types.Round) *validator.Validator {
	return f.Committee.Validators()[f.ProposerIndex(round)]
}

func (f *FakeState) CommitTestBlocks(count int) {
	for i := 0; i < count; i++ {
		blk, cert := f.GenerateTestBlock(f.LastHeight + 1)
		_ = f.CommitBlock(blk, cert)
	}
}

func (f *FakeState) AddTestBlock(height types.Height, opts ...testsuite.BlockMakerOption) *block.Block {
	blk, _ := f.GenerateTestBlock(height, opts...)
	f.TestBlocks[height] = blk

	return blk
}

func (f *FakeState) AddTestAccount(opts ...testsuite.AccountMakerOption) (crypto.Address, *account.Account) {
	acc, addr := f.GenerateTestAccount(opts...)
	f.TestAccounts[addr] = acc

	return addr, acc
}

func (f *FakeState) AddTestValidator(opts ...testsuite.ValidatorMakerOption) *validator.Validator {
	val := f.GenerateTestValidator(opts...)
	f.TestValidators[val.Address()] = val

	return val
}
