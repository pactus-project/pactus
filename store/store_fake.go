package store

import (
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/sortition"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/types/protocol"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util/testsuite"
	gomock "go.uber.org/mock/gomock"
)

type FakeStore struct {
	*MockStore

	FakeBlocks       map[types.Height]*block.Block
	FakeCertificates map[types.Height]*certificate.Certificate
	FakeAccounts     map[crypto.Address]*account.Account
	FakeValidators   map[crypto.Address]*validator.Validator
}

func NewFakeStore(ts *testsuite.TestSuite) *FakeStore {
	fake := &FakeStore{
		MockStore:        NewMockStore(ts.MockController()),
		FakeBlocks:       make(map[types.Height]*block.Block),
		FakeCertificates: make(map[types.Height]*certificate.Certificate),
		FakeAccounts:     make(map[crypto.Address]*account.Account),
		FakeValidators:   make(map[crypto.Address]*validator.Validator),
	}

	fake.EXPECT().LastCertificate().DoAndReturn(
		func() *certificate.Certificate {
			if len(fake.FakeCertificates) == 0 {
				return nil
			}

			lastHeight := types.Height(0)
			for h := range fake.FakeCertificates {
				if h > lastHeight {
					lastHeight = h
				}
			}

			return fake.FakeCertificates[lastHeight]
		},
	).AnyTimes()

	fake.EXPECT().Block(gomock.Any()).DoAndReturn(
		func(height types.Height) (*CommittedBlock, error) {
			blk, ok := fake.FakeBlocks[height]
			if !ok {
				return nil, ErrNotFound
			}

			data, _ := blk.Bytes()

			return &CommittedBlock{
				Data:      data,
				Height:    height,
				BlockHash: blk.Hash(),
			}, nil
		},
	).AnyTimes()

	fake.EXPECT().BlockHash(gomock.Any()).DoAndReturn(
		func(height types.Height) hash.Hash {
			blk, ok := fake.FakeBlocks[height]
			if !ok {
				return hash.UndefHash
			}

			return blk.Hash()
		},
	).AnyTimes()

	fake.EXPECT().BlockHeight(gomock.Any()).DoAndReturn(
		func(h hash.Hash) types.Height {
			for height, blk := range fake.FakeBlocks {
				if blk.Hash() == h {
					return height
				}
			}

			return 0
		},
	).AnyTimes()

	fake.EXPECT().Transaction(gomock.Any()).DoAndReturn(
		func(txID tx.ID) (*CommittedTx, error) {
			for _, blk := range fake.FakeBlocks {
				for _, trx := range blk.Transactions() {
					if txID == trx.ID() {

						data, _ := trx.Bytes()

						return &CommittedTx{
							Data:      data,
							TxID:      txID,
							Height:    blk.Height(),
							BlockTime: blk.Header().UnixTime(),
						}, nil
					}
				}
			}

			return nil, ErrNotFound
		},
	).AnyTimes()

	fake.EXPECT().HasAccount(gomock.Any()).DoAndReturn(
		func(addr crypto.Address) bool {
			_, ok := fake.FakeAccounts[addr]

			return ok
		},
	).AnyTimes()

	fake.EXPECT().Account(gomock.Any()).DoAndReturn(
		func(addr crypto.Address) (*account.Account, error) {
			acc, ok := fake.FakeAccounts[addr]
			if !ok {
				return nil, ErrNotFound
			}

			return acc, nil
		},
	).AnyTimes()

	fake.EXPECT().TotalAccounts().DoAndReturn(
		func() int32 {
			return int32(len(fake.FakeAccounts))
		},
	).AnyTimes()

	fake.EXPECT().HasValidator(gomock.Any()).DoAndReturn(
		func(addr crypto.Address) bool {
			_, ok := fake.FakeValidators[addr]

			return ok
		},
	).AnyTimes()

	fake.EXPECT().ValidatorAddresses().DoAndReturn(
		func() []crypto.Address {
			addrs := make([]crypto.Address, 0, len(fake.FakeValidators))
			for addr := range fake.FakeValidators {
				addrs = append(addrs, addr)
			}

			return addrs
		},
	).AnyTimes()

	fake.EXPECT().Validator(gomock.Any()).DoAndReturn(
		func(addr crypto.Address) (*validator.Validator, error) {
			val, ok := fake.FakeValidators[addr]
			if !ok {
				return nil, ErrNotFound
			}

			return val, nil
		},
	).AnyTimes()

	fake.EXPECT().ValidatorByNumber(gomock.Any()).DoAndReturn(
		func(num int32) (*validator.Validator, error) {
			for _, val := range fake.FakeValidators {
				if val.Number() == num {
					return val, nil
				}
			}

			return nil, ErrNotFound
		},
	).AnyTimes()

	fake.EXPECT().IterateValidators(gomock.Any()).DoAndReturn(
		func(consumer func(*validator.Validator) bool) {
			for _, val := range fake.FakeValidators {
				if consumer(val) {
					return
				}
			}
		},
	).AnyTimes()

	fake.EXPECT().IterateAccounts(gomock.Any()).DoAndReturn(
		func(consumer func(crypto.Address, *account.Account) bool) {
			for addr, acc := range fake.FakeAccounts {
				if consumer(addr, acc) {
					return
				}
			}
		},
	).AnyTimes()

	fake.EXPECT().TotalValidators().DoAndReturn(
		func() int32 {
			return int32(len(fake.FakeValidators))
		},
	).AnyTimes()
	fake.EXPECT().ActiveValidators().DoAndReturn(
		func() int32 {
			return int32(len(fake.FakeValidators))
		},
	).AnyTimes()

	fake.EXPECT().IsPruned().Return(false).AnyTimes()
	fake.EXPECT().PruningHeight().Return(types.Height(0)).AnyTimes()

	fake.EXPECT().UpdateAccount(gomock.Any(), gomock.Any()).DoAndReturn(
		func(addr crypto.Address, acc *account.Account) {
			fake.FakeAccounts[addr] = acc
		},
	).AnyTimes()

	fake.EXPECT().UpdateValidator(gomock.Any()).DoAndReturn(
		func(val *validator.Validator) {
			fake.FakeValidators[val.Address()] = val
		},
	).AnyTimes()

	fake.EXPECT().SaveBlock(gomock.Any(), gomock.Any()).DoAndReturn(
		func(blk *block.Block, cert *certificate.Certificate) {
			fake.FakeBlocks[blk.Height()] = blk
			fake.FakeCertificates[cert.Height()] = cert
		},
	).AnyTimes()

	fake.EXPECT().UpdateValidatorProtocolVersion(gomock.Any(), gomock.Any()).DoAndReturn(
		func(addr crypto.Address, ver protocol.Version) {
			val, ok := fake.FakeValidators[addr]
			if ok {
				val.UpdateProtocolVersion(ver)
			}
		},
	).AnyTimes()

	fake.EXPECT().SortitionSeed(gomock.Any()).DoAndReturn(
		func(height types.Height) *sortition.VerifiableSeed {
			blk, ok := fake.FakeBlocks[height]
			if ok {
				seed := blk.Header().SortitionSeed()

				return &seed
			}

			return &sortition.UndefVerifiableSeed
		},
	).AnyTimes()

	fake.EXPECT().WriteBatch().Return(nil).AnyTimes()
	fake.EXPECT().Close().Return().AnyTimes()

	return fake
}
