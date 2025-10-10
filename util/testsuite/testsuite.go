package testsuite

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/committee"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/ed25519"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/sortition"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/protocol"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util"
	"golang.org/x/exp/slices"
)

// TestSuite provides a set of helper functions for testing purposes.
// All the random values are generated based on a logged seed.
// By using a pre-generated seed, it is possible to reproduce failed tests
// by re-evaluating all the random values. This helps in identifying and debugging
// failures in testing conditions.
type TestSuite struct {
	Seed int64
	Rand *rand.Rand
}

func GenerateSeed() int64 {
	return time.Now().UTC().UnixNano()
}

// NewTestSuiteFromSeed creates a new TestSuite with the given seed.
func NewTestSuiteFromSeed(seed int64) *TestSuite {
	return &TestSuite{
		Seed: seed,
		//nolint:gosec // to reproduce the failed tests
		Rand: rand.New(rand.NewSource(seed)),
	}
}

// NewTestSuite creates a new TestSuite by generating new seed.
func NewTestSuite(t *testing.T) *TestSuite {
	t.Helper()

	seed := GenerateSeed()
	t.Logf("%v seed is %v", t.Name(), seed)

	return &TestSuite{
		Seed: seed,
		//nolint:gosec // to reproduce the failed tests
		Rand: rand.New(rand.NewSource(seed)),
	}
}

// RandBool returns a random boolean value.
func (ts *TestSuite) RandBool() bool {
	return ts.RandInt64(2) == 0
}

// RandInt8 returns a random int8 between 0 and max: [0, max).
func (ts *TestSuite) RandInt8(max int8) int8 {
	return int8(ts.RandUint64(uint64(max)))
}

// RandUint8 returns a random uint8 between 0 and max: [0, max).
func (ts *TestSuite) RandUint8(max uint8) uint8 {
	return uint8(ts.RandUint64(uint64(max)))
}

// RandInt16 returns a random int16 between 0 and max: [0, max).
func (ts *TestSuite) RandInt16(max int16) int16 {
	return int16(ts.RandUint64(uint64(max)))
}

// RandUint16 returns a random uint16 between 0 and max: [0, max).
func (ts *TestSuite) RandUint16(max uint16) uint16 {
	return uint16(ts.RandUint64(uint64(max)))
}

// RandInt32 returns a random int32 between 0 and max: [0, max).
func (ts *TestSuite) RandInt32(max int32) int32 {
	return int32(ts.RandUint64(uint64(max)))
}

// RandUint32 returns a random uint32 between 0 and max: [0, max).
func (ts *TestSuite) RandUint32(max uint32) uint32 {
	return uint32(ts.RandUint64(uint64(max)))
}

// RandInt64 returns a random int64 between 0 and max: [0, max).
func (ts *TestSuite) RandInt64(max int64) int64 {
	return ts.Rand.Int63n(max)
}

// RandUint64 returns a random uint64 between 0 and max: [0, max).
func (ts *TestSuite) RandUint64(max uint64) uint64 {
	return uint64(ts.RandInt64(int64(max)))
}

// RandInt returns a random int between 0 and max: [0, max).
func (ts *TestSuite) RandInt(max int) int {
	return int(ts.RandInt64(int64(max)))
}

// RandInt16NonZero returns a random int16 between 1 and max+1: [1, max+1).
func (ts *TestSuite) RandInt16NonZero(max int16) int16 {
	return ts.RandInt16(max) + 1
}

// RandUint16NonZero returns a random uint16 between 1 and max+1: [1, max+1).
func (ts *TestSuite) RandUint16NonZero(max uint16) uint16 {
	return ts.RandUint16(max) + 1
}

// RandInt32NonZero returns a random int32 between 1 and max+1: [1, max+1).
func (ts *TestSuite) RandInt32NonZero(max int32) int32 {
	return ts.RandInt32(max) + 1
}

// RandUint32NonZero returns a random uint32 between 1 and max+1: [1, max+1).
func (ts *TestSuite) RandUint32NonZero(max uint32) uint32 {
	return ts.RandUint32(max) + 1
}

// RandInt64NonZero returns a random int64 between 1 and max+1: [1, max+1).
func (ts *TestSuite) RandInt64NonZero(max int64) int64 {
	return ts.RandInt64(max) + 1
}

// RandUint64NonZero returns a random uint64 between 1 and max+1: [1, max+1).
func (ts *TestSuite) RandUint64NonZero(max uint64) uint64 {
	return ts.RandUint64(max) + 1
}

// RandIntNonZero returns a random int between 1 and max+1: [1, max+1).
func (ts *TestSuite) RandIntNonZero(max int) int {
	return ts.RandInt(max) + 1
}

// RandHeight returns a random number between [1000, 1000000] for block height.
func (ts *TestSuite) RandHeight() uint32 {
	return ts.RandUint32NonZero(1e6-1000) + 1000
}

// RandRound returns a random number between [0, 10) for block round.
func (ts *TestSuite) RandRound() int16 {
	return ts.RandInt16(10)
}

// RandAmount returns a random amount between [1e9, max).
// If max is not set, it defaults to 100e9.
func (ts *TestSuite) RandAmount(max ...amount.Amount) amount.Amount {
	maxAmt := amount.Amount(100e9) // default max amount
	if len(max) > 0 {
		maxAmt = max[0]
	}

	return ts.RandAmountRange(1e9, maxAmt)
}

// RandAmountRange returns a random amount between [min, max).
func (ts *TestSuite) RandAmountRange(min, max amount.Amount) amount.Amount {
	amt := amount.Amount(ts.RandInt64NonZero(int64(max - min)))

	return amt + min
}

// RandFee returns a random fee between [1e7, max).
// If max is not set, it defaults to 1e9.
func (ts *TestSuite) RandFee(max ...amount.Amount) amount.Amount {
	maxFee := amount.Amount(1e9) // default max fee
	if len(max) > 0 {
		maxFee = max[0]
	}

	return ts.RandAmountRange(1e7, maxFee)
}

// RandBytes returns a slice of random bytes of the given length.
func (ts *TestSuite) RandBytes(length int) []byte {
	buf := make([]byte, length)
	_, _ = ts.Rand.Read(buf)

	return buf
}

// RandSlice generates a random non-repeating slice of int32 elements with the specified length.
func (ts *TestSuite) RandSlice(length int) []int32 {
	slice := []int32{}
	for {
		randInt := ts.RandInt32(1000)
		if !slices.Contains(slice, randInt) {
			slice = append(slice, randInt)
		}

		if len(slice) == length {
			return slice
		}
	}
}

// RandString generates a random string of the given length.
func (ts *TestSuite) RandString(length int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[ts.RandInt(52)]
	}

	return string(b)
}

// DecodingHex decodes the input string from hexadecimal format and returns the resulting byte slice.
func (*TestSuite) DecodingHex(in string) []byte {
	d, _ := hex.DecodeString(in)

	return d
}

func (ts *TestSuite) RandKeyPair() (crypto.PublicKey, crypto.PrivateKey) {
	switch ts.RandInt(2) {
	case 0:
		return ts.RandBLSKeyPair()
	case 1:
		return ts.RandEd25519KeyPair()
	}

	// Impossible
	return nil, nil
}

// RandBLSKeyPair generates a random BLS key pair for testing purposes.
func (ts *TestSuite) RandBLSKeyPair() (*bls.PublicKey, *bls.PrivateKey) {
	buf := ts.RandBytes(bls.PrivateKeySize)
	prv, _ := bls.PrivateKeyFromBytes(buf)
	pub := prv.PublicKeyNative()

	return pub, prv
}

// RandEd25519KeyPair generates a random Ed25519 key pair for testing purposes.
func (ts *TestSuite) RandEd25519KeyPair() (*ed25519.PublicKey, *ed25519.PrivateKey) {
	buf := ts.RandBytes(ed25519.PrivateKeySize)
	prv, _ := ed25519.PrivateKeyFromBytes(buf)
	pub := prv.PublicKeyNative()

	return pub, prv
}

// RandValKey generates a random validator key for testing purposes.
func (ts *TestSuite) RandValKey() *bls.ValidatorKey {
	_, prv := ts.RandBLSKeyPair()

	return bls.NewValidatorKey(prv)
}

// RandBLSSignature generates a random BLS signature for testing purposes.
func (ts *TestSuite) RandBLSSignature() *bls.Signature {
	_, prv := ts.RandBLSKeyPair()
	sig := prv.Sign(ts.RandBytes(8))

	return sig.(*bls.Signature)
}

// RandEd25519Signature generates a random BLS signature for testing purposes.
func (ts *TestSuite) RandEd25519Signature() *ed25519.Signature {
	_, prv := ts.RandEd25519KeyPair()
	sig := prv.Sign(ts.RandBytes(8))

	return sig.(*ed25519.Signature)
}

// RandHash generates a random hash for testing purposes.
func (ts *TestSuite) RandHash() hash.Hash {
	return hash.CalcHash(util.Int64ToSlice(ts.RandInt64(util.MaxInt64)))
}

// RandAccAddress generates a random account address for testing purposes.
func (ts *TestSuite) RandAccAddress() crypto.Address {
	useBLSAddress := ts.RandBool()
	if useBLSAddress {
		return crypto.NewAddress(crypto.AddressTypeBLSAccount, ts.RandBytes(20))
	}

	return crypto.NewAddress(crypto.AddressTypeEd25519Account, ts.RandBytes(20))
}

// RandValAddress generates a random validator address for testing purposes.
func (ts *TestSuite) RandValAddress() crypto.Address {
	addr := crypto.NewAddress(crypto.AddressTypeValidator, ts.RandBytes(20))

	return addr
}

// RandSeed generates a random VerifiableSeed for testing purposes.
func (ts *TestSuite) RandSeed() sortition.VerifiableSeed {
	sig := ts.RandBLSSignature()
	seed, _ := sortition.VerifiableSeedFromBytes(sig.Bytes())

	return seed
}

// RandProof generates a random Proof for testing purposes.
func (ts *TestSuite) RandProof() sortition.Proof {
	_, prv := ts.RandBLSKeyPair()
	sig := prv.Sign(ts.RandHash().Bytes())
	proof, _ := sortition.ProofFromBytes(sig.Bytes())

	return proof
}

// RandPeerID returns a random peer ID.
func (ts *TestSuite) RandPeerID() peer.ID {
	s := ts.RandBytes(32)
	id := [34]byte{0x12, 32}
	copy(id[2:], s)

	return peer.ID(id[:])
}

// RandMultiAddress returns a random MultiAddress.
func (ts *TestSuite) RandMultiAddress() string {
	return fmt.Sprintf("/dns/%s/udp/1234", ts.RandString(12))
}

type AccountMaker struct {
	Number  int32
	Balance amount.Amount
	Address crypto.Address
}

// NewAccountMaker creates a new instance of AccountMaker with random values.
func (ts *TestSuite) NewAccountMaker() *AccountMaker {
	return &AccountMaker{
		Number:  ts.RandInt32NonZero(100000),
		Balance: ts.RandAmountRange(100e9, 1000e9),
		Address: ts.RandAccAddress(),
	}
}

// AccountWithAddress sets the address for the generated test account.
func AccountWithAddress(address crypto.Address) func(*AccountMaker) {
	return func(am *AccountMaker) {
		am.Address = address
	}
}

// AccountWithNumber sets the account number for the generated test account.
func AccountWithNumber(number int32) func(*AccountMaker) {
	return func(am *AccountMaker) {
		am.Number = number
	}
}

// AccountWithBalance sets the balance for the generated test account.
func AccountWithBalance(balance amount.Amount) func(*AccountMaker) {
	return func(am *AccountMaker) {
		am.Balance = balance
	}
}

// GenerateTestAccount generates an account for testing purposes.
func (ts *TestSuite) GenerateTestAccount(options ...func(*AccountMaker)) (*account.Account, crypto.Address) {
	amk := ts.NewAccountMaker()
	for _, opt := range options {
		opt(amk)
	}
	acc := account.NewAccount(amk.Number)
	acc.AddToBalance(amk.Balance)

	return acc, amk.Address
}

type ValidatorMaker struct {
	Number    int32
	Stake     amount.Amount
	PublicKey *bls.PublicKey
}

// NewValidatorMaker creates a new instance of ValidatorMaker with random values.
func (ts *TestSuite) NewValidatorMaker() *ValidatorMaker {
	return &ValidatorMaker{
		Number:    ts.RandInt32(100000),
		Stake:     ts.RandAmountRange(100e9, 1000e9),
		PublicKey: ts.RandValKey().PublicKey(),
	}
}

// ValidatorWithNumber sets the validator number for the generated test validator.
func ValidatorWithNumber(number int32) func(*ValidatorMaker) {
	return func(vm *ValidatorMaker) {
		vm.Number = number
	}
}

// ValidatorWithStake sets the stake for the generated test account.
func ValidatorWithStake(stake amount.Amount) func(*ValidatorMaker) {
	return func(vm *ValidatorMaker) {
		vm.Stake = stake
	}
}

// ValidatorWithPublicKey sets the public Key for the generated test account.
func ValidatorWithPublicKey(publicKey *bls.PublicKey) func(*ValidatorMaker) {
	return func(vm *ValidatorMaker) {
		vm.PublicKey = publicKey
	}
}

// GenerateTestValidator generates a validator for testing purposes.
func (ts *TestSuite) GenerateTestValidator(options ...func(*ValidatorMaker)) *validator.Validator {
	vmk := ts.NewValidatorMaker()
	for _, opt := range options {
		opt(vmk)
	}

	val := validator.NewValidator(vmk.PublicKey, vmk.Number)
	val.AddToStake(vmk.Stake)

	return val
}

type BlockMaker struct {
	Version   protocol.Version
	Txs       block.Txs
	Proposer  crypto.Address
	Time      time.Time
	StateHash hash.Hash
	PrevHash  hash.Hash
	Seed      sortition.VerifiableSeed
	PrevCert  *certificate.Certificate
}

// NewBlockMaker creates a new BlockMaker instance.
func (ts *TestSuite) NewBlockMaker() *BlockMaker {
	txs := block.NewTxs()
	tx0 := ts.GenerateTestSubsidyTx()
	tx1 := ts.GenerateTestTransferTx()
	tx2 := ts.GenerateTestSortitionTx()
	tx3 := ts.GenerateTestBondTx()
	tx4 := ts.GenerateTestUnbondTx()
	tx5 := ts.GenerateTestWithdrawTx()

	txs.Append(tx0)
	txs.Append(tx1)
	txs.Append(tx2)
	txs.Append(tx3)
	txs.Append(tx4)
	txs.Append(tx5)

	return &BlockMaker{
		Version:  protocol.ProtocolVersion2,
		Txs:      txs,
		Proposer: ts.RandValAddress(),
		Time:     time.Now(),
		PrevHash: ts.RandHash(),
		Seed:     ts.RandSeed(),
		PrevCert: nil,
	}
}

// BlockWithVersion sets version to the block.
func BlockWithVersion(ver protocol.Version) func(*BlockMaker) {
	return func(bm *BlockMaker) {
		bm.Version = ver
	}
}

// BlockWithProposer sets proposer address to the block.
func BlockWithProposer(addr crypto.Address) func(*BlockMaker) {
	return func(bm *BlockMaker) {
		bm.Proposer = addr
	}
}

// BlockWithTime sets block creation time to the block.
func BlockWithTime(t time.Time) func(*BlockMaker) {
	return func(bm *BlockMaker) {
		bm.Time = t
	}
}

// BlockWithStateHash sets state hash to the block.
func BlockWithStateHash(h hash.Hash) func(*BlockMaker) {
	return func(bm *BlockMaker) {
		bm.StateHash = h
	}
}

// BlockWithPrevHash sets previous block hash to the block.
func BlockWithPrevHash(h hash.Hash) func(*BlockMaker) {
	return func(bm *BlockMaker) {
		bm.PrevHash = h
	}
}

// BlockWithSeed sets verifiable seed to the block.
func BlockWithSeed(seed sortition.VerifiableSeed) func(*BlockMaker) {
	return func(bm *BlockMaker) {
		bm.Seed = seed
	}
}

// BlockWithPrevCert sets previous block certificate to the block.
func BlockWithPrevCert(cert *certificate.Certificate) func(*BlockMaker) {
	return func(bm *BlockMaker) {
		bm.PrevCert = cert
	}
}

// BlockWithTransactions adds transactions to the block.
func BlockWithTransactions(txs block.Txs) func(*BlockMaker) {
	return func(bm *BlockMaker) {
		bm.Txs = txs
	}
}

// GenerateTestBlock generates a block for testing purposes with optional configuration.
func (ts *TestSuite) GenerateTestBlock(height uint32, options ...func(*BlockMaker)) (
	*block.Block, *certificate.Certificate,
) {
	bmk := ts.NewBlockMaker()
	bmk.PrevCert = ts.GenerateTestCertificate(height - 1)

	if height == 1 {
		bmk.PrevCert = nil
		bmk.PrevHash = hash.UndefHash
	}

	for _, opt := range options {
		opt(bmk)
	}

	header := block.NewHeader(bmk.Version, bmk.Time, bmk.PrevHash, bmk.PrevHash, bmk.Seed, bmk.Proposer)
	blk := block.NewBlock(header, bmk.PrevCert, bmk.Txs)

	blockCert := ts.GenerateTestCertificate(height)

	return blk, blockCert
}

// GenerateTestCertificate generates a block certificate for testing purposes.
func (ts *TestSuite) GenerateTestCertificate(height uint32) *certificate.Certificate {
	sig := ts.RandBLSSignature()

	cert := certificate.NewCertificate(height, ts.RandRound())

	committers := ts.RandSlice(4)
	absentees := []int32{committers[3]}
	cert.SetSignature(committers, absentees, sig)

	return cert
}

type ProposalMaker struct {
	ProposerKey *bls.ValidatorKey
}

// NewProposalMaker creates a new NewProposalMaker instance.
func (ts *TestSuite) NewProposalMaker() *ProposalMaker {
	return &ProposalMaker{
		ProposerKey: ts.RandValKey(),
	}
}

// ProposalWithKey sets the private key of the proposer.
func ProposalWithKey(key *bls.ValidatorKey) func(*ProposalMaker) {
	return func(pm *ProposalMaker) {
		pm.ProposerKey = key
	}
}

// GenerateTestProposal generates a proposal for testing purposes.
func (ts *TestSuite) GenerateTestProposal(height uint32, round int16,
	options ...func(*ProposalMaker),
) *proposal.Proposal {
	pmk := ts.NewProposalMaker()

	for _, opt := range options {
		opt(pmk)
	}

	blk, _ := ts.GenerateTestBlock(height, BlockWithProposer(pmk.ProposerKey.Address()))
	prop := proposal.NewProposal(height, round, blk)
	ts.HelperSignProposal(pmk.ProposerKey, prop)

	return prop
}

type TransactionMaker struct {
	LockTime   uint32
	Amount     amount.Amount
	Fee        amount.Amount
	Signer     crypto.PrivateKey
	ValPubKey  *bls.PublicKey
	Recipients []payload.BatchRecipient
}

func (tm *TransactionMaker) SignerAccountAddress() crypto.Address {
	blsPub, ok := tm.Signer.PublicKey().(*bls.PublicKey)
	if ok {
		return blsPub.AccountAddress()
	}
	ed25519Pub := tm.Signer.PublicKey().(*ed25519.PublicKey)

	return ed25519Pub.AccountAddress()
}

func (tm *TransactionMaker) SignerValidatorAddress() crypto.Address {
	blsPub := tm.Signer.PublicKey().(*bls.PublicKey)

	return blsPub.ValidatorAddress()
}

// NewTransactionMaker creates a new TransactionMaker instance.
func (ts *TestSuite) NewTransactionMaker() *TransactionMaker {
	numOfRecipients := ts.RandInt(6) + 2
	recipients := make([]payload.BatchRecipient, numOfRecipients)

	for i := 0; i < numOfRecipients; i++ {
		recipients[i] = payload.BatchRecipient{
			To:     ts.RandAccAddress(),
			Amount: ts.RandAmount(10e9),
		}
	}

	return &TransactionMaker{
		LockTime:   ts.RandHeight(),
		Amount:     ts.RandAmount(),
		Fee:        ts.RandFee(),
		Signer:     nil,
		ValPubKey:  nil,
		Recipients: recipients,
	}
}

// TransactionWithLockTime sets lock-time to the transaction.
func TransactionWithLockTime(lockTime uint32) func(*TransactionMaker) {
	return func(tm *TransactionMaker) {
		tm.LockTime = lockTime
	}
}

// TransactionWithAmount sets amount to the transaction.
func TransactionWithAmount(amt amount.Amount) func(*TransactionMaker) {
	return func(tm *TransactionMaker) {
		tm.Amount = amt
	}
}

// TransactionWithFee sets fee to the transaction.
func TransactionWithFee(fee amount.Amount) func(*TransactionMaker) {
	return func(tm *TransactionMaker) {
		tm.Fee = fee
	}
}

// TransactionWithSigner sets the BLS signer to sign the test transaction.
func TransactionWithSigner(signer crypto.PrivateKey) func(*TransactionMaker) {
	return func(tm *TransactionMaker) {
		tm.Signer = signer
	}
}

// TransactionWithValidatorPublicKey sets the Validator's public key for the Bond transaction.
func TransactionWithValidatorPublicKey(pubKey *bls.PublicKey) func(*TransactionMaker) {
	return func(tm *TransactionMaker) {
		tm.ValPubKey = pubKey
	}
}

// TransactionWithRecipients sets the recipients for the Bath Transfer transaction.
func TransactionWithRecipients(recipients []payload.BatchRecipient) func(*TransactionMaker) {
	return func(tm *TransactionMaker) {
		tm.Recipients = recipients
	}
}

// GenerateTestTransferTx generates a transfer transaction for testing purposes.
func (ts *TestSuite) GenerateTestTransferTx(options ...func(*TransactionMaker)) *tx.Tx {
	tmk := ts.NewTransactionMaker()

	for _, opt := range options {
		opt(tmk)
	}

	if tmk.Signer == nil {
		useBLSSigner := ts.RandBool()
		if useBLSSigner {
			_, prv := ts.RandBLSKeyPair()
			tmk.Signer = prv
		} else {
			_, prv := ts.RandEd25519KeyPair()
			tmk.Signer = prv
		}
	}

	sender := tmk.SignerAccountAddress()
	trx := tx.NewTransferTx(tmk.LockTime, sender, ts.RandAccAddress(), tmk.Amount, tmk.Fee)
	ts.HelperSignTransaction(tmk.Signer, trx)

	return trx
}

// GenerateTestBatchTransferTx generate a batch transfer transaction for test.
func (ts *TestSuite) GenerateTestBatchTransferTx(options ...func(*TransactionMaker)) *tx.Tx {
	tmk := ts.NewTransactionMaker()

	for _, opt := range options {
		opt(tmk)
	}

	if tmk.Signer == nil {
		useBLSSigner := ts.RandBool()
		if useBLSSigner {
			_, prv := ts.RandBLSKeyPair()
			tmk.Signer = prv
		} else {
			_, prv := ts.RandEd25519KeyPair()
			tmk.Signer = prv
		}
	}

	numOfRecip := ts.RandInt(6) + 2
	recipients := make([]payload.BatchRecipient, numOfRecip)

	for i := 0; i < numOfRecip; i++ {
		recipients[i] = payload.BatchRecipient{
			To:     ts.RandAccAddress(),
			Amount: ts.RandAmount(10e9),
		}
	}

	trx := tx.NewBatchTransferTx(tmk.LockTime, tmk.SignerAccountAddress(), tmk.Recipients, tmk.Fee)
	ts.HelperSignTransaction(tmk.Signer, trx)

	return trx
}

// GenerateTestSubsidyTx creates a subsidy transaction for testing.
func (ts *TestSuite) GenerateTestSubsidyTx(options ...func(tm *TransactionMaker)) *tx.Tx {
	tmk := ts.NewTransactionMaker()
	for _, opt := range options {
		opt(tmk)
	}
	if tmk.LockTime == 0 {
		tmk.LockTime = ts.RandHeight()
	}

	return tx.NewSubsidyTx(tmk.LockTime, tmk.Recipients)
}

// GenerateTestBondTx generates a bond transaction for testing purposes.
func (ts *TestSuite) GenerateTestBondTx(options ...func(*TransactionMaker)) *tx.Tx {
	tmk := ts.NewTransactionMaker()

	for _, opt := range options {
		opt(tmk)
	}

	if tmk.Signer == nil {
		useBLSSigner := ts.RandBool()
		if useBLSSigner {
			_, prv := ts.RandBLSKeyPair()
			tmk.Signer = prv
		} else {
			_, prv := ts.RandEd25519KeyPair()
			tmk.Signer = prv
		}
	}

	sender := tmk.SignerAccountAddress()
	receiver := ts.RandValAddress()
	if tmk.ValPubKey != nil {
		receiver = tmk.ValPubKey.ValidatorAddress()
	}
	trx := tx.NewBondTx(tmk.LockTime, sender, receiver, tmk.ValPubKey, tmk.Amount, tmk.Fee)
	ts.HelperSignTransaction(tmk.Signer, trx)

	return trx
}

// GenerateTestSortitionTx generates a sortition transaction for testing purposes.
func (ts *TestSuite) GenerateTestSortitionTx(options ...func(*TransactionMaker)) *tx.Tx {
	tmk := ts.NewTransactionMaker()

	for _, opt := range options {
		opt(tmk)
	}

	if tmk.Signer == nil {
		_, prv := ts.RandBLSKeyPair()
		tmk.Signer = prv
	}

	proof := ts.RandProof()
	sender := tmk.SignerValidatorAddress()
	trx := tx.NewSortitionTx(tmk.LockTime, sender, proof)
	ts.HelperSignTransaction(tmk.Signer, trx)

	return trx
}

// GenerateTestUnbondTx generates an unbond transaction for testing purposes.
func (ts *TestSuite) GenerateTestUnbondTx(options ...func(*TransactionMaker)) *tx.Tx {
	tmk := ts.NewTransactionMaker()

	for _, opt := range options {
		opt(tmk)
	}

	if tmk.Signer == nil {
		_, prv := ts.RandBLSKeyPair()
		tmk.Signer = prv
	}

	sender := tmk.SignerValidatorAddress()
	trx := tx.NewUnbondTx(tmk.LockTime, sender)
	ts.HelperSignTransaction(tmk.Signer, trx)

	return trx
}

// GenerateTestWithdrawTx generates a withdraw transaction for testing purposes.
func (ts *TestSuite) GenerateTestWithdrawTx(options ...func(*TransactionMaker)) *tx.Tx {
	tmk := ts.NewTransactionMaker()

	for _, opt := range options {
		opt(tmk)
	}

	if tmk.Signer == nil {
		_, prv := ts.RandBLSKeyPair()
		tmk.Signer = prv
	}

	sender := tmk.SignerValidatorAddress()
	trx := tx.NewWithdrawTx(tmk.LockTime, sender, ts.RandAccAddress(), tmk.Amount, tmk.Fee)
	ts.HelperSignTransaction(tmk.Signer, trx)

	return trx
}

// GenerateTestPrecommitVote generates a precommit vote for testing purposes.
func (ts *TestSuite) GenerateTestPrecommitVote(height uint32, round int16) (*vote.Vote, *bls.ValidatorKey) {
	valKey := ts.RandValKey()
	vote := vote.NewPrecommitVote(
		ts.RandHash(),
		height, round,
		valKey.Address())
	ts.HelperSignVote(valKey, vote)

	return vote, valKey
}

// GenerateTestPrepareVote generates a prepare vote for testing purposes.
func (ts *TestSuite) GenerateTestPrepareVote(height uint32, round int16) (*vote.Vote, *bls.ValidatorKey) {
	valKey := ts.RandValKey()
	vote := vote.NewPrepareVote(
		ts.RandHash(),
		height, round,
		valKey.Address())
	ts.HelperSignVote(valKey, vote)

	return vote, valKey
}

// GenerateTestCommittee generates a committee for testing purposes.
// All committee members have the same power.
func (ts *TestSuite) GenerateTestCommittee(num int) (committee.Committee, []*bls.ValidatorKey) {
	valKeys := make([]*bls.ValidatorKey, num)
	vals := make([]*validator.Validator, num)
	for index := int32(0); index < int32(num); index++ {
		valKey := ts.RandValKey()
		val := ts.GenerateTestValidator(
			ValidatorWithNumber(index),
			ValidatorWithPublicKey(valKey.PublicKey()))
		valKeys[index] = valKey
		vals[index] = val

		val.UpdateLastBondingHeight(1 + uint32(index))
		val.UpdateLastSortitionHeight(1 + uint32(index))
		val.SubtractFromStake(val.Stake())
		val.AddToStake(10e9)
	}

	cmt, _ := committee.NewCommittee(vals, num, vals[0].Address())

	return cmt, valKeys
}

func (*TestSuite) HelperSignVote(valKey *bls.ValidatorKey, v *vote.Vote) {
	sig := valKey.Sign(v.SignBytes())
	v.SetSignature(sig)
}

func (*TestSuite) HelperSignProposal(valKey *bls.ValidatorKey, p *proposal.Proposal) {
	sig := valKey.Sign(p.SignBytes())
	p.SetSignature(sig)
}

func (*TestSuite) HelperSignTransaction(prv crypto.PrivateKey, trx *tx.Tx) {
	sig := prv.Sign(trx.SignBytes())
	trx.SetSignature(sig)
	trx.SetPublicKey(prv.PublicKey())
}

func FindFreePort() int {
	var freePort int
	for {
		// Find a free TCP port
		listenerTCP, err := util.NetworkListen(context.Background(), "tcp", "127.0.0.1:0")
		if err != nil {
			continue
		}

		freePort = listenerTCP.Addr().(*net.TCPAddr).Port
		_ = listenerTCP.Close()

		udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("127.0.0.1:%d", freePort))
		if err != nil {
			continue
		}
		udpConn, err := net.ListenUDP("udp", udpAddr)
		if err != nil {
			continue
		}
		_ = udpConn.Close()

		break
	}

	return freePort
}
