package state

import (
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/stretchr/testify/require"
)

// The tests in this file document issues that are not fixed yet, because fixing them
// changes the consensus rules or the behavior of the node, and needs a decision from the maintainers.
// Each test describes the expected behavior. Remove the `t.Skip` call once the issue is fixed.

// The subsidy transaction is a batch transfer, and a batch transfer is rejected if a recipient
// appears twice. `createSubsidyTx` doesn't merge the recipients, so a proposer whose reward address
// is the foundation address of the height, or whose delegation owner is its own reward address,
// builds a block that all the peers reject (and `BasicCheck` of its own proposal fails).
//
// A fix needs to merge the recipients in `createSubsidyTx` and to accept the merged subsidy in
// `checkSubsidy`. That makes valid some blocks that are invalid today, so it must be activated with
// a protocol version.
func TestKnownIssueSubsidyRecipientsMustBeUnique(t *testing.T) {
	t.Skip("known issue: the subsidy transaction has duplicated recipients if the reward address collides")

	td := setup(t)

	proposerAddr := td.state.Proposer(0).Address()
	foundationAddr := td.state.params.FoundationAddress(td.state.LastBlockHeight() + 1)

	tests := []struct {
		name       string
		rewardAddr crypto.Address
		dlgOwner   crypto.Address
		share      amount.Amount
	}{
		{"reward address is the foundation address", foundationAddr, crypto.TreasuryAddress, 0},
		{"delegation owner is the reward address", td.RandAccAddress(), crypto.TreasuryAddress, 0.2e9},
		{"delegation owner is the foundation address", td.RandAccAddress(), foundationAddr, 0.2e9},
		{"reward address is the foundation address, delegated", foundationAddr, td.RandAccAddress(), 0.2e9},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			owner := tt.dlgOwner
			if tt.name == "delegation owner is the reward address" {
				owner = tt.rewardAddr
			}

			val, err := td.state.ValidatorByAddress(proposerAddr)
			require.NoError(t, err)
			if tt.share > 0 {
				val.SetDelegation(owner, tt.share, td.RandHeight())
			} else {
				val.SetDelegation(crypto.TreasuryAddress, 0, 0)
			}
			td.state.store.UpdateValidator(val)

			trx := td.state.createSubsidyTx(proposerAddr, tt.rewardAddr, 0)

			require.NoError(t, trx.BasicCheck())
		})
	}
}
