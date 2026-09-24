package executor

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/require"
)

// clampAmount maps any int64 into [0, MaxNanoPAC], the range a stored
// balance or deposit can hold.
func clampAmount(value int64) amount.Amount {
	if value < 0 {
		value = -(value + 1)
	}

	return amount.Amount(value % (amount.MaxNanoPAC + 1))
}

// FuzzAnchorCheckExecute guards the rule the commit path relies on: once
// Check accepts an anchor, Execute must not panic, must conserve coins
// (minus the fee) and must leave an account that can be stored.
// A rejected anchor must leave the account untouched.
func FuzzAnchorCheckExecute(f *testing.F) {
	minDeposit := int64(MinAnchorDeposit)
	maxDeposit := int64(amount.MaxNanoPAC)

	add := func(action, hashLen, uriLen, anchorType uint8,
		deposit, fee, balance, locked int64, height, unixTime uint32,
	) {
		f.Add(action, hashLen, uriLen, anchorType, deposit, fee, balance, locked, height, unixTime)
	}

	// action, hashLen, uriLen, anchorType, deposit, fee, balance, locked, height, unixTime
	add(0, 32, 3, 0, minDeposit, 1, 10*minDeposit, 0, 10, 100)
	add(0, 64, 128, 0xFF, 0, 0, minDeposit, minDeposit, 20, 200)
	add(0, 32, 0, 1, 0, 1, 0, 0, 1, 1)
	add(0, 32, 0, 1, maxDeposit, 1, maxDeposit, 0, 1, 1)
	add(0, 32, 0, 1, 1, 0, minDeposit, maxDeposit, 1, 1)
	add(1, 0, 0, 0, 0, 1, 0, minDeposit, 5, 50)
	add(1, 0, 0, 0, 0, minDeposit+1, 0, minDeposit, 5, 50)
	add(1, 0, 0, 0, 0, 0, maxDeposit, minDeposit, 5, 50)
	add(1, 0, 0, 0, 0, -1, 0, minDeposit, 5, 50)
	add(2, 32, 0, 0, minDeposit, 1, 10*minDeposit, 0, 5, 50)
	add(0, 31, 129, 0, minDeposit, 1, 10*minDeposit, 0, 5, 50)

	f.Fuzz(func(t *testing.T, action, hashLen, uriLen, anchorType uint8,
		deposit, fee, balance, locked int64, height, unixTime uint32,
	) {
		td := setup(t)
		td.useAnchor(unixTime)
		td.sbx.FakeHeight = types.Height(height)

		acc, addr := td.addTestAccount(t, testsuite.AccountWithBalance(clampAmount(balance)))
		if locked > 0 {
			require.NoError(t, acc.SetAnchor(account.AnchorData{
				RootHash:      bytes.Repeat([]byte{0x5A}, 32),
				LockedDeposit: amount.Amount((locked-1)%amount.MaxNanoPAC + 1), // in [1, MaxNanoPAC]
			}))
		}

		uri := strings.Repeat("u", int(uriLen))
		if uriLen > 0 && anchorType%2 == 1 {
			uri = uri[:len(uri)-1] + string([]byte{0xFF}) // invalid UTF-8
		}
		trx := tx.NewAnchorTx(td.sbx.CurrentHeight(), addr, action,
			bytes.Repeat([]byte{0xAB}, int(hashLen)), uri, anchorType,
			amount.Amount(deposit), amount.Amount(fee))

		oldLocked := acc.LockedDeposit()
		oldTotal := acc.Balance() + oldLocked
		oldCreatedH, oldCreatedT := acc.CreatedAtHeight(), acc.CreatedAtTime()
		hadAnchor := acc.HasAnchor()
		before := accountBytes(t, acc)

		exe, err := MakeExecutor(trx, td.sbx)
		require.NoError(t, err)
		if err := exe.Check(td.sbx, true); err != nil {
			requireUnchanged(t, td.sbx.Account(addr), before)

			return
		}

		require.NotPanics(t, func() {
			exe.Execute(td.sbx)
		})

		got := td.sbx.Account(addr)
		require.Equal(t, oldTotal-amount.Amount(fee), got.Balance()+got.LockedDeposit())

		raw := accountBytes(t, got)
		decoded, err := account.FromBytes(raw)
		require.NoError(t, err)
		require.Equal(t, got.Hash(), decoded.Hash())
		require.Equal(t, got.HasAnchor(), len(raw) > 12)

		switch action {
		case payload.AnchorActionSet:
			require.True(t, got.HasAnchor())
			require.Equal(t, oldLocked+amount.Amount(deposit), got.LockedDeposit())
			require.Equal(t, types.Height(height), got.UpdatedAtHeight())
			require.Equal(t, unixTime, got.UpdatedAtTime())
			if hadAnchor {
				// An existing anchor stays valid below the minimum (PIP-50 section 3).
				require.Equal(t, oldCreatedH, got.CreatedAtHeight())
				require.Equal(t, oldCreatedT, got.CreatedAtTime())
			} else {
				require.GreaterOrEqual(t, got.LockedDeposit(), MinAnchorDeposit)
				require.Equal(t, types.Height(height), got.CreatedAtHeight())
				require.Equal(t, unixTime, got.CreatedAtTime())
			}
		case payload.AnchorActionDelete:
			require.False(t, got.HasAnchor())
			require.Len(t, raw, 12)
		default:
			require.Failf(t, "invalid action accepted", "action %d", action)
		}
	})
}
