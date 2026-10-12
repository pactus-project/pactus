package state

import (
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/store"
	"github.com/pactus-project/pactus/types/account"
)

// AnchorAccount is an account that currently has an anchor.
type AnchorAccount struct {
	Address crypto.Address
	Account *account.Account
}

// listAnchors reads one page of anchor holders through the store's anchor index.
// An anchor deleted between the two reads is skipped.
func listAnchors(reader store.Reader, skip, count uint32) ([]AnchorAccount, uint32) {
	addrs, total := reader.AnchorAddresses(skip, count)
	items := make([]AnchorAccount, 0, len(addrs))
	for _, addr := range addrs {
		acc, err := reader.Account(addr)
		if err != nil || !acc.HasAnchor() {
			continue
		}
		items = append(items, AnchorAccount{Address: addr, Account: acc})
	}

	return items, total
}
