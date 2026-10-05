package txpool

import (
	"bytes"
	"io"
	"strings"
	"testing"
	"time"
	"unsafe"

	"github.com/pactus-project/gopkg/pipeline"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/execution"
	"github.com/pactus-project/pactus/execution/executor"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/state/param"
	"github.com/pactus-project/pactus/store"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/protocol"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const anchorTestHeight = types.Height(20_000)

type anchorLab struct {
	*testsuite.TestSuite

	t         *testing.T
	pool      *txPool
	sbx       *sandbox.MockSandbox
	broadcast pipeline.Pipeline[message.Message]
	events    pipeline.Pipeline[any]
	accounts  map[crypto.Address]*account.Account
	committed map[tx.ID]struct{}
	banned    map[crypto.Address]bool
	indexed   map[crypto.Address]bool
	version   protocol.Version
	height    types.Height
}

func newAnchorLab(t *testing.T, version protocol.Version, cfg *Config) *anchorLab {
	t.Helper()

	lab := &anchorLab{
		TestSuite: testsuite.NewTestSuite(t),
		t:         t,
		accounts:  make(map[crypto.Address]*account.Account),
		committed: make(map[tx.ID]struct{}),
		banned:    make(map[crypto.Address]bool),
		indexed:   make(map[crypto.Address]bool),
		version:   version,
		height:    anchorTestHeight,
	}

	previous := executor.DefaultFactory
	t.Cleanup(func() {
		executor.DefaultFactory = previous
	})
	mockExe := executor.NewMockExecutor(lab.MockController())
	mockExe.EXPECT().Check(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockExe.EXPECT().Execute(gomock.Any()).Return().AnyTimes()
	executor.DefaultFactory = func(trx *tx.Tx, sbx sandbox.Sandbox) (executor.Executor, error) {
		if trx.Payload().Type() == payload.TypeAnchor {
			return executor.MakeExecutorImpl(trx, sbx)
		}

		return mockExe, nil
	}

	lab.broadcast = pipeline.New[message.Message](t.Context())
	lab.events = pipeline.New[any](t.Context())
	lab.sbx = lab.bindSandbox()

	fakeStore := store.NewFakeStore(lab.TestSuite)
	fakeStore.EXPECT().HasPublicKey(gomock.Any()).DoAndReturn(func(addr crypto.Address) bool {
		return lab.indexed[addr]
	}).AnyTimes()

	config := testDefaultConfig()
	if cfg != nil {
		config = cfg
	}
	poolInt := NewTxPool(t.Context(), config, fakeStore, lab.broadcast, lab.events)
	poolInt.SetNewSandboxAndRecheck(lab.sbx)
	lab.pool = poolInt.(*txPool)

	return lab
}

func (lab *anchorLab) bindSandbox() *sandbox.MockSandbox {
	sbx := sandbox.NewMockSandbox(lab.MockController())
	sbx.EXPECT().CurrentHeight().DoAndReturn(func() types.Height { return lab.height }).AnyTimes()
	sbx.EXPECT().CurrentUnixTime().Return(uint32(1_700_000_000)).AnyTimes()
	sbx.EXPECT().IsBanned(gomock.Any()).DoAndReturn(func(addr crypto.Address) bool {
		return lab.banned[addr]
	}).AnyTimes()
	sbx.EXPECT().Params().Return(param.FromGenesis(genesis.MainnetGenesis())).AnyTimes()
	sbx.EXPECT().BlockVersion().DoAndReturn(func() protocol.Version { return lab.version }).AnyTimes()
	sbx.EXPECT().Account(gomock.Any()).DoAndReturn(func(addr crypto.Address) *account.Account {
		return lab.accounts[addr]
	}).AnyTimes()
	sbx.EXPECT().RecentTransaction(gomock.Any()).DoAndReturn(func(id tx.ID) bool {
		_, ok := lab.committed[id]

		return ok
	}).AnyTimes()
	sbx.EXPECT().UpdateAccount(gomock.Any(), gomock.Any()).Return().AnyTimes()
	sbx.EXPECT().CommitTransaction(gomock.Any()).DoAndReturn(func(trx *tx.Tx) {
		lab.committed[trx.ID()] = struct{}{}
	}).AnyTimes()

	return sbx
}

func (lab *anchorLab) fund(balance amount.Amount) (crypto.Address, crypto.PrivateKey, *account.Account) {
	pub, prv := lab.RandBLSKeyPair()
	acc := account.NewAccount(int32(len(lab.accounts) + 1))
	acc.AddToBalance(balance)
	addr := pub.AccountAddress()
	lab.accounts[addr] = acc

	return addr, prv, acc
}

func (lab *anchorLab) sign(prv crypto.PrivateKey, trx *tx.Tx) *tx.Tx {
	lab.HelperSignTransaction(prv, trx)

	return trx
}

func (lab *anchorLab) setTx( //nolint:revive // the anchor fields stay one call
	addr crypto.Address, prv crypto.PrivateKey, deposit, fee amount.Amount,
	root []byte, uri string, anchorType uint8, lockTime types.Height, memo string,
) *tx.Tx {
	trx := tx.NewAnchorTx(lockTime, addr, payload.AnchorActionSet, root, uri, anchorType, deposit, fee)
	if memo != "" {
		trx = tx.NewAnchorTx(lockTime, addr, payload.AnchorActionSet, root, uri, anchorType, deposit, fee, tx.WithMemo(memo))
	}

	return lab.sign(prv, trx)
}

func (lab *anchorLab) deleteTx(addr crypto.Address, prv crypto.PrivateKey, fee amount.Amount) *tx.Tx {
	return lab.sign(prv, tx.NewAnchorTx(lab.height, addr, payload.AnchorActionDelete, nil, "", 0, 0, fee))
}

func anchorRoot(fill byte, length int) []byte {
	return bytes.Repeat([]byte{fill}, length)
}

func (lab *anchorLab) fee() amount.Amount {
	return lab.pool.config.fixedFee()
}

func (lab *anchorLab) requireQuiet(t *testing.T) {
	t.Helper()

	select {
	case <-lab.broadcast.UnsafeGetChannel():
		require.Fail(t, "transaction was broadcast")
	default:
	}
	select {
	case <-lab.events.UnsafeGetChannel():
		require.Fail(t, "transaction was published")
	default:
	}
}

func (lab *anchorLab) requirePublished(t *testing.T, txID tx.ID) {
	t.Helper()

	timer := time.NewTimer(time.Second)
	defer timer.Stop()

	gotBroadcast := false
	gotEvent := false
	for !gotBroadcast || !gotEvent {
		select {
		case <-timer.C:
			require.Fail(t, "publish timeout")
		case msg := <-lab.broadcast.UnsafeGetChannel():
			if msg.Type() != message.TypeTransaction {
				continue
			}
			body := msg.(*message.TransactionsMessage)
			require.Equal(t, txID, body.Transactions[0].ID())
			gotBroadcast = true
		case ev := <-lab.events.UnsafeGetChannel():
			published, ok := ev.(*tx.Tx)
			require.True(t, ok)
			require.Equal(t, txID, published.ID())
			gotEvent = true
		}
	}
}

func (lab *anchorLab) requirePool(t *testing.T, trx *tx.Tx, size int) {
	t.Helper()

	require.Equal(t, size, lab.pool.Size())
	require.True(t, lab.pool.HasTx(trx.ID()))
	require.Equal(t, trx, lab.pool.PendingTx(trx.ID()))
	found := false
	for _, pending := range lab.pool.AllPendingTxs() {
		if pending.ID() == trx.ID() {
			found = true
		}
	}
	require.True(t, found)
	require.Equal(t, 1, lab.pool.pools[payload.TypeAnchor].list.Size())
}

func (lab *anchorLab) requireAbsent(t *testing.T, trx *tx.Tx, size int) {
	t.Helper()

	require.Equal(t, size, lab.pool.Size())
	require.False(t, lab.pool.HasTx(trx.ID()))
}

func (lab *anchorLab) placeAnchor(acc *account.Account, locked amount.Amount, root byte) {
	require.NoError(lab.t, acc.SetAnchor(account.AnchorData{
		RootHash:        anchorRoot(root, 32),
		LockedDeposit:   locked,
		CreatedAtHeight: 1,
		CreatedAtTime:   1,
		UpdatedAtHeight: 1,
		UpdatedAtTime:   1,
	}))
}

func setLockedDeposit(acc *account.Account, locked amount.Amount) {
	type data struct {
		Number  int32
		Balance amount.Amount
		Anchor  *account.AnchorData
	}
	raw := (*data)(unsafe.Pointer(acc))
	raw.Anchor.LockedDeposit = locked
}

type unknownPayload struct {
	from crypto.Address
}

func (p unknownPayload) Signer() crypto.Address { return p.from }
func (unknownPayload) Value() amount.Amount     { return 0 }
func (unknownPayload) Type() payload.Type       { return payload.Type(8) }
func (unknownPayload) SerializeSize() int       { return 1 }
func (unknownPayload) Encode(w io.Writer) error {
	_, err := w.Write([]byte{0x08})

	return err
}
func (unknownPayload) Decode(payload.DecodeContext, io.Reader) error { return nil }
func (unknownPayload) BasicCheck() error                             { return nil }
func (unknownPayload) LogString() string                             { return "unknown" }

func swapPayload(trx *tx.Tx, pld payload.Payload) {
	type data struct {
		basicChecked bool
		Version      uint8
		LockTime     types.Height
		Fee          amount.Amount
		Memo         string
		Payload      payload.Payload
		Signature    crypto.Signature
		PublicKey    crypto.PublicKey
	}
	type exposed struct {
		memorizedID *tx.ID
		data        data
	}
	raw := (*exposed)(unsafe.Pointer(trx))
	raw.memorizedID = nil
	raw.data.basicChecked = false
	raw.data.Payload = pld
}

func TestAnchorEntersOnce(t *testing.T) {
	lab := newAnchorLab(t, protocol.ProtocolVersion5, nil)
	addr, prv, acc := lab.fund(executor.MinAnchorDeposit*4 + lab.fee()*4)
	before := acc.Balance()
	root := anchorRoot(0x21, 32)
	trx := lab.setTx(addr, prv, executor.MinAnchorDeposit, lab.fee(), root, "once", 0, lab.height, "")

	require.False(t, trx.IsFreeTx())
	require.NoError(t, lab.pool.AppendTx(trx))
	lab.requirePool(t, trx, 1)
	lab.requireQuiet(t)
	require.Equal(t, before-executor.MinAnchorDeposit-lab.fee(), acc.Balance())
	require.Equal(t, executor.MinAnchorDeposit, acc.LockedDeposit())
	require.Equal(t, root, acc.RootHash())
	require.Zero(t, lab.pool.pools[payload.TypeTransfer].list.Size())

	prepared := lab.pool.PrepareBlockTransactions()
	require.Len(t, prepared, 1)
	require.Equal(t, trx.ID(), prepared[0].ID())

	raw, err := trx.Bytes()
	require.NoError(t, err)
	clone, err := tx.FromBytes(raw)
	require.NoError(t, err)

	trx.Payload().(*payload.AnchorPayload).RootHash[0] ^= 0xff
	require.Equal(t, byte(0x21), acc.RootHash()[0])
	require.ErrorIs(t, lab.pool.AppendTx(clone), execution.TransactionCommittedError{ID: clone.ID()})
	require.Equal(t, 1, lab.pool.Size())
	require.True(t, lab.pool.HasTx(trx.ID()))
	require.Equal(t, executor.MinAnchorDeposit, acc.LockedDeposit())

	priced := lab.setTx(addr, prv, 1, lab.fee()+1, anchorRoot(0x22, 32), "plus", 0, lab.height, "")
	require.NoError(t, lab.pool.AppendTx(priced))
	require.Equal(t, before-executor.MinAnchorDeposit-lab.fee()-1-(lab.fee()+1), acc.Balance())
	require.Equal(t, executor.MinAnchorDeposit+1, acc.LockedDeposit())

	other, otherPrv, otherAcc := lab.fund(executor.MinAnchorDeposit + lab.fee())
	otherBefore := otherAcc.Balance()
	senderBefore := acc.Balance()
	otherTx := lab.setTx(other, otherPrv, executor.MinAnchorDeposit, lab.fee(),
		anchorRoot(0x23, 32), "other", 0, lab.height, "")
	require.NoError(t, lab.pool.AppendTx(otherTx))
	require.Equal(t, senderBefore, acc.Balance())
	require.Equal(t, otherBefore-executor.MinAnchorDeposit-lab.fee(), otherAcc.Balance())

	broadcast := lab.setTx(addr, prv, 0, lab.fee(), anchorRoot(0x24, 32), "air", 0, lab.height, "")
	lock := acc.LockedDeposit()
	balance := acc.Balance()
	require.NoError(t, lab.pool.AppendTxAndBroadcast(broadcast))
	lab.requirePublished(t, broadcast.ID())
	require.True(t, lab.pool.HasTx(broadcast.ID()))
	require.Equal(t, lock, acc.LockedDeposit())
	require.Equal(t, balance-lab.fee(), acc.Balance())

	unbond := lab.GenerateTestUnbondTx()
	sortition := lab.GenerateTestSortitionTx()
	require.NoError(t, lab.pool.AppendTx(unbond))
	require.NoError(t, lab.pool.AppendTx(sortition))
	require.True(t, lab.pool.HasTx(unbond.ID()))
	require.True(t, lab.pool.HasTx(sortition.ID()))
}

func TestAnchorSequences(t *testing.T) {
	lab := newAnchorLab(t, protocol.ProtocolVersion5, nil)
	addr, prv, acc := lab.fund(executor.MinAnchorDeposit*3 + lab.fee()*6)
	lab.placeAnchor(acc, executor.MinAnchorDeposit, 0x31)
	start := acc.Balance()
	oldLock := acc.LockedDeposit()

	first := lab.setTx(addr, prv, executor.MinAnchorDeposit, lab.fee(),
		anchorRoot(0x32, 32), "first", 0, lab.height, "one")
	second := lab.setTx(addr, prv, 1, lab.fee(), anchorRoot(0x33, 32), "second", 0, lab.height, "two")
	require.NoError(t, lab.pool.AppendTx(first))
	require.NoError(t, lab.pool.AppendTx(second))
	require.Equal(t, oldLock+executor.MinAnchorDeposit+1, acc.LockedDeposit())
	require.Equal(t, start-executor.MinAnchorDeposit-1-lab.fee()*2, acc.Balance())

	del := lab.deleteTx(addr, prv, lab.fee())
	require.NoError(t, lab.pool.AppendTx(del))
	require.False(t, acc.HasAnchor())
	require.Equal(t, amount.Amount(0), acc.LockedDeposit())
	require.Equal(t, start+oldLock-lab.fee()*3, acc.Balance())
	raw, err := acc.Bytes()
	require.NoError(t, err)
	require.Len(t, raw, 12)

	prepared := lab.pool.PrepareBlockTransactions()
	require.Equal(t, first.ID(), prepared[0].ID())
	require.Equal(t, second.ID(), prepared[1].ID())
	require.Equal(t, del.ID(), prepared[2].ID())
}

func TestAnchorDeleteThenSet(t *testing.T) {
	lab := newAnchorLab(t, protocol.ProtocolVersion5, nil)
	addr, prv, acc := lab.fund(executor.MinAnchorDeposit + lab.fee()*2)
	lab.placeAnchor(acc, executor.MinAnchorDeposit, 0x41)
	start := acc.Balance()
	oldLock := acc.LockedDeposit()

	del := lab.deleteTx(addr, prv, lab.fee())
	require.NoError(t, lab.pool.AppendTx(del))
	fresh := lab.setTx(addr, prv, executor.MinAnchorDeposit, lab.fee(), anchorRoot(0x42, 32), "new", 0, lab.height, "")
	require.NoError(t, lab.pool.AppendTx(fresh))
	require.Equal(t, executor.MinAnchorDeposit, acc.LockedDeposit())
	require.Equal(t, start+oldLock-lab.fee()*2-executor.MinAnchorDeposit, acc.Balance())
	require.Equal(t, anchorRoot(0x42, 32), acc.RootHash())
}

func TestAnchorDeleteRefundsLock(t *testing.T) {
	lab := newAnchorLab(t, protocol.ProtocolVersion5, nil)
	addr, prv, acc := lab.fund(5)
	lab.placeAnchor(acc, lab.fee(), 0x51)
	require.NoError(t, lab.pool.AppendTx(lab.deleteTx(addr, prv, lab.fee())))
	require.False(t, acc.HasAnchor())
	require.Equal(t, amount.Amount(5), acc.Balance())
	require.Equal(t, amount.Amount(0), acc.LockedDeposit())

	again, againPrv, againAcc := lab.fund(0)
	lab.placeAnchor(againAcc, lab.fee()-1, 0x52)
	rejected := lab.deleteTx(again, againPrv, lab.fee())
	require.ErrorIs(t, lab.pool.AppendTx(rejected), executor.ErrInsufficientFunds)
	lab.requireAbsent(t, rejected, 1)
	require.True(t, againAcc.HasAnchor())
	require.Equal(t, lab.fee()-1, againAcc.LockedDeposit())

	rich, richPrv, richAcc := lab.fund(9)
	lab.placeAnchor(richAcc, executor.MinAnchorDeposit, 0x53)
	huge := tx.NewAnchorTx(lab.height, rich, payload.AnchorActionDelete, nil, "", 0, 0, lab.fee())
	huge.Payload().(*payload.AnchorPayload).Deposit = amount.MaxNanoPAC
	lab.sign(richPrv, huge)
	require.NoError(t, lab.pool.AppendTx(huge))
	require.Equal(t, amount.Amount(9)+executor.MinAnchorDeposit-lab.fee(), richAcc.Balance())
	require.False(t, richAcc.HasAnchor())
}

func TestAnchorLowFeeSticksUntilNewSandbox(t *testing.T) {
	lab := newAnchorLab(t, protocol.ProtocolVersion5, nil)
	addr, prv, acc := lab.fund(executor.MinAnchorDeposit*2 + lab.fee()*2)
	before := acc.Balance()
	ghost := lab.setTx(addr, prv, executor.MinAnchorDeposit, 0, anchorRoot(0x61, 32), "ghost", 0, lab.height, "")

	require.ErrorIs(t, lab.pool.AppendTx(ghost), InvalidFeeError{MinimumFee: lab.fee()})
	lab.requireAbsent(t, ghost, 0)
	require.Equal(t, before-executor.MinAnchorDeposit, acc.Balance())
	require.Equal(t, executor.MinAnchorDeposit, acc.LockedDeposit())

	err := lab.pool.AppendTxAndBroadcast(ghost)
	require.ErrorIs(t, err, execution.TransactionCommittedError{ID: ghost.ID()})
	lab.requireQuiet(t)
	lab.requireAbsent(t, ghost, 0)

	follow := lab.setTx(addr, prv, executor.MinAnchorDeposit, lab.fee(), anchorRoot(0x62,
		32), "follow", 0, lab.height, "fixed")
	require.NoError(t, lab.pool.AppendTx(follow))
	require.True(t, lab.pool.HasTx(follow.ID()))
	require.False(t, lab.pool.HasTx(ghost.ID()))
	require.Equal(t, executor.MinAnchorDeposit*2, acc.LockedDeposit())

	lab.pool.removeTx(follow.ID())
	resetAcc := account.NewAccount(1)
	resetAcc.AddToBalance(before)
	lab.accounts = map[crypto.Address]*account.Account{addr: resetAcc}
	lab.committed = map[tx.ID]struct{}{}
	lab.sbx = lab.bindSandbox()
	lab.pool.SetNewSandboxAndRecheck(lab.sbx)
	require.False(t, lab.pool.HasTx(ghost.ID()))
	require.False(t, lab.pool.HasTx(follow.ID()))
	require.ErrorIs(t, lab.pool.AppendTx(ghost), InvalidFeeError{MinimumFee: lab.fee()})
	require.False(t, lab.pool.HasTx(ghost.ID()))
	require.Equal(t, before-executor.MinAnchorDeposit, resetAcc.Balance())
	require.Equal(t, executor.MinAnchorDeposit, resetAcc.LockedDeposit())
}

func TestAnchorRejectionsLeaveThePool(t *testing.T) {
	for _, version := range []protocol.Version{
		protocol.ProtocolVersionUnknown,
		protocol.ProtocolVersion3,
		protocol.ProtocolVersion4,
	} {
		lab := newAnchorLab(t, version, nil)
		addr, prv, acc := lab.fund(executor.MinAnchorDeposit + lab.fee())
		before := acc.Balance()
		for _, fee := range []amount.Amount{0, lab.fee()} {
			setTx := lab.setTx(addr, prv, executor.MinAnchorDeposit, fee, anchorRoot(0x71, 32), "no", 0, lab.height, "")
			require.ErrorIs(t, lab.pool.AppendTx(setTx), executor.InvalidPayloadTypeError{PayloadType: payload.TypeAnchor})
			require.ErrorIs(t, lab.pool.AppendTxAndBroadcast(setTx),
				executor.InvalidPayloadTypeError{PayloadType: payload.TypeAnchor})
			lab.requireQuiet(t)
			del := lab.deleteTx(addr, prv, fee)
			require.ErrorIs(t, lab.pool.AppendTx(del), executor.InvalidPayloadTypeError{PayloadType: payload.TypeAnchor})
			lab.requireAbsent(t, setTx, 0)
			require.Equal(t, before, acc.Balance())
			require.False(t, acc.HasAnchor())
		}
	}

	lab := newAnchorLab(t, protocol.ProtocolVersion5, nil)
	addr, prv, acc := lab.fund(amount.Amount(amount.MaxNanoPAC))
	before := acc.Balance()

	below := lab.setTx(addr, prv, executor.MinAnchorDeposit, lab.fee()-1, anchorRoot(0x72, 32), "low", 0, lab.height, "")
	require.ErrorIs(t, lab.pool.AppendTx(below), InvalidFeeError{MinimumFee: lab.fee()})
	lab.requireAbsent(t, below, 0)
	require.Equal(t, executor.MinAnchorDeposit, acc.LockedDeposit())

	airAddr, airPrv, airAcc := lab.fund(executor.MinAnchorDeposit)
	air := lab.setTx(airAddr, airPrv, executor.MinAnchorDeposit, 0, anchorRoot(0x72, 32), "air", 0, lab.height, "air")
	require.NoError(t, lab.pool.AppendTxAndBroadcast(air))
	lab.requirePublished(t, air.ID())
	require.False(t, lab.pool.HasTx(air.ID()))
	require.True(t, airAcc.HasAnchor())
	require.Equal(t, amount.Amount(0), airAcc.Balance())

	cases := append(make([]*tx.Tx, 0, 14),
		lab.setTx(addr, prv, executor.MinAnchorDeposit, lab.fee(), anchorRoot(0x73, 32), "zero-lock", 0, 0, ""),
		func() *tx.Tx {
			trx := tx.NewAnchorTx(lab.height, addr, payload.AnchorActionSet, anchorRoot(0x74,
				32), "nosig", 0, executor.MinAnchorDeposit, lab.fee())

			return trx
		}(),
		lab.setTx(addr, prv, executor.MinAnchorDeposit, -1, anchorRoot(0x75, 32), "neg", 0, lab.height, ""),
		lab.setTx(addr, prv, executor.MinAnchorDeposit, amount.Amount(amount.MaxNanoPAC)+1,
			anchorRoot(0x76, 32), "bigfee", 0, lab.height, ""),
		lab.setTx(addr, prv, executor.MinAnchorDeposit, lab.fee(), anchorRoot(0x77,
			32), "memo", 0, lab.height, strings.Repeat("m", 65)),
		lab.setTx(addr, prv, executor.MinAnchorDeposit, lab.fee(), anchorRoot(0x78, 31), "short", 0, lab.height, ""),
		lab.setTx(addr, prv, executor.MinAnchorDeposit, lab.fee(), anchorRoot(0x79, 65), "long", 0, lab.height, ""),
		lab.setTx(addr, prv, executor.MinAnchorDeposit, lab.fee(), anchorRoot(0x7a,
			32), strings.Repeat("u", 129), 0, lab.height, ""),
		lab.setTx(addr, prv, executor.MinAnchorDeposit, lab.fee(), anchorRoot(0x7b, 32), "bad\xC3", 0, lab.height, ""),
		lab.sign(prv, tx.NewAnchorTx(lab.height, addr, 2, anchorRoot(0x7c, 32),
			"act", 0, executor.MinAnchorDeposit, lab.fee())),
		lab.setTx(addr, prv, -1, lab.fee(), anchorRoot(0x7d, 32), "negdep", 0, lab.height, ""),
		lab.sign(prv, tx.NewAnchorTx(lab.height, crypto.TreasuryAddress, payload.AnchorActionSet,
			anchorRoot(0x7e, 32), "treasury", 0, executor.MinAnchorDeposit, lab.fee())),
		lab.sign(prv, tx.NewAnchorTx(lab.height, lab.RandValAddress(), payload.AnchorActionSet,
			anchorRoot(0x7f, 32), "val", 0, executor.MinAnchorDeposit, lab.fee())),
	)
	otherPub, _ := lab.RandBLSKeyPair()
	mismatch := tx.NewAnchorTx(lab.height, otherPub.AccountAddress(), payload.AnchorActionSet,
		anchorRoot(0x70, 32), "mismatch", 0, executor.MinAnchorDeposit, lab.fee())
	lab.sign(prv, mismatch)
	cases = append(cases, mismatch)

	for _, trx := range cases {
		err := lab.pool.AppendTx(trx)
		require.Error(t, err)
		require.NotErrorIs(t, err, ErrPayloadTypeNotSupported)
		lab.requireAbsent(t, trx, 0)
		require.Error(t, lab.pool.AppendTxAndBroadcast(trx))
	}
	lab.requireQuiet(t)
	require.Equal(t, before-executor.MinAnchorDeposit-(lab.fee()-1), acc.Balance())

	unknown := lab.setTx(addr, prv, executor.MinAnchorDeposit, lab.fee(), anchorRoot(0x11, 32), "swap", 0, lab.height, "")
	swapPayload(unknown, unknownPayload{from: addr})
	lab.sign(prv, unknown)
	require.ErrorIs(t, lab.pool.AppendTx(unknown), ErrPayloadTypeNotSupported)
	require.ErrorIs(t, lab.pool.AppendTxAndBroadcast(unknown), ErrPayloadTypeNotSupported)
	lab.requireQuiet(t)
	for _, amt := range []amount.Amount{-1, 0, 1, amount.Amount(amount.MaxNanoPAC)} {
		require.Equal(t, amount.Amount(0), lab.pool.EstimatedFee(amt, payload.Type(8)))
		require.Equal(t, lab.fee(), lab.pool.EstimatedFee(amt, payload.TypeAnchor))
		require.Equal(t, lab.fee(), lab.pool.EstimatedFee(amt, payload.TypeTransfer))
	}

	missingAddr, missingPrv, _ := lab.fund(0)
	delete(lab.accounts, missingAddr)
	missing := lab.setTx(missingAddr, missingPrv, executor.MinAnchorDeposit,
		lab.fee(), anchorRoot(0x12, 32), "gone", 0, lab.height, "")
	require.ErrorIs(t, lab.pool.AppendTx(missing), executor.AccountNotFoundError{Address: missingAddr})

	bannedAddr, bannedPrv, bannedAcc := lab.fund(executor.MinAnchorDeposit + lab.fee())
	lab.banned[bannedAddr] = true
	bannedTx := lab.setTx(bannedAddr, bannedPrv, executor.MinAnchorDeposit,
		lab.fee(), anchorRoot(0x13, 32), "ban", 0, lab.height, "")
	require.ErrorIs(t, lab.pool.AppendTx(bannedTx), execution.SignerBannedError{Address: bannedAddr})
	require.False(t, bannedAcc.HasAnchor())
	require.Equal(t, executor.MinAnchorDeposit+lab.fee(), bannedAcc.Balance())

	cleanAddr, cleanPrv, cleanAcc := lab.fund(executor.MinAnchorDeposit + lab.fee())
	cleanBefore := cleanAcc.Balance()
	small := lab.setTx(cleanAddr, cleanPrv, 0, lab.fee(), anchorRoot(0x14, 32), "zero", 0, lab.height, "")
	require.ErrorIs(t, lab.pool.AppendTx(small), executor.ErrAnchorDepositTooSmall)
	tiny := lab.setTx(cleanAddr, cleanPrv, executor.MinAnchorDeposit-1, lab.fee(),
		anchorRoot(0x15, 32), "tiny", 0, lab.height, "")
	require.ErrorIs(t, lab.pool.AppendTx(tiny), executor.ErrAnchorDepositTooSmall)
	require.False(t, cleanAcc.HasAnchor())
	require.Equal(t, cleanBefore, cleanAcc.Balance())

	poorAddr, poorPrv, poorAcc := lab.fund(executor.MinAnchorDeposit)
	poor := lab.setTx(poorAddr, poorPrv, executor.MinAnchorDeposit, lab.fee(),
		anchorRoot(0x16, 32), "poor", 0, lab.height, "")
	require.ErrorIs(t, lab.pool.AppendTx(poor), executor.ErrInsufficientFunds)
	require.False(t, poorAcc.HasAnchor())

	bareAddr, barePrv, bareAcc := lab.fund(lab.fee())
	bareBefore := bareAcc.Balance()
	noAnchor := lab.deleteTx(bareAddr, barePrv, lab.fee())
	require.ErrorIs(t, lab.pool.AppendTx(noAnchor), executor.ErrAnchorNotFound)
	require.Equal(t, bareBefore, bareAcc.Balance())

	interval := lab.sbx.Params().TransactionToLiveInterval
	expired := lab.setTx(addr, prv, executor.MinAnchorDeposit, lab.fee(),
		anchorRoot(0x17, 32), "old", 0, lab.height-types.Height(interval)-1, "")
	require.ErrorIs(t, lab.pool.AppendTx(expired), execution.LockTimeExpiredError{LockTime: expired.LockTime()})
	require.False(t, lab.pool.HasTx(expired.ID()))
}

// PIP-50 section 7: anchors come after all other transactions in the block
// template, so they are the first left out of a full block.
func TestAnchorsComeLastInTheBlockTemplate(t *testing.T) {
	lab := newAnchorLab(t, protocol.ProtocolVersion5, nil)
	addr, prv, _ := lab.fund(executor.MinAnchorDeposit + lab.fee())

	anchorTx := lab.setTx(addr, prv, executor.MinAnchorDeposit, lab.fee(), anchorRoot(0x61, 32),
		"first in", 0, lab.height, "")
	require.NoError(t, lab.pool.AppendTx(anchorTx))
	transfer := lab.GenerateTestTransferTx(
		testsuite.TransactionWithLockTime(lab.height),
		testsuite.TransactionWithFee(lab.fee()),
	)
	require.NoError(t, lab.pool.AppendTx(transfer))

	prepared := lab.pool.PrepareBlockTransactions()
	require.Len(t, prepared, 2)
	require.Equal(t, transfer.ID(), prepared[0].ID())
	require.Equal(t, anchorTx.ID(), prepared[1].ID())
}

func TestAnchorPoolLimits(t *testing.T) {
	conf := testDefaultConfig()
	conf.MaxSize = 10
	lab := newAnchorLab(t, protocol.ProtocolVersion5, conf)
	addr, prv, acc := lab.fund(executor.MinAnchorDeposit*8 + lab.fee()*8)
	capSize := lab.pool.config.anchorPoolSize()
	kept := make([]*tx.Tx, 0, capSize+1)
	for i := 0; i < capSize+1; i++ {
		trx := lab.setTx(addr, prv, executor.MinAnchorDeposit, lab.fee(), anchorRoot(byte(0x80+i),
			32), "full", 0, lab.height, strings.Repeat("n", i+1))
		require.NoError(t, lab.pool.AppendTx(trx))
		kept = append(kept, trx)
	}
	require.Equal(t, capSize, lab.pool.pools[payload.TypeAnchor].list.Size())
	require.False(t, lab.pool.HasTx(kept[0].ID()))
	require.True(t, lab.pool.HasTx(kept[1].ID()))
	require.True(t, lab.pool.HasTx(kept[capSize].ID()))

	transfer := lab.GenerateTestTransferTx(
		testsuite.TransactionWithLockTime(lab.height),
		testsuite.TransactionWithFee(lab.fee()),
	)
	require.NoError(t, lab.pool.AppendTx(transfer))
	extra := lab.setTx(addr, prv, 1, lab.fee(), anchorRoot(0x90, 32), "extra", 0, lab.height, "extra")
	require.NoError(t, lab.pool.AppendTx(extra))
	require.True(t, lab.pool.HasTx(transfer.ID()))
	require.False(t, lab.pool.HasTx(kept[1].ID()))
	require.Equal(t, capSize, lab.pool.pools[payload.TypeAnchor].list.Size())

	for i := 0; i < lab.pool.config.transferPoolSize()+1; i++ {
		trx := lab.GenerateTestTransferTx(
			testsuite.TransactionWithLockTime(lab.height),
			testsuite.TransactionWithFee(lab.fee()),
		)
		require.NoError(t, lab.pool.AppendTx(trx))
	}
	require.True(t, lab.pool.HasTx(extra.ID()))

	broken := lab.setTx(addr, prv, 1, lab.fee(), anchorRoot(0x91, 32), "broken", 0, lab.height, "broken")
	broken.Payload().(*payload.AnchorPayload).RootHash[0] ^= 0xff
	broken.SetSignature(broken.Signature())
	require.Error(t, lab.pool.AppendTx(broken))
	require.True(t, lab.pool.HasTx(extra.ID()))
	require.Equal(t, capSize, lab.pool.pools[payload.TypeAnchor].list.Size())
	require.True(t, acc.HasAnchor())
}

func TestAnchorCommittedAndRecheck(t *testing.T) {
	lab := newAnchorLab(t, protocol.ProtocolVersion5, nil)
	addr, prv, acc := lab.fund(executor.MinAnchorDeposit*2 + lab.fee()*2)
	first := lab.setTx(addr, prv, executor.MinAnchorDeposit, lab.fee(), anchorRoot(0xa1, 32), "keep", 0, lab.height, "a")
	second := lab.setTx(addr, prv, 1, lab.fee(), anchorRoot(0xa2, 32), "stay", 0, lab.height, "b")
	require.NoError(t, lab.pool.AppendTx(first))
	require.NoError(t, lab.pool.AppendTx(second))
	transfer := lab.GenerateTestTransferTx(
		testsuite.TransactionWithLockTime(lab.height),
		testsuite.TransactionWithFee(lab.fee()),
	)
	require.NoError(t, lab.pool.AppendTx(transfer))
	require.Equal(t, first.SerializeSize(), lab.pool.getPendingConsumption(addr)-second.SerializeSize())

	subsidy := lab.makeSubsidy(first.LockTime())
	blk, _ := lab.GenerateTestBlock(1, testsuite.BlockWithTransactions([]*tx.Tx{subsidy, first}))
	lab.pool.HandleCommittedBlock(blk)
	require.False(t, lab.pool.HasTx(first.ID()))
	require.True(t, lab.pool.HasTx(second.ID()))
	require.True(t, lab.pool.HasTx(transfer.ID()))
	require.Equal(t, executor.MinAnchorDeposit+1, acc.LockedDeposit())

	empty, _ := lab.GenerateTestBlock(2, testsuite.BlockWithTransactions([]*tx.Tx{subsidy}))
	lab.pool.HandleCommittedBlock(empty)
	require.True(t, lab.pool.HasTx(second.ID()))

	require.ErrorIs(t, lab.pool.AppendTx(first), execution.TransactionCommittedError{ID: first.ID()})
	require.False(t, lab.pool.HasTx(first.ID()))

	fresh := account.NewAccount(acc.Number())
	fresh.AddToBalance(1 + lab.fee())
	lab.placeAnchor(fresh, executor.MinAnchorDeposit, 0xaa)
	lab.recheck(protocol.ProtocolVersion5, lab.height, map[crypto.Address]*account.Account{addr: fresh}, false)
	require.True(t, lab.pool.HasTx(second.ID()))
	require.True(t, lab.pool.HasTx(transfer.ID()))
	require.Equal(t, executor.MinAnchorDeposit+1, fresh.LockedDeposit())

	untouched := account.NewAccount(acc.Number())
	untouched.AddToBalance(executor.MinAnchorDeposit)
	lab.recheck(protocol.ProtocolVersion4, lab.height, map[crypto.Address]*account.Account{addr: untouched}, false)
	require.False(t, lab.pool.HasTx(second.ID()))
	require.True(t, lab.pool.HasTx(transfer.ID()))
	require.False(t, untouched.HasAnchor())
	require.Equal(t, executor.MinAnchorDeposit, untouched.Balance())
}

func TestAnchorRecheckDropsBadState(t *testing.T) {
	lab := newAnchorLab(t, protocol.ProtocolVersion5, nil)
	addr, prv, _ := lab.fund(executor.MinAnchorDeposit + lab.fee())
	trx := lab.setTx(addr, prv, executor.MinAnchorDeposit, lab.fee(), anchorRoot(0xb1, 32), "re", 0, lab.height, "")
	require.NoError(t, lab.pool.AppendTx(trx))

	interval := lab.sbx.Params().TransactionToLiveInterval
	later := account.NewAccount(1)
	later.AddToBalance(executor.MinAnchorDeposit + lab.fee())
	lab.recheck(protocol.ProtocolVersion5, lab.height+types.Height(interval)+1,
		map[crypto.Address]*account.Account{addr: later}, false)
	require.False(t, lab.pool.HasTx(trx.ID()))
	require.False(t, later.HasAnchor())

	lab.height = anchorTestHeight
	require.NoError(t, lab.pool.AppendTx(trx))
	bannedAcc := account.NewAccount(1)
	bannedAcc.AddToBalance(executor.MinAnchorDeposit + lab.fee())
	lab.banned[addr] = true
	lab.recheck(protocol.ProtocolVersion5, anchorTestHeight, map[crypto.Address]*account.Account{addr: bannedAcc}, true)
	require.False(t, lab.pool.HasTx(trx.ID()))
	require.Equal(t, executor.MinAnchorDeposit+lab.fee(), bannedAcc.Balance())

	lab.banned[addr] = false
	ready := account.NewAccount(1)
	ready.AddToBalance(executor.MinAnchorDeposit + lab.fee())
	lab.accounts = map[crypto.Address]*account.Account{addr: ready}
	require.NoError(t, lab.pool.AppendTx(trx))
	lab.recheck(protocol.ProtocolVersion5, anchorTestHeight, map[crypto.Address]*account.Account{}, false)
	require.False(t, lab.pool.HasTx(trx.ID()))

	ready = account.NewAccount(1)
	ready.AddToBalance(executor.MinAnchorDeposit + lab.fee())
	lab.accounts = map[crypto.Address]*account.Account{addr: ready}
	lab.committed = map[tx.ID]struct{}{}
	require.NoError(t, lab.pool.AppendTx(trx))
	poor := account.NewAccount(1)
	poor.AddToBalance(1)
	lab.recheck(protocol.ProtocolVersion5, lab.height, map[crypto.Address]*account.Account{addr: poor}, false)
	require.False(t, lab.pool.HasTx(trx.ID()))
	require.Equal(t, amount.Amount(1), poor.Balance())
	require.False(t, poor.HasAnchor())
}

func TestAnchorAbnormal(t *testing.T) {
	lab := newAnchorLab(t, protocol.Version(6), nil)
	require.Equal(t, protocol.Version(6), lab.version)
	addr, prv, acc := lab.fund(executor.MinAnchorDeposit*2 + lab.fee()*4)

	zero := lab.setTx(addr, prv, executor.MinAnchorDeposit, lab.fee(), anchorRoot(0x00,
		32), "bad\nuri", 0x00, lab.height, "")
	require.NoError(t, lab.pool.AppendTx(zero))
	require.Equal(t, "bad\nuri", acc.ManifestURI())
	opaque := lab.setTx(addr, prv, 0, lab.fee(), anchorRoot(0xab, 32), "", 0xFF, lab.height, "")
	require.NoError(t, lab.pool.AppendTx(opaque))
	require.Equal(t, byte(0xFF), acc.AnchorType())

	edPub, edPrv := lab.RandEd25519KeyPair()
	edAcc := account.NewAccount(8)
	edAcc.AddToBalance(executor.MinAnchorDeposit + lab.fee())
	edAddr := edPub.AccountAddress()
	lab.accounts[edAddr] = edAcc
	require.NoError(t, lab.pool.AppendTx(lab.setTx(edAddr, edPrv, executor.MinAnchorDeposit,
		lab.fee(), anchorRoot(0xc1, 32), "ed", 0, lab.height, "")))
	require.True(t, edAcc.HasAnchor())

	secPub, secPrv := lab.RandSecp256k1KeyPair()
	secAcc := account.NewAccount(9)
	secAcc.AddToBalance(executor.MinAnchorDeposit + lab.fee())
	secAddr := secPub.AccountAddress()
	lab.accounts[secAddr] = secAcc
	require.NoError(t, lab.pool.AppendTx(lab.setTx(secAddr, secPrv, executor.MinAnchorDeposit,
		lab.fee(), anchorRoot(0xc2, 32), "sec", 0, lab.height, "")))
	require.True(t, secAcc.HasAnchor())

	delAddr, delPrv, delAcc := lab.fund(3)
	lab.placeAnchor(delAcc, executor.MinAnchorDeposit, 0xc3)
	firstDel := lab.deleteTx(delAddr, delPrv, lab.fee())
	require.NoError(t, lab.pool.AppendTx(firstDel))
	secondDel := lab.sign(delPrv, tx.NewAnchorTx(lab.height, delAddr, payload.AnchorActionDelete,
		nil, "", 0, 0, lab.fee(), tx.WithMemo("again")))
	require.ErrorIs(t, lab.pool.AppendTx(secondDel), executor.ErrAnchorNotFound)
	require.Equal(t, amount.Amount(3)+executor.MinAnchorDeposit-lab.fee(), delAcc.Balance())
	require.True(t, lab.pool.HasTx(firstDel.ID()))
	require.False(t, lab.pool.HasTx(secondDel.ID()))

	cleanAddr, cleanPrv, cleanAcc := lab.fund(executor.MinAnchorDeposit + lab.fee())
	miss := lab.deleteTx(cleanAddr, cleanPrv, lab.fee())
	require.ErrorIs(t, lab.pool.AppendTx(miss), executor.ErrAnchorNotFound)
	require.NoError(t, lab.pool.AppendTx(lab.setTx(cleanAddr, cleanPrv, executor.MinAnchorDeposit,
		lab.fee(), anchorRoot(0xc4, 32), "after", 0, lab.height, "")))
	require.True(t, cleanAcc.HasAnchor())

	balAddr, balPrv, balAcc := lab.fund(11)
	lab.placeAnchor(balAcc, executor.MinAnchorDeposit, 0xc5)
	require.NoError(t, lab.pool.AppendTx(lab.deleteTx(balAddr, balPrv, lab.fee())))
	require.Equal(t, amount.Amount(11)+executor.MinAnchorDeposit-lab.fee(), balAcc.Balance())
	require.False(t, balAcc.HasAnchor())

	upAddr, upPrv, upAcc := lab.fund(lab.fee() * 2)
	lab.placeAnchor(upAcc, executor.MinAnchorDeposit, 0xc6)
	firstUp := lab.setTx(upAddr, upPrv, 0, lab.fee(), anchorRoot(0xc7, 32), "u1", 0, lab.height, "u1")
	secondUp := lab.setTx(upAddr, upPrv, 0, lab.fee(), anchorRoot(0xc8, 32), "u2", 0, lab.height, "u2")
	require.NoError(t, lab.pool.AppendTx(firstUp))
	require.NoError(t, lab.pool.AppendTx(secondUp))
	require.True(t, lab.pool.HasTx(firstUp.ID()))
	require.True(t, lab.pool.HasTx(secondUp.ID()))
	require.Equal(t, anchorRoot(0xc8, 32), upAcc.RootHash())
	require.Equal(t, executor.MinAnchorDeposit, upAcc.LockedDeposit())

	future := lab.setTx(addr, prv, 0, lab.fee(), anchorRoot(0xc9, 32), "future", 0, lab.height+1, "future")
	require.NoError(t, lab.pool.AppendTx(future))
	far := lab.setTx(addr, prv, 0, lab.fee(), anchorRoot(0xca, 32), "far", 0, types.Height(^uint32(0)), "far")
	require.NoError(t, lab.pool.AppendTx(far))
}

func TestAnchorValueBounds(t *testing.T) {
	maxAmt := amount.Amount(amount.MaxNanoPAC)
	lab := newAnchorLab(t, protocol.ProtocolVersion5, nil)
	addr, prv, acc := lab.fund(maxAmt)

	exact := lab.setTx(addr, prv, executor.MinAnchorDeposit, lab.fee(), anchorRoot(0xd1, 32), "min", 0, lab.height, "")
	poor, poorPrv, poorAcc := lab.fund(executor.MinAnchorDeposit + lab.fee() - 1)
	require.ErrorIs(t, lab.pool.AppendTx(lab.setTx(poor, poorPrv, executor.MinAnchorDeposit,
		lab.fee(), anchorRoot(0xd2, 32), "short", 0, lab.height, "")), executor.ErrInsufficientFunds)
	require.False(t, poorAcc.HasAnchor())
	pocket, pocketPrv, pocketAcc := lab.fund(executor.MinAnchorDeposit + lab.fee())
	require.NoError(t, lab.pool.AppendTx(lab.setTx(pocket, pocketPrv, executor.MinAnchorDeposit,
		lab.fee(), anchorRoot(0xd3, 32), "exact", 0, lab.height, "")))
	require.Equal(t, amount.Amount(0), pocketAcc.Balance())
	require.Equal(t, executor.MinAnchorDeposit, pocketAcc.LockedDeposit())

	require.NoError(t, lab.pool.AppendTx(exact))

	topAddr, topPrv, topAcc := lab.fund(maxAmt)
	top := lab.setTx(topAddr, topPrv, maxAmt-lab.fee(), lab.fee(), anchorRoot(0xd5,
		32), strings.Repeat("a", 128), 0, lab.height, strings.Repeat("m", 64))
	require.NoError(t, lab.pool.AppendTx(top))
	require.Equal(t, amount.Amount(0), topAcc.Balance())
	require.Equal(t, maxAmt-lab.fee(), topAcc.LockedDeposit())

	overAddr, overPrv, overAcc := lab.fund(maxAmt)
	over := lab.setTx(overAddr, overPrv, maxAmt-lab.fee()+1, lab.fee(), anchorRoot(0xd6, 32), "over", 0, lab.height, "")
	require.ErrorIs(t, lab.pool.AppendTx(over), executor.ErrAmountOverflow)
	require.False(t, overAcc.HasAnchor())
	require.Equal(t, maxAmt, overAcc.Balance())

	lockedAddr, lockedPrv, lockedAcc := lab.fund(maxAmt)
	lab.placeAnchor(lockedAcc, maxAmt, 0xd7)
	require.NoError(t, lab.pool.AppendTx(lab.setTx(lockedAddr, lockedPrv,
		0, lab.fee(), anchorRoot(0xd8, 32), "", 0, lab.height, "")))
	require.Equal(t, maxAmt, lockedAcc.LockedDeposit())
	bump := lab.setTx(lockedAddr, lockedPrv, 1, lab.fee(), anchorRoot(0xd9, 32), "bump", 0, lab.height, "bump")
	require.ErrorIs(t, lab.pool.AppendTx(bump), executor.ErrAmountOverflow)
	require.Equal(t, maxAmt, lockedAcc.LockedDeposit())

	feeAddr, feePrv, feeAcc := lab.fund(maxAmt)
	lab.placeAnchor(feeAcc, 1, 0xda)
	maxFee := lab.setTx(feeAddr, feePrv, 0, maxAmt, anchorRoot(0xdb, 32), "maxfee", 0, lab.height, "maxfee")
	require.NoError(t, lab.pool.AppendTx(maxFee))
	require.Equal(t, amount.Amount(0), feeAcc.Balance())
	require.Equal(t, amount.Amount(1), feeAcc.LockedDeposit())

	refund := executor.MinAnchorDeposit - lab.fee()
	okAddr, okPrv, okAcc := lab.fund(maxAmt - refund)
	lab.placeAnchor(okAcc, executor.MinAnchorDeposit, 0xdc)
	require.NoError(t, lab.pool.AppendTx(lab.deleteTx(okAddr, okPrv, lab.fee())))
	require.Equal(t, maxAmt, okAcc.Balance())
	require.False(t, okAcc.HasAnchor())

	badAddr, badPrv, badAcc := lab.fund(maxAmt - refund + 1)
	lab.placeAnchor(badAcc, executor.MinAnchorDeposit, 0xdd)
	badDel := lab.deleteTx(badAddr, badPrv, lab.fee())
	require.ErrorIs(t, lab.pool.AppendTx(badDel), executor.ErrAmountOverflow)
	require.True(t, badAcc.HasAnchor())
	require.Equal(t, maxAmt-refund+1, badAcc.Balance())

	hash64 := lab.setTx(addr, prv, 0, lab.fee(), anchorRoot(0xde, 64), "", 0, lab.height, "")
	require.NoError(t, lab.pool.AppendTx(hash64))
	require.Len(t, acc.RootHash(), 64)

	interval := lab.sbx.Params().TransactionToLiveInterval
	oldest := lab.height - types.Height(interval)
	require.NoError(t, lab.pool.AppendTx(lab.setTx(addr, prv, 0, lab.fee(),
		anchorRoot(0xdf, 32), "edge", 0, oldest, "edge")))
	require.NoError(t, lab.pool.AppendTx(lab.setTx(addr, prv, 0, lab.fee(),
		anchorRoot(0xe0, 32), "now", 0, lab.height, "now")))
	require.ErrorIs(t,
		lab.pool.AppendTx(lab.setTx(addr, prv, 0, lab.fee(), anchorRoot(0xe1, 32), "past", 0, oldest-1, "past")),
		execution.LockTimeExpiredError{LockTime: oldest - 1})

	negAddr, negPrv, negAcc := lab.fund(lab.fee())
	lab.placeAnchor(negAcc, 1, 0xe2)
	setLockedDeposit(negAcc, -1)
	neg := lab.setTx(negAddr, negPrv, 0, lab.fee(), anchorRoot(0xe3, 32), "neg", 0, lab.height, "neg")
	require.ErrorIs(t, lab.pool.AppendTx(neg), executor.ErrAmountOverflow)
	require.Equal(t, amount.Amount(-1), negAcc.LockedDeposit())
	lab.requireAbsent(t, neg, lab.pool.Size())
}

func TestAnchorConsumption(t *testing.T) {
	cfg := testConsumptionalConfig()
	cfg.Fee.FixedFee = 0.01
	cfg.MaxSize = 100
	lab := newAnchorLab(t, protocol.ProtocolVersion5, cfg)
	addr, prv, _ := lab.fund(executor.MinAnchorDeposit + lab.fee()*4)
	lab.indexed[addr] = true
	trx := lab.setTx(addr, prv, executor.MinAnchorDeposit, lab.fee(), anchorRoot(0xf1, 32), "use", 0, lab.height, "")
	lab.pool.consumptionMap[addr] = cfg.Fee.DailyLimit - trx.SerializeSize() - 1
	require.NoError(t, lab.pool.AppendTx(trx))
	require.Equal(t, trx.SerializeSize(), lab.pool.getPendingConsumption(addr))

	other, otherPrv, _ := lab.fund(executor.MinAnchorDeposit + lab.fee()*4)
	lab.indexed[other] = true
	next := lab.setTx(other, otherPrv, executor.MinAnchorDeposit, lab.fee(),
		anchorRoot(0xf2, 32), "limit", 0, lab.height, "")
	lab.pool.consumptionMap[other] = cfg.Fee.DailyLimit - next.SerializeSize()
	require.Error(t, lab.pool.AppendTx(next))
	require.False(t, lab.pool.HasTx(next.ID()))
	require.Equal(t, lab.fee(), lab.pool.EstimatedFee(1, payload.TypeAnchor))

	stranger, strangerPrv, _ := lab.fund(executor.MinAnchorDeposit + lab.fee()*4)
	outside := lab.setTx(stranger, strangerPrv, executor.MinAnchorDeposit,
		lab.fee(), anchorRoot(0xf3, 32), "out", 0, lab.height, "")
	require.Error(t, lab.pool.AppendTx(outside))
	require.False(t, lab.pool.HasTx(outside.ID()))

	countLab := newAnchorLab(t, protocol.ProtocolVersion5, cfg)
	countAddr, countPrv, _ := countLab.fund(countLab.fee() + executor.MinAnchorDeposit)
	countLab.indexed[countAddr] = true
	counted := countLab.setTx(countAddr, countPrv, executor.MinAnchorDeposit,
		countLab.fee(), anchorRoot(0xf4, 32), "count", 0, countLab.height,
		"")
	require.NoError(t, countLab.pool.AppendTx(counted))
	subsidy := countLab.GenerateTestSubsidyTx(testsuite.TransactionWithLockTime(countLab.height))
	blk, _ := countLab.GenerateTestBlock(1, testsuite.BlockWithTransactions([]*tx.Tx{subsidy, counted}))
	countLab.pool.HandleCommittedBlock(blk)
	require.Equal(t, counted.SerializeSize(), countLab.pool.consumptionMap[countAddr])
	require.NotPanics(t, func() {
		countLab.pool.HandleCommittedBlock(blk)
	})
}

func TestAnchorRealSandbox(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	fakeStore := store.NewFakeStore(ts)
	fakeStore.EXPECT().IsBanned(gomock.Any()).Return(false).AnyTimes()
	fakeStore.EXPECT().RecentTransaction(gomock.Any()).DoAndReturn(func(id tx.ID) bool {
		_, ok := fakeStore.RecentTxs[id]

		return ok
	}).AnyTimes()

	params := param.FromGenesis(genesis.MainnetGenesis())
	params.BlockVersion = protocol.ProtocolVersion5
	fee, _ := amount.NewAmount(0.01)
	pub, prv := ts.RandBLSKeyPair()
	addr := pub.AccountAddress()
	acc := account.NewAccount(0)
	start := executor.MinAnchorDeposit + fee*4
	acc.AddToBalance(start)
	fakeStore.FakeAccounts[addr] = acc

	previous := executor.DefaultFactory
	t.Cleanup(func() {
		executor.DefaultFactory = previous
	})
	executor.DefaultFactory = executor.MakeExecutorImpl

	sbx := sandbox.NewSandbox(anchorTestHeight-1, params.BlockVersion, 0, fakeStore, params, nil, 0)
	broadcast := pipeline.New[message.Message](t.Context())
	events := pipeline.New[any](t.Context())
	pool := NewTxPool(t.Context(), testDefaultConfig(), fakeStore, broadcast, events).(*txPool)
	pool.SetNewSandboxAndRecheck(sbx)

	root := bytes.Repeat([]byte{0x44}, 32)
	setTx := tx.NewAnchorTx(sbx.CurrentHeight(), addr, payload.AnchorActionSet,
		root, "real", 0, executor.MinAnchorDeposit, fee)
	ts.HelperSignTransaction(prv, setTx)
	require.NoError(t, pool.AppendTx(setTx))
	seen := sbx.Account(addr)
	require.Equal(t, executor.MinAnchorDeposit, seen.LockedDeposit())
	require.Equal(t, start-executor.MinAnchorDeposit-fee, seen.Balance())
	require.Equal(t, root, seen.RootHash())

	over := tx.NewTransferTx(sbx.CurrentHeight(), addr, ts.RandAccAddress(), start, fee)
	ts.HelperSignTransaction(prv, over)
	require.ErrorIs(t, pool.AppendTx(over), executor.ErrInsufficientFunds)
	require.Equal(t, executor.MinAnchorDeposit, sbx.Account(addr).LockedDeposit())

	remain := sbx.Account(addr).Balance() - fee
	spend := tx.NewTransferTx(sbx.CurrentHeight(), addr, ts.RandAccAddress(), remain, fee)
	ts.HelperSignTransaction(prv, spend)
	require.NoError(t, pool.AppendTx(spend))
	require.Equal(t, amount.Amount(0), sbx.Account(addr).Balance())
	require.Equal(t, executor.MinAnchorDeposit, sbx.Account(addr).LockedDeposit())

	del := tx.NewAnchorTx(sbx.CurrentHeight(), addr, payload.AnchorActionDelete, nil, "", 0, 0, fee)
	ts.HelperSignTransaction(prv, del)
	require.NoError(t, pool.AppendTx(del))
	require.False(t, sbx.Account(addr).HasAnchor())
	require.Equal(t, executor.MinAnchorDeposit-fee, sbx.Account(addr).Balance())
	raw, err := sbx.Account(addr).Bytes()
	require.NoError(t, err)
	require.Len(t, raw, 12)

	require.ErrorIs(t, pool.AppendTx(setTx), execution.TransactionCommittedError{ID: setTx.ID()})
	require.Equal(t, executor.MinAnchorDeposit-fee, sbx.Account(addr).Balance())

	ghostPub, ghostPrv := ts.RandBLSKeyPair()
	ghostAddr := ghostPub.AccountAddress()
	ghostAcc := account.NewAccount(1)
	ghostAcc.AddToBalance(executor.MinAnchorDeposit)
	fakeStore.FakeAccounts[ghostAddr] = ghostAcc
	ghost := tx.NewAnchorTx(sbx.CurrentHeight(), ghostAddr, payload.AnchorActionSet,
		bytes.Repeat([]byte{0x45}, 32), "ghost", 0, executor.MinAnchorDeposit, 0)
	ts.HelperSignTransaction(ghostPrv, ghost)
	require.ErrorIs(t, pool.AppendTx(ghost), InvalidFeeError{MinimumFee: fee})
	require.False(t, pool.HasTx(ghost.ID()))
	require.True(t, sbx.Account(ghostAddr).HasAnchor())
	require.Equal(t, amount.Amount(0), sbx.Account(ghostAddr).Balance())
}

func TestAnchorSandboxOutlivesTheQueue(t *testing.T) {
	conf := testDefaultConfig()
	conf.MaxSize = 30 // three anchor slots, so the queue keeps more than one tx
	lab := newAnchorLab(t, protocol.ProtocolVersion5, conf)
	count := lab.pool.config.anchorPoolSize() + 1
	addr, prv, acc := lab.fund(executor.MinAnchorDeposit + amount.Amount(count) + lab.fee()*amount.Amount(count))
	queued := make([]*tx.Tx, 0, count)
	for i := range count {
		deposit := amount.Amount(1)
		if i == 0 {
			deposit = executor.MinAnchorDeposit
		}
		trx := lab.setTx(addr, prv, deposit, lab.fee(), anchorRoot(byte(0x30+i), 32),
			"q", 0, lab.height, strings.Repeat("q", i+1))
		require.NoError(t, lab.pool.AppendTx(trx))
		queued = append(queued, trx)
	}

	locked := executor.MinAnchorDeposit + amount.Amount(count-1)
	require.False(t, lab.pool.HasTx(queued[0].ID()))
	require.Equal(t, locked, acc.LockedDeposit())
	pending := 0
	for _, trx := range queued[1:] {
		pending += trx.SerializeSize()
	}
	require.Equal(t, pending, lab.pool.getPendingConsumption(addr))
	require.ErrorIs(t, lab.pool.AppendTx(queued[0]), execution.TransactionCommittedError{ID: queued[0].ID()})
	require.Error(t, lab.pool.AppendTxAndBroadcast(queued[0]))
	lab.requireQuiet(t)
	require.Equal(t, locked, acc.LockedDeposit())

	balance := acc.Balance()
	lab.pool.removeTx(queued[1].ID())
	require.False(t, lab.pool.HasTx(queued[1].ID()))
	require.Equal(t, balance, acc.Balance())
	require.Equal(t, locked, acc.LockedDeposit())

	require.NoError(t, lab.pool.AppendTx(lab.deleteTx(addr, prv, lab.fee())))
	require.False(t, acc.HasAnchor())
	require.Equal(t, amount.Amount(0), acc.LockedDeposit())
	require.Equal(t, balance+locked-lab.fee(), acc.Balance())
	require.True(t, lab.pool.HasTx(queued[count-1].ID()))
	require.False(t, lab.pool.HasTx(queued[0].ID()))

	priorAddr, priorPrv, priorAcc := lab.fund(executor.MinAnchorDeposit + lab.fee())
	prior := lab.setTx(priorAddr, priorPrv, executor.MinAnchorDeposit, lab.fee(),
		anchorRoot(0x3a, 32), "prior", 0, lab.height, "prior")
	lab.committed[prior.ID()] = struct{}{}
	priorBalance := priorAcc.Balance()
	require.ErrorIs(t, lab.pool.AppendTx(prior), execution.TransactionCommittedError{ID: prior.ID()})
	require.Equal(t, priorBalance, priorAcc.Balance())
	require.False(t, priorAcc.HasAnchor())
	require.False(t, lab.pool.HasTx(prior.ID()))

	ghostAddr, ghostPrv, ghostAcc := lab.fund(executor.MinAnchorDeposit)
	ghost := lab.setTx(ghostAddr, ghostPrv, executor.MinAnchorDeposit, 0,
		anchorRoot(0x3b, 32), "ghost", 0, lab.height, "ghost")
	require.ErrorIs(t, lab.pool.AppendTx(ghost), InvalidFeeError{MinimumFee: lab.fee()})
	require.False(t, lab.pool.HasTx(ghost.ID()))
	require.Equal(t, executor.MinAnchorDeposit, ghostAcc.LockedDeposit())
	require.NoError(t, lab.pool.AppendTx(lab.deleteTx(ghostAddr, ghostPrv, lab.fee())))
	require.False(t, ghostAcc.HasAnchor())
	require.False(t, lab.pool.HasTx(ghost.ID()))
	require.Equal(t, executor.MinAnchorDeposit-lab.fee(), ghostAcc.Balance())
}

func TestAnchorPendingBytesRaiseTheFee(t *testing.T) {
	cfg := testConsumptionalConfig()
	cfg.Fee.FixedFee = 0.01
	cfg.MaxSize = 100
	lab := newAnchorLab(t, protocol.ProtocolVersion5, cfg)
	repairFee := amount.Amount(1_000_000_000)
	addr, prv, acc := lab.fund(executor.MinAnchorDeposit + 1 + lab.fee()*2 + repairFee)
	lab.indexed[addr] = true
	first := lab.setTx(addr, prv, executor.MinAnchorDeposit, lab.fee(),
		anchorRoot(0xf5, 32), "first", 0, lab.height, "first")
	second := lab.setTx(addr, prv, 1, lab.fee(), anchorRoot(0xf6, 32), "second", 0, lab.height, "second")
	lab.pool.consumptionMap[addr] = cfg.Fee.DailyLimit - first.SerializeSize() - 1
	require.NoError(t, lab.pool.AppendTx(first))

	spent := lab.pool.consumptionMap[addr] + second.SerializeSize() + first.SerializeSize()
	err := lab.pool.AppendTx(second)
	require.ErrorIs(t, err, InvalidFeeError{MinimumFee: lab.fee() + consumptionFee(cfg, spent)})
	require.False(t, lab.pool.HasTx(second.ID()))
	require.True(t, lab.pool.HasTx(first.ID()))
	require.Equal(t, executor.MinAnchorDeposit+1, acc.LockedDeposit())
	require.Equal(t, repairFee, acc.Balance())

	transfer := lab.GenerateTestTransferTx(
		testsuite.TransactionWithLockTime(lab.height),
		testsuite.TransactionWithFee(lab.fee()),
		testsuite.TransactionWithSigner(prv),
	)
	transferSpent := lab.pool.consumptionMap[addr] + transfer.SerializeSize() + lab.pool.getPendingConsumption(addr)
	err = lab.pool.AppendTx(transfer)
	require.ErrorIs(t, err, InvalidFeeError{MinimumFee: lab.fee() + consumptionFee(cfg, transferSpent)})
	require.False(t, lab.pool.HasTx(transfer.ID()))
	require.Equal(t, lab.fee(), lab.pool.EstimatedFee(transfer.Payload().Value(), payload.TypeTransfer))
	require.Equal(t, repairFee, acc.Balance())

	repaired := lab.setTx(addr, prv, 0, repairFee, anchorRoot(0xf7, 32), "repair", 0, lab.height, "repair")
	require.NoError(t, lab.pool.AppendTx(repaired))
	require.True(t, lab.pool.HasTx(repaired.ID()))
	require.False(t, lab.pool.HasTx(second.ID()))
	require.Equal(t, executor.MinAnchorDeposit+1, acc.LockedDeposit())
	require.Equal(t, amount.Amount(0), acc.Balance())

	window := newAnchorLab(t, protocol.ProtocolVersion5, cfg)
	windowAddr, windowPrv, _ := window.fund(executor.MinAnchorDeposit + window.fee())
	window.indexed[windowAddr] = true
	counted := window.setTx(windowAddr, windowPrv, executor.MinAnchorDeposit, window.fee(),
		anchorRoot(0xf8, 32), "window", 0, window.height, "")
	require.NoError(t, window.pool.AppendTx(counted))
	fakeStore := window.pool.store.(*store.FakeStore)
	for height := types.Height(1); height <= types.Height(cfg.ConsumptionWindow)+1; height++ {
		txs := []*tx.Tx{window.makeSubsidy(height)}
		if height == 1 {
			txs = append(txs, counted)
		}
		blk, cert := window.GenerateTestBlock(height, testsuite.BlockWithTransactions(txs))
		fakeStore.SaveBlock(blk, cert)
		window.pool.HandleCommittedBlock(blk)
	}
	_, stillCounted := window.pool.consumptionMap[windowAddr]
	require.False(t, stillCounted)
}

func consumptionFee(cfg *Config, consumption int) amount.Amount {
	coefficient := consumption / cfg.Fee.DailyLimit
	fee, _ := amount.NewAmount(float64(coefficient) * float64(consumption) * cfg.Fee.UnitPrice)

	return fee
}

func (lab *anchorLab) makeSubsidy(lockTime types.Height) *tx.Tx {
	return lab.GenerateTestSubsidyTx(testsuite.TransactionWithLockTime(lockTime))
}

func (lab *anchorLab) recheck(version protocol.Version, height types.Height,
	accounts map[crypto.Address]*account.Account, banned bool,
) {
	lab.version = version
	lab.height = height
	lab.accounts = accounts
	if banned {
		for addr := range accounts {
			lab.banned[addr] = true
		}
	}
	lab.committed = map[tx.ID]struct{}{}
	lab.sbx = lab.bindSandbox()
	lab.pool.SetNewSandboxAndRecheck(lab.sbx)
}
