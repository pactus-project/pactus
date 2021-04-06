package state

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/genesis"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/param"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/tx/payload"
	"github.com/zarbchain/zarb-go/txpool"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/validator"
)

var tState1 *state
var tState2 *state
var tState3 *state
var tState4 *state
var tValSigner1 crypto.Signer
var tValSigner2 crypto.Signer
var tValSigner3 crypto.Signer
var tValSigner4 crypto.Signer
var tGenTime time.Time
var tCommonTxPool *txpool.MockTxPool

func setup(t *testing.T) {
	if tState1 != nil {
		tState1.Close()
		tState2.Close()
		tState3.Close()
		tState4.Close()
	}
	logger.InitLogger(logger.TestConfig())

	_, _, priv1 := crypto.GenerateTestKeyPair()
	_, _, priv2 := crypto.GenerateTestKeyPair()
	_, _, priv3 := crypto.GenerateTestKeyPair()
	_, _, priv4 := crypto.GenerateTestKeyPair()

	tValSigner1 = crypto.NewSigner(priv1)
	tValSigner2 = crypto.NewSigner(priv2)
	tValSigner3 = crypto.NewSigner(priv3)
	tValSigner4 = crypto.NewSigner(priv4)

	tGenTime = util.RoundNow(10)
	tCommonTxPool = txpool.MockingTxPool()

	acc := account.NewAccount(crypto.TreasuryAddress, 0)
	acc.AddToBalance(21 * 1e14) // 2,100,000,000,000,000
	val1 := validator.NewValidator(tValSigner1.PublicKey(), 0, 0)
	val2 := validator.NewValidator(tValSigner2.PublicKey(), 1, 0)
	val3 := validator.NewValidator(tValSigner3.PublicKey(), 2, 0)
	val4 := validator.NewValidator(tValSigner4.PublicKey(), 3, 0)
	params := param.DefaultParams()
	params.CommitteeSize = 4
	gnDoc := genesis.MakeGenesis(tGenTime, []*account.Account{acc}, []*validator.Validator{val1, val2, val3, val4}, params)

	st1, err := LoadOrNewState(TestConfig(), gnDoc, tValSigner1, tCommonTxPool)
	require.NoError(t, err)
	st2, err := LoadOrNewState(TestConfig(), gnDoc, tValSigner2, tCommonTxPool)
	require.NoError(t, err)
	st3, err := LoadOrNewState(TestConfig(), gnDoc, tValSigner3, tCommonTxPool)
	require.NoError(t, err)
	st4, err := LoadOrNewState(TestConfig(), gnDoc, tValSigner4, tCommonTxPool)
	require.NoError(t, err)

	tState1, _ = st1.(*state)
	tState2, _ = st2.(*state)
	tState3, _ = st3.(*state)
	tState4, _ = st4.(*state)
}

func makeBlockAndCertificate(t *testing.T, round int, signers ...crypto.Signer) (block.Block, block.Certificate) {
	var st *state
	if tState1.committee.IsProposer(tState1.signer.Address(), round) {
		st = tState1
	} else if tState1.committee.IsProposer(tState2.signer.Address(), round) {
		st = tState2
	} else if tState1.committee.IsProposer(tState3.signer.Address(), round) {
		st = tState3
	} else {
		st = tState4
	}

	b, err := st.ProposeBlock(round)
	require.NoError(t, err)
	c := makeCertificateAndSign(t, b.Hash(), round, signers...)

	return *b, c
}

func makeCertificateAndSign(t *testing.T, blockHash crypto.Hash, round int, signers ...crypto.Signer) block.Certificate {
	sigs := make([]crypto.Signature, len(signers))
	sb := block.CertificateSignBytes(blockHash, round)
	committers := []int{0, 1, 2, 3}
	signedBy := []int{}

	for i, s := range signers {
		if s.Address().EqualsTo(tValSigner1.Address()) {
			signedBy = append(signedBy, 0)
		}

		if s.Address().EqualsTo(tValSigner2.Address()) {
			signedBy = append(signedBy, 1)
		}

		if s.Address().EqualsTo(tValSigner3.Address()) {
			signedBy = append(signedBy, 2)
		}

		if s.Address().EqualsTo(tValSigner4.Address()) {
			signedBy = append(signedBy, 3)
		}
		sigs[i] = s.SignData(sb)
	}

	absences := util.Subtracts(committers, signedBy)
	return *block.NewCertificate(blockHash, round, committers, absences, crypto.Aggregate(sigs))
}

func CommitBlockForAllStates(t *testing.T, b block.Block, c block.Certificate) {
	assert.NoError(t, tState1.CommitBlock(tState1.lastInfo.BlockHeight()+1, b, c))
	assert.NoError(t, tState2.CommitBlock(tState2.lastInfo.BlockHeight()+1, b, c))
	assert.NoError(t, tState3.CommitBlock(tState3.lastInfo.BlockHeight()+1, b, c))
	assert.NoError(t, tState4.CommitBlock(tState4.lastInfo.BlockHeight()+1, b, c))
}

func moveToNextHeightForAllStates(t *testing.T) {
	b, c := makeBlockAndCertificate(t, 0, tValSigner1, tValSigner2, tValSigner3, tValSigner4)
	CommitBlockForAllStates(t, b, c)
}

func TestProposeBlockAndValidation(t *testing.T) {
	setup(t)
	moveToNextHeightForAllStates(t)

	b, err := tState1.ProposeBlock(0)
	assert.Error(t, err)
	assert.Nil(t, b)

	trx := tx.NewSendTx(crypto.UndefHash, 1, tValSigner1.Address(), tValSigner2.Address(), 1000, 1000, "")
	tValSigner1.SignMsg(trx)
	assert.NoError(t, tCommonTxPool.AppendTx(trx))

	b, err = tState2.ProposeBlock(0)
	assert.NoError(t, err)
	assert.NotNil(t, b)
	assert.Equal(t, b.TxIDs().Len(), 2)

	err = tState1.ValidateBlock(*b)
	require.NoError(t, err)

	// Propose and validate again
	b, err = tState2.ProposeBlock(0)
	assert.NoError(t, err)
	assert.NotNil(t, b)
	assert.Equal(t, b.TxIDs().Len(), 2)

	err = tState1.ValidateBlock(*b)
	require.NoError(t, err)
}

func TestBlockSubsidy(t *testing.T) {
	interval := 2100000
	assert.Equal(t, int64(5*1e8), calcBlockSubsidy(1, interval))
	assert.Equal(t, int64(5*1e8), calcBlockSubsidy((1*interval)-1, interval))
	assert.Equal(t, int64(2.5*1e8), calcBlockSubsidy((1*interval), interval))
	assert.Equal(t, int64(2.5*1e8), calcBlockSubsidy((2*interval)-1, interval))
	assert.Equal(t, int64(1.25*1e8), calcBlockSubsidy((2*interval), interval))
}

func TestBlockSubsidyTx(t *testing.T) {
	setup(t)

	// Without mintbase address in config
	trx := tState1.createSubsidyTx(7)
	assert.True(t, trx.IsMintbaseTx())
	assert.Equal(t, trx.Payload().Value(), calcBlockSubsidy(1, tState1.params.SubsidyReductionInterval)+7)
	assert.Equal(t, trx.Payload().(*payload.SendPayload).Receiver, tValSigner1.Address())

	// With ivalid mintbase address in config
	tState1.config.MintbaseAddress = "invalid"
	tState1.Close()
	_, err := LoadOrNewState(tState1.config, tState1.genDoc, tValSigner1, tCommonTxPool)
	assert.Error(t, err)

	// With mintbase address in config
	addr, _, _ := crypto.GenerateTestKeyPair()
	tState1.config.MintbaseAddress = addr.String()
	tState1.Close()
	st, err := LoadOrNewState(tState1.config, tState1.genDoc, tValSigner1, tCommonTxPool)
	assert.NoError(t, err)
	trx = st.(*state).createSubsidyTx(0)
	assert.Equal(t, trx.Payload().(*payload.SendPayload).Receiver, addr)
}

func TestCommitBlocks(t *testing.T) {
	setup(t)

	b1, c1 := makeBlockAndCertificate(t, 1, tValSigner1, tValSigner2, tValSigner3)
	invBlock, _ := block.GenerateTestBlock(nil, nil)
	assert.Error(t, tState1.CommitBlock(1, *invBlock, c1))
	// No error here but block is ignored, because the height is invalid
	assert.NoError(t, tState1.CommitBlock(2, b1, c1))
	assert.NoError(t, tState1.CommitBlock(1, b1, c1))

	assert.Equal(t, tState1.LastBlockHash(), b1.Hash())
	assert.Equal(t, tState1.LastBlockTime(), b1.Header().Time())
	assert.Equal(t, tState1.LastCertificate().Hash(), c1.Hash())
	assert.Equal(t, tState1.LastBlockHeight(), 1)
	assert.Equal(t, tState1.GenesisHash(), tState2.GenesisHash())
}

func TestCommitSandbox(t *testing.T) {

	t.Run("Certificate new account", func(t *testing.T) {
		setup(t)

		addr, _, _ := crypto.GenerateTestKeyPair()
		sb := tState1.makeSandbox()
		newAcc := sb.MakeNewAccount(addr)
		newAcc.AddToBalance(1)
		tState1.commitSandbox(sb, 0)

		assert.True(t, tState1.store.HasAccount(addr))
	})

	t.Run("Certificate new validator", func(t *testing.T) {
		setup(t)

		addr, pub, _ := crypto.GenerateTestKeyPair()
		sb := tState1.makeSandbox()
		newVal := sb.MakeNewValidator(pub)
		newVal.AddToStake(1)
		sb.UpdateValidator(newVal)
		tState1.commitSandbox(sb, 0)

		assert.True(t, tState1.store.HasValidator(addr))
		assert.Equal(t, sb.TotalStakeChange(), int64(1))
		assert.Equal(t, tState1.sortition.TotalStake(), int64(1))
	})

	t.Run("Modify account", func(t *testing.T) {
		setup(t)

		sb := tState1.makeSandbox()
		acc := sb.Account(crypto.TreasuryAddress)
		acc.SubtractFromBalance(1)
		sb.UpdateAccount(acc)
		tState1.commitSandbox(sb, 0)

		acc1, _ := tState1.store.Account(crypto.TreasuryAddress)
		assert.Equal(t, acc1.Balance(), acc.Balance())
	})

	t.Run("Modify validator", func(t *testing.T) {
		setup(t)

		sb := tState1.makeSandbox()
		val := sb.Validator(tValSigner2.Address())
		val.AddToStake(2)
		sb.UpdateValidator(val)
		tState1.commitSandbox(sb, 0)

		val1, _ := tState1.store.Validator(tValSigner2.Address())
		assert.Equal(t, val1.Stake(), val.Stake())
		assert.Equal(t, sb.TotalStakeChange(), int64(2))
	})

	t.Run("Move committee", func(t *testing.T) {
		setup(t)

		nextProposer := tState1.committee.Proposer(1)

		sb := tState1.makeSandbox()
		tState1.commitSandbox(sb, 0)

		assert.Equal(t, tState1.committee.Proposer(0).Address(), nextProposer.Address())
	})

	t.Run("Move committee next round", func(t *testing.T) {
		setup(t)

		nextNextProposer := tState1.committee.Proposer(2)

		sb := tState1.makeSandbox()
		tState1.commitSandbox(sb, 1)

		assert.Equal(t, tState1.committee.Proposer(0).Address(), nextNextProposer.Address())
	})
}

func TestUpdateLastCertificate(t *testing.T) {
	setup(t)
	b1, c1 := makeBlockAndCertificate(t, 0, tValSigner1, tValSigner3, tValSigner4)
	b11, c11 := makeBlockAndCertificate(t, 0, tValSigner1, tValSigner2, tValSigner3, tValSigner4)
	_, c12 := makeBlockAndCertificate(t, 1, tValSigner1, tValSigner2, tValSigner3, tValSigner4)

	CommitBlockForAllStates(t, b1, c1)

	assert.Equal(t, b1.Hash(), b11.Hash())
	assert.Equal(t, tState1.lastInfo.Certificate().Hash(), c1.Hash())
	assert.Error(t, tState1.UpdateLastCertificate(&c12))
	assert.NoError(t, tState1.UpdateLastCertificate(&c1))
	assert.Equal(t, tState1.lastInfo.Certificate().Hash(), c1.Hash())
	assert.NoError(t, tState1.UpdateLastCertificate(&c11))
	assert.Equal(t, tState1.lastInfo.Certificate().Hash(), c11.Hash())
}

func TestInvalidProposerProposeBlock(t *testing.T) {
	setup(t)

	_, err := tState2.ProposeBlock(0)
	assert.Error(t, err)
	_, err = tState2.ProposeBlock(1)
	assert.NoError(t, err)
}

func TestBlockProposal(t *testing.T) {
	setup(t)
	moveToNextHeightForAllStates(t)

	t.Run("validity of proposed block", func(t *testing.T) {
		b, err := tState2.ProposeBlock(0)
		assert.NoError(t, err)
		assert.NoError(t, tState1.ValidateBlock(*b)) // State1 check state2's proposed block
	})

	t.Run("Tx pool has two subsidy transactions", func(t *testing.T) {
		trx := tState2.createSubsidyTx(0)
		assert.NoError(t, tState2.txPool.AppendTx(trx))

		// Moving to the next round
		b, err := tState3.ProposeBlock(1)
		assert.NoError(t, err)
		assert.NoError(t, tState1.ValidateBlock(*b))
	})
}

func TestInvalidBlock(t *testing.T) {
	setup(t)

	b, _ := block.GenerateTestBlock(nil, nil)
	assert.Error(t, tState1.ValidateBlock(*b))
}

func TestForkDetection(t *testing.T) {
	setup(t)

	b1, c1 := makeBlockAndCertificate(t, 0, tValSigner1, tValSigner2, tValSigner3)
	b2, c2 := makeBlockAndCertificate(t, 1, tValSigner1, tValSigner2, tValSigner3)
	assert.NoError(t, tState1.CommitBlock(1, b1, c1))
	assert.NoError(t, tState1.CommitBlock(1, b1, c1))
	assert.Error(t, tState1.CommitBlock(1, b2, c2))
}

func TestNodeShutdown(t *testing.T) {
	setup(t)
	b1, c1 := makeBlockAndCertificate(t, 0, tValSigner1, tValSigner2, tValSigner3)

	// Should not panic or crash
	tState1.Close()
	assert.Error(t, tState1.CommitBlock(1, b1, c1))
	b, _ := block.GenerateTestBlock(nil, nil)
	assert.Error(t, tState1.ValidateBlock(*b))
	_, err := tState1.ProposeBlock(0)
	assert.Error(t, err)
}

func TestSortition(t *testing.T) {
	setup(t)

	addr, pub, priv := crypto.GenerateTestKeyPair()
	signer := crypto.NewSigner(priv)

	st, err := LoadOrNewState(TestConfig(), tState1.genDoc, signer, tCommonTxPool)
	assert.NoError(t, err)
	st1 := st.(*state)

	assert.False(t, st1.evaluateSortition()) //  not a validator

	// Certificate 14 blocks
	height := 1
	for ; height < 12; height++ {
		if height == 4 {
			trx := tx.NewBondTx(crypto.UndefHash, 1, tValSigner1.Address(), pub, 1000, 1000, "")
			tValSigner1.SignMsg(trx)
			assert.NoError(t, tCommonTxPool.AppendTx(trx))
		}

		b, c := makeBlockAndCertificate(t, 0, tValSigner1, tValSigner2, tValSigner3, tValSigner4)
		CommitBlockForAllStates(t, b, c)
		require.NoError(t, st1.CommitBlock(height, b, c))
	}

	assert.False(t, st1.evaluateSortition()) //  bonding period

	// Certificate another block
	b, c := makeBlockAndCertificate(t, 0, tValSigner1, tValSigner2, tValSigner3, tValSigner4)
	CommitBlockForAllStates(t, b, c)
	require.NoError(t, st1.CommitBlock(height, b, c))
	height++

	assert.True(t, st1.evaluateSortition())           //  ok
	assert.False(t, tState1.committee.Contains(addr)) // still not in the set

	// ---------------------------------------------
	// Certificate another block, new validator should be in the set now
	b, c = makeBlockAndCertificate(t, 0, tValSigner1, tValSigner2, tValSigner3, tValSigner4)
	CommitBlockForAllStates(t, b, c)
	require.NoError(t, st1.CommitBlock(height, b, c))

	assert.False(t, st1.evaluateSortition()) // already in the set
	assert.False(t, tState1.committee.Contains(tValSigner1.Address()))
	assert.True(t, tState1.committee.Contains(addr))

	// ---------------------------------------------
	// Let's save and load tState1
	committeeHash := tState1.committee.CommitteeHash()
	tState1.Close()
	state1, _ := LoadOrNewState(tState1.config, tState1.genDoc, tValSigner1, tCommonTxPool)

	assert.Equal(t, state1.(*state).committee.CommitteeHash(), committeeHash)

	// ---------------------------------------------
	// Let's commit another block with new Validator set
	b1, err := st1.ProposeBlock(3)
	require.NoError(t, err)
	require.NotNil(t, b1)

	sigs := make([]crypto.Signature, 4)
	sb := block.CertificateSignBytes(b1.Hash(), 3)

	sigs[0] = tValSigner2.SignData(sb)
	sigs[1] = tValSigner3.SignData(sb)
	sigs[2] = tValSigner4.SignData(sb)
	sigs[3] = signer.SignData(sb)
	c1 := block.NewCertificate(b1.Hash(), 3, []int{4, 1, 2, 3}, []int{}, crypto.Aggregate(sigs))

	height++
	require.NoError(t, st1.CommitBlock(height, *b1, *c1))
	require.NoError(t, tState2.CommitBlock(height, *b1, *c1))
}

func TestValidateBlockTime(t *testing.T) {
	setup(t)

	fmt.Printf("BlockTimeInSecond: %d\n", tState1.params.BlockTimeInSecond)
	roundedNow := util.RoundNow(10)
	tState1.lastInfo.SetBlockTime(roundedNow.Add(-1 * time.Minute))

	// Time not rounded
	assert.Error(t, tState1.validateBlockTime(roundedNow.Add(-15*time.Second)))
	assert.Error(t, tState1.validateBlockTime(roundedNow.Add(-5*time.Second)))
	assert.Error(t, tState1.validateBlockTime(roundedNow.Add(5*time.Second)))
	assert.Error(t, tState1.validateBlockTime(roundedNow.Add(15*time.Second)))

	// Too early
	assert.Error(t, tState1.validateBlockTime(tState1.lastInfo.BlockTime().Add(-20*time.Second)))
	assert.Error(t, tState1.validateBlockTime(tState1.lastInfo.BlockTime().Add(-10*time.Second)))
	assert.Error(t, tState1.validateBlockTime(tState1.lastInfo.BlockTime()))

	// Ok
	assert.NoError(t, tState1.validateBlockTime(roundedNow.Add(10*time.Second)))
	assert.NoError(t, tState1.validateBlockTime(roundedNow.Add(20*time.Second)))

	// Too late
	assert.Error(t, tState1.validateBlockTime(roundedNow.Add(30*time.Second)))
	assert.Error(t, tState1.validateBlockTime(roundedNow.Add(40*time.Second)))
}

func TestInvalidBlockVersion(t *testing.T) {
	setup(t)

	tState1.params.BlockVersion = 2
	b, _ := tState1.ProposeBlock(0)
	assert.Error(t, tState2.ValidateBlock(*b))
}

func TestInvalidBlockTime(t *testing.T) {
	setup(t)
	moveToNextHeightForAllStates(t)

	validBlock, _ := makeBlockAndCertificate(t, 0, tValSigner1, tValSigner2, tValSigner3, tValSigner4)

	invalidBlock := block.MakeBlock(
		validBlock.Header().Version(),
		validBlock.Header().Time().Add(30*time.Second),
		validBlock.TxIDs(),
		validBlock.Header().LastBlockHash(),
		validBlock.Header().CommitteeHash(),
		validBlock.Header().StateHash(),
		validBlock.Header().LastReceiptsHash(),
		validBlock.LastCertificate(),
		validBlock.Header().SortitionSeed(),
		validBlock.Header().ProposerAddress())

	assert.NoError(t, tState1.ValidateBlock(validBlock))
	assert.Error(t, tState1.ValidateBlock(invalidBlock))

}

func TestValidatorHelpers(t *testing.T) {
	setup(t)

	t.Run("Should return nil for NonExisting Validator Address", func(t *testing.T) {
		_, _, priv5 := crypto.GenerateTestKeyPair()
		nonExistenceValidator := tState1.Validator(priv5.PublicKey().Address())
		assert.Nil(t, nonExistenceValidator, "State 1 returned Non nil For nonExisting validator")
		nonExistenceValidator = tState2.Validator(priv5.PublicKey().Address())
		assert.Nil(t, nonExistenceValidator, "State 2 returned Non nil For nonExisting validator")
		nonExistenceValidator = tState3.Validator(priv5.PublicKey().Address())
		assert.Nil(t, nonExistenceValidator, "State 3 returned Non nil For nonExisting validator")
		nonExistenceValidator = tState4.Validator(priv5.PublicKey().Address())
		assert.Nil(t, nonExistenceValidator, "State 4 returned Non nil For nonExisting validator")
	})

	t.Run("Should return validator for valid committee Validator Address", func(t *testing.T) {
		existingValidator := tState4.Validator(tValSigner1.Address())
		assert.NotNil(t, existingValidator)
		assert.Equal(t, 0, existingValidator.Number())
	})

	t.Run("Should return validator for corresponding Validator number", func(t *testing.T) {
		existingValidator := tState4.ValidatorByNumber(1)
		assert.NotNil(t, existingValidator)
		assert.Equal(t, tValSigner2.Address(), existingValidator.Address())
	})

	t.Run("Should return nil for invalid Validator number", func(t *testing.T) {
		nonExistenceValidator := tState4.ValidatorByNumber(10)
		assert.Nil(t, nonExistenceValidator)
	})
}
