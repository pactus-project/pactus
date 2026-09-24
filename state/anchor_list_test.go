package state

import (
	"math"
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/store"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/require"
)

func TestListAnchorsOrdersByNumber(t *testing.T) {
	td := setup(t)
	fakeStore := td.state.store.(*store.FakeStore)

	plain := account.NewAccount(2)
	fakeStore.FakeAccounts[td.RandAccAddress()] = plain

	third := anchored(t, 30, 0x33)
	first := anchored(t, 10, 0x11)
	second := anchored(t, 20, 0x22)
	fakeStore.FakeAccounts[td.RandAccAddress()] = third
	fakeStore.FakeAccounts[td.RandAccAddress()] = first
	fakeStore.FakeAccounts[td.RandAccAddress()] = second

	items, total := td.state.ListAnchors(0, 20)
	require.Equal(t, uint32(3), total)
	require.Equal(t, []int32{10, 20, 30}, numbers(items))

	page, total := td.state.ListAnchors(2, 50)
	require.Equal(t, uint32(3), total)
	require.Equal(t, []int32{30}, numbers(page))

	empty, total := td.state.ListAnchors(math.MaxUint32, 1)
	require.Equal(t, uint32(3), total)
	require.Empty(t, empty)

	none, total := td.state.ListAnchors(0, 0)
	require.Equal(t, uint32(3), total)
	require.Empty(t, none)
}

func TestListAnchorsSignedOrderAndMutation(t *testing.T) {
	td := setup(t)
	fakeStore := td.state.store.(*store.FakeStore)

	low := anchored(t, -1, 0x01)
	mid := anchored(t, 0, 0x02)
	high := anchored(t, math.MaxInt32, 0x03)
	// Address bytes run the opposite way from the account numbers.
	fakeStore.FakeAccounts[filledAddress(0x00)] = high
	fakeStore.FakeAccounts[filledAddress(0x80)] = mid
	fakeStore.FakeAccounts[filledAddress(0xFF)] = low
	for number := int32(1); number <= 5; number++ {
		fakeStore.FakeAccounts[td.RandAccAddress()] = account.NewAccount(number)
	}

	items, total := td.state.ListAnchors(0, 20)
	require.Equal(t, uint32(3), total)
	require.Equal(t, []int32{-1, 0, math.MaxInt32}, numbers(items))

	kept := items[0].Address
	fakeStore.FakeAccounts[items[1].Address].ClearAnchor()
	items, total = td.state.ListAnchors(0, 20)
	require.Equal(t, uint32(2), total)
	require.Equal(t, kept, items[0].Address)
	require.NotEqual(t, items[0].Address, items[1].Address)

	require.NoError(t, fakeStore.FakeAccounts[items[0].Address].SetAnchor(account.AnchorData{
		RootHash:      bytesRepeat(0x44, 32),
		LockedDeposit: 1,
	}))
	again, total := td.state.ListAnchors(0, 20)
	require.Equal(t, uint32(2), total)
	require.Equal(t, items[0].Address, again[0].Address)
	require.Equal(t, bytesRepeat(0x44, 32), again[0].Account.RootHash())
}

func anchored(t *testing.T, number int32, mark byte) *account.Account {
	t.Helper()

	acc := account.NewAccount(number)
	require.NoError(t, acc.SetAnchor(account.AnchorData{
		RootHash:      bytesRepeat(mark, 32),
		LockedDeposit: amount.Amount(1),
	}))

	return acc
}

func numbers(items []AnchorAccount) []int32 {
	out := make([]int32, 0, len(items))
	for _, item := range items {
		out = append(out, item.Account.Number())
	}

	return out
}

func TestListAnchorsStopsAtTheEnd(t *testing.T) {
	td := setup(t)
	fakeStore := td.state.store.(*store.FakeStore)
	for number := int32(0); number < 5; number++ {
		fakeStore.FakeAccounts[td.RandAccAddress()] = anchored(t, number, byte(number))
	}

	got, total := td.state.ListAnchors(2, 2)
	require.Equal(t, uint32(5), total)
	require.Equal(t, []int32{2, 3}, numbers(got))

	got, total = td.state.ListAnchors(2, 50)
	require.Equal(t, uint32(5), total)
	require.Equal(t, []int32{2, 3, 4}, numbers(got))

	got, total = td.state.ListAnchors(2, math.MaxUint32)
	require.Equal(t, uint32(5), total)
	require.Equal(t, []int32{2, 3, 4}, numbers(got))
}

// An anchor can be deleted between reading the index and reading the account.
// Such accounts are skipped, and the total still comes from the index.
func TestListAnchorsSkipsChangedAccounts(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	reader := store.NewMockReader(ts.MockController())
	kept, cleared, gone := ts.RandAccAddress(), ts.RandAccAddress(), ts.RandAccAddress()

	reader.EXPECT().AnchorAddresses(uint32(0), uint32(3)).
		Return([]crypto.Address{kept, cleared, gone}, uint32(3))
	reader.EXPECT().Account(kept).Return(anchored(t, 1, 0x01), nil)
	reader.EXPECT().Account(cleared).Return(account.NewAccount(2), nil)
	reader.EXPECT().Account(gone).Return(nil, store.ErrNotFound)

	items, total := listAnchors(reader, 0, 3)
	require.Equal(t, uint32(3), total)
	require.Len(t, items, 1)
	require.Equal(t, kept, items[0].Address)
}

func filledAddress(fill byte) crypto.Address {
	data := bytesRepeat(fill, 20)

	return crypto.NewAddress(crypto.AddressTypeBLSAccount, data)
}

func bytesRepeat(fill byte, n int) []byte {
	out := make([]byte, n)
	for i := range out {
		out[i] = fill
	}

	return out
}
