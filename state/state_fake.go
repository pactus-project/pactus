package state

import (
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

type FakeState struct {
	*MockState
	*testsuite.TestSuite

	lk sync.RWMutex

	FakeCommittee *committee.FakeCommittee
	FakeStore     *store.FakeStore
	FakeHeight    types.Height
	FakeScore     float64
	StateParams   *param.Params
	GenDoc        *genesis.Genesis
	ErrCommit     error
	ErrValidator  error
}

func NewFakeState(ts *testsuite.TestSuite) *FakeState {
	mock := NewMockState(ts.MockController())

	genDoc := genesis.MainnetGenesis()
	stateParams := param.FromGenesis(genesis.MainnetGenesis())
	stateParams.BlockVersion = protocol.ProtocolVersionLatest
	fakeCommittee := committee.NewFakeCommittee(ts)
	fakeStore := store.NewFakeStore(ts)

	fake := &FakeState{
		MockState:     mock,
		TestSuite:     ts,
		FakeCommittee: fakeCommittee,
		FakeHeight:    0,
		FakeScore:     0.987,
		GenDoc:        genDoc,
		StateParams:   stateParams,
		FakeStore:     fakeStore,
	}

	mock.EXPECT().LastBlockHeight().DoAndReturn(func() types.Height {
		fake.lk.RLock()
		defer fake.lk.RUnlock()

		return fake.FakeHeight
	}).AnyTimes()

	mock.EXPECT().LastBlockHash().DoAndReturn(func() hash.Hash {
		fake.lk.RLock()
		defer fake.lk.RUnlock()

		cBlk, _ := fake.FakeStore.Block(fake.FakeHeight)
		if cBlk != nil {
			return cBlk.BlockHash
		}

		return hash.UndefHash
	}).AnyTimes()

	mock.EXPECT().Genesis().DoAndReturn(func() *genesis.Genesis {
		return fake.GenDoc
	}).AnyTimes()

	mock.EXPECT().LastBlockTime().DoAndReturn(func() time.Time {
		fake.lk.RLock()
		defer fake.lk.RUnlock()

		cBlk, _ := fake.FakeStore.Block(fake.FakeHeight)
		if cBlk != nil {
			blk, _ := cBlk.ToBlock()

			return blk.Header().Time()
		}

		return fake.GenDoc.GenesisTime()
	}).AnyTimes()

	mock.EXPECT().Params().DoAndReturn(func() *param.Params {
		return fake.StateParams
	}).AnyTimes()

	mock.EXPECT().LastCertificate().DoAndReturn(func() *certificate.Certificate {
		return ts.GenerateTestCertificate(fake.FakeHeight)
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
			return fake.FakeScore
		},
	).AnyTimes()

	mock.EXPECT().ProposeBlock(gomock.Any(), gomock.Any()).DoAndReturn(
		func(valKey *bls.ValidatorKey, _ crypto.Address) (*block.Block, error) {
			blk, _ := ts.GenerateTestBlock(fake.FakeHeight+1,
				testsuite.BlockWithProposer(valKey.Address()))

			return blk, nil
		},
	).AnyTimes()

	mock.EXPECT().CommitBlock(gomock.Any(), gomock.Any()).DoAndReturn(
		func(blk *block.Block, cert *certificate.Certificate) error {
			fake.lk.Lock()
			defer fake.lk.Unlock()

			fake.FakeStore.SaveBlock(blk, cert)

			if cert.Height() > fake.FakeHeight {
				fake.FakeHeight = cert.Height()
			}

			return fake.ErrCommit
		},
	).AnyTimes()

	mock.EXPECT().BlockHash(gomock.Any()).DoAndReturn(
		func(height types.Height) hash.Hash {
			fake.lk.Lock()
			defer fake.lk.Unlock()

			cBlk, _ := fake.FakeStore.Block(height)
			if cBlk != nil {
				return cBlk.BlockHash
			}

			return hash.UndefHash
		},
	).AnyTimes()

	mock.EXPECT().ValidateBlock(gomock.Any(), gomock.Any()).DoAndReturn(
		func(*block.Block, types.Round) error {
			return fake.ErrValidator
		},
	).AnyTimes()

	mock.EXPECT().CommittedBlock(gomock.Any()).DoAndReturn(
		func(height types.Height) (*store.CommittedBlock, error) {
			fake.lk.Lock()
			defer fake.lk.Unlock()

			return fake.FakeStore.Block(height)
		},
	).AnyTimes()

	mock.EXPECT().CommitteeValidators().DoAndReturn(
		func() []*validator.Validator {
			return fake.FakeCommittee.Validators()
		},
	).AnyTimes()

	mock.EXPECT().IsInCommittee(gomock.Any()).DoAndReturn(
		func(addr crypto.Address) bool {
			return fake.FakeCommittee.Contains(addr)
		},
	).AnyTimes()

	mock.EXPECT().BlockHeight(gomock.Any()).DoAndReturn(
		func(h hash.Hash) types.Height {
			return fake.FakeStore.BlockHeight(h)
		},
	).AnyTimes()

	mock.EXPECT().CommittedTx(gomock.Any()).DoAndReturn(
		func(txID tx.ID) (*store.CommittedTx, error) {
			return fake.FakeStore.Transaction(txID)
		},
	).AnyTimes()

	mock.EXPECT().AccountByAddress(gomock.Any()).DoAndReturn(
		func(addr crypto.Address) (*account.Account, error) {
			return fake.FakeStore.Account(addr)
		},
	).AnyTimes()

	mock.EXPECT().ValidatorByAddress(gomock.Any()).DoAndReturn(
		func(addr crypto.Address) (*validator.Validator, error) {
			return fake.FakeStore.Validator(addr)
		},
	).AnyTimes()

	mock.EXPECT().ValidatorByNumber(gomock.Any()).DoAndReturn(
		func(num int32) (*validator.Validator, error) {
			return fake.FakeStore.ValidatorByNumber(num)
		},
	).AnyTimes()

	mock.EXPECT().ValidatorAddresses().DoAndReturn(
		func() []crypto.Address {
			return fake.FakeStore.ValidatorAddresses()
		},
	).AnyTimes()

	mock.EXPECT().CommitteeInfo().DoAndReturn(
		func() *CommitteeInfo {
			return &CommitteeInfo{
				Validators:       fake.FakeCommittee.Validators(),
				ProtocolVersions: fake.FakeCommittee.ProtocolVersions(),
				CommitteePower:   fake.FakeCommittee.Power(),
				TotalPower:       fake.FakeCommittee.Power(),
			}
		},
	).AnyTimes()

	mock.EXPECT().ChainInfo().DoAndReturn(
		func() *ChainInfo {
			fake.lk.Lock()
			defer fake.lk.Unlock()

			lastBlockHash := hash.UndefHash
			lastBlockTime := genDoc.GenesisTime()

			cBlk, _ := fake.FakeStore.Block(fake.FakeHeight)
			if cBlk != nil {
				blk, _ := cBlk.ToBlock()

				lastBlockHash = blk.Hash()
				lastBlockTime = blk.Header().Time()
			}

			return &ChainInfo{
				LastBlockHeight:  fake.FakeHeight,
				LastBlockHash:    lastBlockHash,
				LastBlockTime:    lastBlockTime,
				TotalPower:       fake.FakeCommittee.Power(),
				CommitteePower:   fake.FakeCommittee.Power(),
				CommitteeSize:    fake.FakeCommittee.Size(),
				TotalAccounts:    fake.FakeStore.TotalAccounts(),
				TotalValidators:  fake.FakeStore.TotalValidators(),
				ActiveValidators: fake.FakeStore.ActiveValidators(),
				AverageScore:     fake.FakeScore,
			}
		},
	).AnyTimes()

	return fake
}

func (f *FakeState) ProposerIndex(round types.Round) int {
	len := f.FakeCommittee.Size()
	i := int(f.FakeHeight)%len + int(round)%len

	return i % len
}

func (f *FakeState) Proposer(round types.Round) *validator.Validator {
	return f.FakeCommittee.Validators()[f.ProposerIndex(round)]
}

func (f *FakeState) CommitTestBlocks(count int) {
	blkTime := f.GenDoc.GenesisTime()
	for i := 0; i < count; i++ {
		blk, cert := f.GenerateTestBlock(f.FakeHeight+1, testsuite.BlockWithTime(blkTime))
		_ = f.CommitBlock(blk, cert)

		blkTime = blkTime.Add(10 * time.Second)
	}
}

func (f *FakeState) AddTestBlock(blk *block.Block, cert *certificate.Certificate) {
	f.FakeStore.SaveBlock(blk, cert)
}

func (f *FakeState) AddTestAccount(addr crypto.Address, acc *account.Account) {
	f.FakeStore.UpdateAccount(addr, acc)
}

func (f *FakeState) AddTestValidator(val *validator.Validator) {
	f.FakeStore.UpdateValidator(val)
}
