package executor

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/require"
)

// This file checks the anchor executor against an independent model written
// from the PIP-50 text (sections 5 and 6), over random operation sequences.

const maxNano = amount.Amount(amount.MaxNanoPAC)

// errBasicCheck stands for any payload.BasicCheckError.
var errBasicCheck = errors.New("basic check error")

type modelAnchor struct {
	root     []byte
	uri      string
	kind     uint8
	locked   amount.Amount
	createdH types.Height
	createdT uint32
	updatedH types.Height
	updatedT uint32
}

type modelAccount struct {
	balance amount.Amount
	anchor  *modelAnchor
}

type anchorModel struct {
	td       *testData
	addrs    []crypto.Address
	accounts map[crypto.Address]*modelAccount
	fees     amount.Amount
	supply   amount.Amount
}

// modelOp is one random operation with the result the PIP predicts.
type modelOp struct {
	kind    string
	trx     *tx.Tx
	wantErr error
	apply   func()
}

func newAnchorModel(t *testing.T, td *testData) *anchorModel {
	t.Helper()

	model := &anchorModel{
		td:       td,
		accounts: make(map[crypto.Address]*modelAccount),
	}
	starts := []amount.Amount{
		0,
		MinAnchorDeposit / 2,
		MinAnchorDeposit,
		3*MinAnchorDeposit + td.RandAmount(),
		50 * MinAnchorDeposit,
		maxNano - MinAnchorDeposit,
	}
	for i := 0; i < 4; i++ {
		balance := starts[td.RandIntMax(len(starts))]
		_, addr := td.addTestAccount(t, testsuite.AccountWithBalance(balance))
		model.addrs = append(model.addrs, addr)
		model.accounts[addr] = &modelAccount{balance: balance}
		model.supply += balance
	}

	return model
}

func (m *anchorModel) pickAmount(options ...amount.Amount) amount.Amount {
	return options[m.td.RandIntMax(len(options))]
}

func (m *anchorModel) randomFee() amount.Amount {
	return m.pickAmount(0, 0, 1, 1, amount.Amount(m.td.RandInt64Max(1e7)), -1, maxNano, maxNano+1)
}

func (m *anchorModel) randomRoot() []byte {
	lengths := []int{0, 31, 32, 32, 33, 48, 64, 64, 65}

	return m.td.RandBytes(lengths[m.td.RandIntMax(len(lengths))])
}

func (m *anchorModel) randomURI() string {
	uris := []string{
		"",
		"ipfs://manifest",
		strings.Repeat("é", 64),             // 128 bytes, valid
		strings.Repeat("a", 129),            // one byte too long
		string([]byte{0xC3}),                // truncated UTF-8
		"https://x/" + string([]byte{0xFF}), // invalid UTF-8
	}

	return uris[m.td.RandIntMax(len(uris))]
}

// isAnchorOwner follows PIP section 4: only user account addresses own a slot.
func isAnchorOwner(addr crypto.Address) bool {
	switch addr.Type() {
	case crypto.AddressTypeBLSAccount, crypto.AddressTypeEd25519Account, crypto.AddressTypeSecp256k1Account:
		return true
	case crypto.AddressTypeTreasury, crypto.AddressTypeValidator:
		return false
	default:
		return false
	}
}

func (m *anchorModel) signer() crypto.Address {
	switch m.td.RandIntMax(20) {
	case 0:
		return crypto.TreasuryAddress
	case 1:
		return m.td.RandAccAddress() // unknown account
	default:
		return m.addrs[m.td.RandIntMax(len(m.addrs))]
	}
}

func (m *anchorModel) nextOp() modelOp {
	switch m.td.RandIntMax(5) {
	case 0, 1:
		return m.setOp()
	case 2:
		return m.deleteOp()
	default:
		return m.transferOp()
	}
}

// setOp follows PIP sections 5 and 6.1.
func (m *anchorModel) setOp() modelOp {
	from := m.signer()
	root, uri := m.randomRoot(), m.randomURI()
	kind := uint8(m.td.RandIntMax(256))
	fee := m.randomFee()
	switch m.td.RandIntMax(3) {
	case 0:
		// Well formed, so the random part exercises deposits, balances and
		// anchor state rather than BasicCheck.
		root = m.td.RandBytes(32 + m.td.RandIntMax(33))
		uri = "ipfs://manifest"
		fee = amount.Amount(m.td.RandInt64Max(1e7))
	case 1:
		// Keep the current content, as a deposit top-up does.
		if acc := m.accounts[from]; acc != nil && acc.anchor != nil {
			root = bytes.Clone(acc.anchor.root)
			uri = acc.anchor.uri
			kind = acc.anchor.kind
			fee = amount.Amount(m.td.RandInt64Max(1e7))
		}
	}

	var balance amount.Amount
	acc := m.accounts[from]
	if acc != nil {
		balance = acc.balance
	}
	deposit := m.pickAmount(-1, 0, 1, MinAnchorDeposit-1, MinAnchorDeposit, MinAnchorDeposit,
		MinAnchorDeposit+amount.Amount(m.td.RandInt64Max(int64(10*MinAnchorDeposit))),
		balance, maxNano, maxNano+1)

	trx := tx.NewAnchorTx(m.td.sbx.CurrentHeight(), from, payload.AnchorActionSet, root, uri, kind, deposit, fee)
	result := modelOp{kind: "update", trx: trx}
	switch {
	case acc == nil || acc.anchor == nil:
		result.kind = "create"
	case bytes.Equal(acc.anchor.root, root) && acc.anchor.uri == uri && acc.anchor.kind == kind:
		result.kind = "keep content"
	}

	switch {
	case acc == nil && from != crypto.TreasuryAddress:
		result.wantErr = AccountNotFoundError{Address: from}
	case !isAnchorOwner(from),
		len(root) < 32, len(root) > 64,
		len(uri) > 128, !utf8.ValidString(uri),
		deposit < 0, deposit > maxNano:
		result.wantErr = errBasicCheck
	case fee < 0, fee > maxNano:
		result.wantErr = ErrAmountOverflow
	case deposit > maxNano-fee:
		result.wantErr = ErrAmountOverflow
	case balance < deposit+fee:
		result.wantErr = ErrInsufficientFunds
	case acc.anchor != nil && acc.anchor.locked > maxNano-deposit:
		result.wantErr = ErrAmountOverflow
	case acc.anchor == nil && deposit < MinAnchorDeposit:
		result.wantErr = ErrAnchorDepositTooSmall
	default:
		height, unixTime := m.td.sbx.CurrentHeight(), m.td.sbx.FakeUnixTime
		result.apply = func() {
			acc.balance -= deposit + fee
			m.fees += fee
			// PIP-50 section 6.1: the content date moves only when the
			// root hash, the manifest URI or the anchor type changes.
			contentChanged := acc.anchor == nil ||
				!bytes.Equal(acc.anchor.root, root) || acc.anchor.uri != uri || acc.anchor.kind != kind
			if acc.anchor == nil {
				acc.anchor = &modelAnchor{createdH: height, createdT: unixTime}
			}
			acc.anchor.root = bytes.Clone(root)
			acc.anchor.uri = uri
			acc.anchor.kind = kind
			acc.anchor.locked += deposit
			if contentChanged {
				acc.anchor.updatedH = height
				acc.anchor.updatedT = unixTime
			}
		}
	}

	return result
}

// deleteOp follows PIP sections 5 and 6.2.
func (m *anchorModel) deleteOp() modelOp {
	from := m.signer()
	fee := m.randomFee()
	trx := tx.NewAnchorTx(m.td.sbx.CurrentHeight(), from, payload.AnchorActionDelete, nil, "", 0, 0, fee)
	result := modelOp{kind: "delete", trx: trx}

	acc := m.accounts[from]
	switch {
	case acc == nil && from != crypto.TreasuryAddress:
		result.wantErr = AccountNotFoundError{Address: from}
	case !isAnchorOwner(from):
		result.wantErr = errBasicCheck
	case fee < 0, fee > maxNano:
		result.wantErr = ErrAmountOverflow
	case acc.anchor == nil:
		result.wantErr = ErrAnchorNotFound
	case acc.anchor.locked < fee:
		result.wantErr = ErrInsufficientFunds
	case acc.balance > maxNano-(acc.anchor.locked-fee):
		result.wantErr = ErrAmountOverflow
	default:
		result.apply = func() {
			acc.balance += acc.anchor.locked - fee
			m.fees += fee
			acc.anchor = nil
		}
	}

	return result
}

// transferOp mixes plain transfers in; they must never touch the locked deposit.
func (m *anchorModel) transferOp() modelOp {
	from := m.addrs[m.td.RandIntMax(len(m.addrs))]
	receiverAddr := m.addrs[m.td.RandIntMax(len(m.addrs))]
	sender, receiver := m.accounts[from], m.accounts[receiverAddr]
	fee := m.pickAmount(0, 1, amount.Amount(m.td.RandInt64Max(1e7)))
	amt := m.pickAmount(0, 1, sender.balance, sender.balance+1,
		amount.Amount(m.td.RandInt64Max(int64(sender.balance)+1)))

	trx := tx.NewTransferTx(m.td.sbx.CurrentHeight(), from, receiverAddr, amt, fee)
	result := modelOp{kind: "transfer", trx: trx}
	if sender.balance < amt+fee {
		result.wantErr = ErrInsufficientFunds
	} else {
		result.apply = func() {
			sender.balance -= amt + fee
			receiver.balance += amt
			m.fees += fee
		}
	}

	return result
}

func requireCheckResult(t *testing.T, got, want error) {
	t.Helper()

	switch {
	case want == nil:
		require.NoError(t, got)
	case errors.Is(want, errBasicCheck):
		var basicErr payload.BasicCheckError
		require.ErrorAs(t, got, &basicErr)
	default:
		require.ErrorIs(t, got, want)
	}
}

// requireMatchesModel compares every account with the model and checks the
// invariants that must hold after any accepted or rejected operation.
func (m *anchorModel) requireMatchesModel(t *testing.T) {
	t.Helper()

	total := m.fees
	for _, addr := range m.addrs {
		want := m.accounts[addr]
		got := m.td.sbx.Account(addr)
		total += got.Balance() + got.LockedDeposit()

		require.Equal(t, want.balance, got.Balance())
		require.Equal(t, want.anchor != nil, got.HasAnchor())

		raw, err := got.Bytes()
		require.NoError(t, err)
		decoded, err := account.FromBytes(raw)
		require.NoError(t, err, "a stored account must decode")
		require.Equal(t, got.Hash(), decoded.Hash())

		if want.anchor == nil {
			require.Len(t, raw, 12, "an account without anchor keeps the 12-byte record")

			continue
		}
		require.Greater(t, len(raw), 12)
		require.Equal(t, want.anchor.root, got.RootHash())
		require.Equal(t, want.anchor.uri, got.ManifestURI())
		require.Equal(t, want.anchor.kind, got.AnchorType())
		require.Equal(t, want.anchor.locked, got.LockedDeposit())
		require.Equal(t, want.anchor.createdH, got.CreatedAtHeight())
		require.Equal(t, want.anchor.createdT, got.CreatedAtTime())
		require.Equal(t, want.anchor.updatedH, got.UpdatedAtHeight())
		require.Equal(t, want.anchor.updatedT, got.UpdatedAtTime())

		require.GreaterOrEqual(t, got.LockedDeposit(), MinAnchorDeposit)
		require.LessOrEqual(t, got.CreatedAtHeight(), got.UpdatedAtHeight())
		require.LessOrEqual(t, got.CreatedAtTime(), got.UpdatedAtTime())
	}
	require.Equal(t, m.supply, total, "anchors and transfers must not mint or burn coins")
}

func TestAnchorModelRandomSequences(t *testing.T) {
	sequences, steps := 300, 40
	if testing.Short() {
		// Long enough sequences to chain create, update and delete.
		sequences, steps = 30, 60
	}

	td := setup(t)
	td.useAnchor(1_700_000_000)
	_, _ = td.addTestAccount(t, testsuite.AccountWithAddress(crypto.TreasuryAddress),
		testsuite.AccountWithBalance(10*MinAnchorDeposit))

	accepted := map[string]int{}
	rejected := 0
	for seq := 0; seq < sequences; seq++ {
		model := newAnchorModel(t, td)
		for step := 0; step < steps; step++ {
			// Several operations can share a block; clocks never go back.
			td.sbx.FakeHeight += types.Height(td.RandIntMax(3))
			td.sbx.FakeUnixTime += uint32(td.RandIntMax(21))

			result := model.nextOp()
			exe, err := MakeExecutor(result.trx, td.sbx)
			if err == nil {
				err = exe.Check(td.sbx, true)
			}
			requireCheckResult(t, err, result.wantErr)

			if result.wantErr == nil {
				exe.Execute(td.sbx)
				result.apply()
				accepted[result.kind]++
			} else {
				rejected++
			}
			model.requireMatchesModel(t)
		}
	}

	for _, kind := range []string{"create", "update", "keep content", "delete", "transfer"} {
		require.Positive(t, accepted[kind], "no accepted %s", kind)
	}
	require.Positive(t, rejected)
	t.Logf("accepted: %v, rejected: %d", accepted, rejected)
}
