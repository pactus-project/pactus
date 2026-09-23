package state

import (
	"cmp"
	"slices"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/account"
)

// AnchorAccount is an account that currently has an anchor.
type AnchorAccount struct {
	Address crypto.Address
	Account *account.Account
}

func gatherAnchors(iter func(func(crypto.Address, *account.Account) bool)) []AnchorAccount {
	items := make([]AnchorAccount, 0)
	iter(func(addr crypto.Address, acc *account.Account) bool {
		if acc.HasAnchor() {
			items = append(items, AnchorAccount{Address: addr, Account: acc})
		}

		return false
	})
	slices.SortFunc(items, func(left, right AnchorAccount) int {
		return cmp.Compare(left.Account.Number(), right.Account.Number())
	})

	return items
}

func pageAnchors(items []AnchorAccount, skip, count uint32) ([]AnchorAccount, uint32) {
	total := uint32(len(items))
	if skip >= total || count == 0 {
		return nil, total
	}

	end := skip + count
	if end < skip || end > total {
		end = total
	}

	return items[skip:end], total
}
