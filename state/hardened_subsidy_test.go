package state

import (
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/state/param"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/protocol"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/stretchr/testify/require"
)

// The subsidy transaction should pay the recipients that `createSubsidyTx` makes.
// Until the hardened rules, only the total, the foundation, and a minimum for the owner were checked.
func TestExactSubsidy(t *testing.T) {
	const blockFee = amount.Amount(0.5e9)

	td := setup(t)
	proposerKey := td.proposerKey(t, 0)
	proposerAddr := proposerKey.Address()
	owner := td.RandAccAddress()
	reward := td.RandAccAddress()
	other := td.RandAccAddress()

	// A block that has a transaction with a fee.
	lockTime := td.state.LastBlockHeight()
	trx := tx.NewTransferTx(lockTime, td.genAccKey.PublicKeyNative().AccountAddress(), td.RandAccAddress(), 1e9, blockFee)
	td.HelperSignTransaction(td.genAccKey, trx)
	td.fakeTxPool.EXPECT().PrepareBlockTransactions().Return(block.Txs{trx}).AnyTimes()

	height := td.state.LastBlockHeight() + 1
	foundation := payload.BatchRecipient{
		To:     td.state.params.FoundationAddress(height),
		Amount: td.state.params.FoundationReward(height),
	}
	blockReward := td.state.params.BlockReward(height)
	maxShare := param.MaxDelegateOwnerRewardShare
	partialShare := amount.Amount(0.2e9)

	// The subsidy transaction is replaced in a block that was proposed with the right one.
	withSubsidy := func(recipients ...payload.BatchRecipient) *block.Block {
		base, err := td.state.ProposeBlock(proposerKey, reward)
		require.NoError(t, err)

		subsidy := tx.NewSubsidyTx(height, recipients)

		return block.MakeBlock(base.Header().Version(), base.Header().Time(),
			append(block.Txs{subsidy}, base.Transactions()[1:]...),
			base.Header().PrevBlockHash(), base.Header().StateRoot(), base.PrevCertificate(),
			base.Header().SortitionSeed(), base.Header().ProposerAddress())
	}

	execute := func(version protocol.Version, blk *block.Block) error {
		sbx := td.state.makeSandbox(version, blk.Header().UnixTime())

		return td.state.executeBlock(blk, sbx, true)
	}

	delegate := func(share amount.Amount) {
		val, err := td.state.ValidatorByAddress(proposerAddr)
		require.NoError(t, err)
		if share == 0 {
			val.SetDelegation(crypto.TreasuryAddress, 0, 0)
		} else {
			val.SetDelegation(owner, share, td.RandHeight())
		}
		td.state.store.UpdateValidator(val)
	}

	recipient := func(to crypto.Address, amt amount.Amount) payload.BatchRecipient {
		return payload.BatchRecipient{To: to, Amount: amt}
	}

	tests := []struct {
		name       string
		share      amount.Amount
		recipients []payload.BatchRecipient
		legacyOk   bool // accepted before the hardened rules
		hardenedOk bool // accepted by the hardened rules
		amountErr  bool // rejected because the total is wrong, before the hardened rules check the recipients
	}{
		{
			name: "no delegation, the whole reward and fee go to the reward address",
			recipients: []payload.BatchRecipient{
				foundation, recipient(reward, blockReward+blockFee-foundation.Amount),
			},
			legacyOk: true, hardenedOk: true,
		},
		{
			name: "no delegation, the reward is sent to another address",
			recipients: []payload.BatchRecipient{
				foundation, recipient(other, blockReward+blockFee-foundation.Amount),
			},
			legacyOk: true, hardenedOk: true,
		},
		{
			name: "no delegation, the reward is split",
			recipients: []payload.BatchRecipient{
				foundation, recipient(reward, 1), recipient(other, blockReward+blockFee-foundation.Amount-1),
			},
			legacyOk: false, hardenedOk: false,
		},
		{
			name: "no delegation, a part of the fee is missing",
			recipients: []payload.BatchRecipient{
				foundation, recipient(reward, blockReward+blockFee-foundation.Amount-1),
			},
			legacyOk: false, hardenedOk: false, amountErr: true,
		},
		{
			name: "partial share, the usual distribution", share: partialShare,
			recipients: []payload.BatchRecipient{
				foundation, recipient(owner, partialShare),
				recipient(reward, blockReward+blockFee-foundation.Amount-partialShare),
			},
			legacyOk: true, hardenedOk: true,
		},
		{
			name: "partial share, the owner gets more than the share", share: partialShare,
			recipients: []payload.BatchRecipient{
				foundation, recipient(owner, partialShare+1),
				recipient(reward, blockReward+blockFee-foundation.Amount-partialShare-1),
			},
			legacyOk: true, hardenedOk: false,
		},
		{
			name: "partial share, the owner gets all the remaining reward", share: partialShare,
			recipients: []payload.BatchRecipient{
				foundation, recipient(owner, blockReward+blockFee-foundation.Amount),
			},
			legacyOk: true, hardenedOk: false,
		},
		{
			name: "partial share, an extra recipient", share: partialShare,
			recipients: []payload.BatchRecipient{
				foundation, recipient(owner, partialShare),
				recipient(reward, blockReward+blockFee-foundation.Amount-partialShare-1), recipient(other, 1),
			},
			legacyOk: true, hardenedOk: false,
		},
		{
			name: "maximum share, the owner gets the fee", share: maxShare,
			recipients: []payload.BatchRecipient{
				foundation, recipient(owner, maxShare+blockFee),
			},
			legacyOk: true, hardenedOk: true,
		},
		{
			name: "maximum share, the proposer keeps the fee of the owner", share: maxShare,
			recipients: []payload.BatchRecipient{
				foundation, recipient(owner, maxShare), recipient(reward, blockFee),
			},
			legacyOk: true, hardenedOk: false,
		},
		{
			name: "maximum share, the proposer keeps a part of the fee of the owner", share: maxShare,
			recipients: []payload.BatchRecipient{
				foundation, recipient(owner, maxShare+blockFee-1), recipient(reward, 1),
			},
			legacyOk: true, hardenedOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			delegate(tt.share)
			blk := withSubsidy(tt.recipients...)

			err := execute(protocol.ProtocolVersion5, blk)
			if tt.legacyOk {
				require.NoError(t, err, "legacy rules")
			} else {
				require.Error(t, err, "legacy rules")
			}

			err = execute(protocol.ProtocolVersion6, blk)
			switch {
			case tt.hardenedOk:
				require.NoError(t, err, "hardened rules")
			case tt.amountErr:
				var amountErr InvalidSubsidyAmountError
				require.ErrorAs(t, err, &amountErr, "hardened rules")
			default:
				require.ErrorIs(t, err, ErrInvalidSubsidyTransaction, "hardened rules")
			}
		})
	}
}
