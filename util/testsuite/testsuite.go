package testsuite

import (
	"encoding/hex"
	"math/rand"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/committee"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/sortition"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/tx"
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

// NewTestSuiteForSeed creates a new TestSuite with the given seed.
func NewTestSuiteForSeed(seed int64) *TestSuite {
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

// RandAmount returns a random amount between [0, 100^e9).
func (ts *TestSuite) RandAmount() amount.Amount {
	return ts.RandAmountRange(0, 1000e9)
}

// RandAmountRange returns a random amount between [min, max).
func (ts *TestSuite) RandAmountRange(min, max amount.Amount) amount.Amount {
	amt := amount.Amount(ts.RandInt64NonZero(int64(max - min)))

	return amt + min
}

// RandFee returns a random fee between [0.1, 1).
func (ts *TestSuite) RandFee() amount.Amount {
	fee := amount.Amount(ts.RandInt64(0.9e9) + 0.1e9)

	return fee
}

// RandBytes returns a slice of random bytes of the given length.
func (ts *TestSuite) RandBytes(length int) []byte {
	buf := make([]byte, length)
	_, err := ts.Rand.Read(buf)
	if err != nil {
		panic(err)
	}

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
	d, err := hex.DecodeString(in)
	if err != nil {
		panic(err)
	}

	return d
}

// RandBLSKeyPair generates a random BLS key pair for testing purposes.
func (ts *TestSuite) RandBLSKeyPair() (*bls.PublicKey, *bls.PrivateKey) {
	buf := make([]byte, bls.PrivateKeySize)
	_, err := ts.Rand.Read(buf)
	if err != nil {
		panic(err)
	}
	prv, _ := bls.PrivateKeyFromBytes(buf)
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

// RandHash generates a random hash for testing purposes.
func (ts *TestSuite) RandHash() hash.Hash {
	return hash.CalcHash(util.Int64ToSlice(ts.RandInt64(util.MaxInt64)))
}

// RandAccAddress generates a random account address for testing purposes.
func (ts *TestSuite) RandAccAddress() crypto.Address {
	addr := crypto.NewAddress(crypto.AddressTypeBLSAccount, ts.RandBytes(20))

	return addr
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

// GenerateTestAccount generates an account for testing purposes.
func (ts *TestSuite) GenerateTestAccount(number int32) (*account.Account, crypto.Address) {
	_, prv := ts.RandBLSKeyPair()
	acc := account.NewAccount(number)
	acc.AddToBalance(ts.RandAmount())

	return acc, prv.PublicKeyNative().AccountAddress()
}

// GenerateTestValidator generates a validator for testing purposes.
func (ts *TestSuite) GenerateTestValidator(number int32) (*validator.Validator, *bls.ValidatorKey) {
	pub, prv := ts.RandBLSKeyPair()
	val := validator.NewValidator(pub, number)
	val.AddToStake(ts.RandAmount())

	return val, bls.NewValidatorKey(prv)
}

type BlockMaker struct {
	Version   uint8
	Txs       block.Txs
	Proposer  crypto.Address
	Time      time.Time
	StateHash hash.Hash
	PrevHash  hash.Hash
	Seed      sortition.VerifiableSeed
	PrevCert  *certificate.BlockCertificate
}

// NewBlockMaker creates a new BlockMaker instance with default values.
func (ts *TestSuite) NewBlockMaker() *BlockMaker {
	txs := block.NewTxs()
	tx1 := ts.GenerateTestTransferTx()
	tx2 := ts.GenerateTestSortitionTx()
	tx3 := ts.GenerateTestBondTx()
	tx4 := ts.GenerateTestUnbondTx()
	tx5 := ts.GenerateTestWithdrawTx()

	txs.Append(tx1)
	txs.Append(tx2)
	txs.Append(tx3)
	txs.Append(tx4)
	txs.Append(tx5)

	return &BlockMaker{
		Version:  1,
		Txs:      txs,
		Proposer: ts.RandValAddress(),
		Time:     util.Now(),
		PrevHash: ts.RandHash(),
		Seed:     ts.RandSeed(),
		PrevCert: nil,
	}
}

// BlockWithVersion sets version to the block.
func BlockWithVersion(ver uint8) func(bm *BlockMaker) {
	return func(bm *BlockMaker) {
		bm.Version = ver
	}
}

// BlockWithProposer sets proposer address to the block.
func BlockWithProposer(addr crypto.Address) func(bm *BlockMaker) {
	return func(bm *BlockMaker) {
		bm.Proposer = addr
	}
}

// BlockWithTime sets block creation time to the block.
func BlockWithTime(t time.Time) func(bm *BlockMaker) {
	return func(bm *BlockMaker) {
		bm.Time = t
	}
}

// BlockWithStateHash sets state hash to the block.
func BlockWithStateHash(h hash.Hash) func(bm *BlockMaker) {
	return func(bm *BlockMaker) {
		bm.StateHash = h
	}
}

// BlockWithPrevHash sets previous block hash to the block.
func BlockWithPrevHash(h hash.Hash) func(bm *BlockMaker) {
	return func(bm *BlockMaker) {
		bm.PrevHash = h
	}
}

// BlockWithSeed sets verifiable seed to the block.
func BlockWithSeed(seed sortition.VerifiableSeed) func(bm *BlockMaker) {
	return func(bm *BlockMaker) {
		bm.Seed = seed
	}
}

// BlockWithPrevCert sets previous block certificate to the block.
func BlockWithPrevCert(cert *certificate.BlockCertificate) func(bm *BlockMaker) {
	return func(bm *BlockMaker) {
		bm.PrevCert = cert
	}
}

// BlockWithTransactions adds transactions to the block.
func BlockWithTransactions(txs block.Txs) func(bm *BlockMaker) {
	return func(bm *BlockMaker) {
		bm.Txs = txs
	}
}

// GenerateTestBlock generates a block for testing purposes with optional configuration.
func (ts *TestSuite) GenerateTestBlock(height uint32, options ...func(bm *BlockMaker)) (
	*block.Block, *certificate.BlockCertificate,
) {
	bm := ts.NewBlockMaker()
	bm.PrevCert = ts.GenerateTestBlockCertificate(height - 1)

	if height == 1 {
		bm.PrevCert = nil
		bm.PrevHash = hash.UndefHash
	}

	for _, opt := range options {
		opt(bm)
	}

	// blockCert := ts.GenerateTestBlockCertificate(height)
	header := block.NewHeader(bm.Version, bm.Time, bm.PrevHash, bm.PrevHash, bm.Seed, bm.Proposer)
	blk := block.NewBlock(header, bm.PrevCert, bm.Txs)

	blockCert := ts.GenerateTestBlockCertificate(height)
	return blk, blockCert
}

// GenerateTestBlockCertificate generates a block certificate for testing purposes.
func (ts *TestSuite) GenerateTestBlockCertificate(height uint32) *certificate.BlockCertificate {
	sig := ts.RandBLSSignature()

	cert := certificate.NewBlockCertificate(height, ts.RandRound())

	committers := ts.RandSlice(6)
	absentees := []int32{committers[5]}
	cert.SetSignature(committers, absentees, sig)

	return cert
}

// GenerateTestPrepareCertificate generates a prepare certificate for testing purposes.
func (ts *TestSuite) GenerateTestPrepareCertificate(height uint32) *certificate.VoteCertificate {
	sig := ts.RandBLSSignature()

	cert := certificate.NewVoteCertificate(height, ts.RandRound())

	committers := ts.RandSlice(6)
	absentees := []int32{committers[5]}
	cert.SetSignature(committers, absentees, sig)

	err := cert.BasicCheck()
	if err != nil {
		panic(err)
	}

	return cert
}

// GenerateTestProposal generates a proposal for testing purposes.
func (ts *TestSuite) GenerateTestProposal(height uint32, round int16) (*proposal.Proposal, *bls.ValidatorKey) {
	valKey := ts.RandValKey()
	blk, _ := ts.GenerateTestBlock(height, BlockWithProposer(valKey.Address()))
	prop := proposal.NewProposal(height, round, blk)
	ts.HelperSignProposal(valKey, prop)

	return prop, valKey
}

type TransactionMaker struct {
	Amount amount.Amount
	Fee    amount.Amount
	PrvKey *bls.PrivateKey
	PubKey *bls.PublicKey
}

// NewTransactionMaker creates a new TransactionMaker instance with default values.
func (ts *TestSuite) NewTransactionMaker() *TransactionMaker {
	pub, prv := ts.RandBLSKeyPair()
	return &TransactionMaker{
		Amount: ts.RandAmount(),
		Fee:    ts.RandFee(),
		PrvKey: prv,
		PubKey: pub,
	}
}

// TransactionWithAmount sets amount to the transaction.
func TransactionWithAmount(amt amount.Amount) func(tm *TransactionMaker) {
	return func(tm *TransactionMaker) {
		tm.Amount = amt
	}
}

// TransactionWithFee sets fee to the transaction.
func TransactionWithFee(fee amount.Amount) func(tm *TransactionMaker) {
	return func(tm *TransactionMaker) {
		tm.Fee = fee
	}
}

// TransactionWithSigner sets signer to the transaction.
func TransactionWithSigner(signer *bls.PrivateKey) func(tm *TransactionMaker) {
	return func(tm *TransactionMaker) {
		tm.PrvKey = signer
		tm.PubKey = signer.PublicKeyNative()
	}
}

// GenerateTestTransferTx generates a transfer transaction for testing purposes.
func (ts *TestSuite) GenerateTestTransferTx(options ...func(tm *TransactionMaker)) *tx.Tx {
	tm := ts.NewTransactionMaker()

	for _, opt := range options {
		opt(tm)
	}
	trx := tx.NewTransferTx(ts.RandHeight(), tm.PubKey.AccountAddress(), ts.RandAccAddress(),
		tm.Amount, tm.Fee, "test send-tx")
	ts.HelperSignTransaction(tm.PrvKey, trx)

	return trx
}

// GenerateTestBondTx generates a bond transaction for testing purposes.
func (ts *TestSuite) GenerateTestBondTx(options ...func(tm *TransactionMaker)) *tx.Tx {
	tm := ts.NewTransactionMaker()

	for _, opt := range options {
		opt(tm)
	}
	trx := tx.NewBondTx(ts.RandHeight(), tm.PubKey.AccountAddress(), ts.RandValAddress(),
		nil, tm.Amount, tm.Fee, "test bond-tx")
	ts.HelperSignTransaction(tm.PrvKey, trx)

	return trx
}

// GenerateTestSortitionTx generates a sortition transaction for testing purposes.
func (ts *TestSuite) GenerateTestSortitionTx(options ...func(tm *TransactionMaker)) *tx.Tx {
	tm := ts.NewTransactionMaker()

	for _, opt := range options {
		opt(tm)
	}
	proof := ts.RandProof()
	trx := tx.NewSortitionTx(ts.RandHeight(), tm.PubKey.ValidatorAddress(), proof)
	ts.HelperSignTransaction(tm.PrvKey, trx)

	return trx
}

// GenerateTestUnbondTx generates an unbond transaction for testing purposes.
func (ts *TestSuite) GenerateTestUnbondTx(options ...func(tm *TransactionMaker)) *tx.Tx {
	tm := ts.NewTransactionMaker()

	for _, opt := range options {
		opt(tm)
	}
	trx := tx.NewUnbondTx(ts.RandHeight(), tm.PubKey.ValidatorAddress(), "test unbond-tx")
	ts.HelperSignTransaction(tm.PrvKey, trx)

	return trx
}

// GenerateTestWithdrawTx generates a withdraw transaction for testing purposes.
func (ts *TestSuite) GenerateTestWithdrawTx(options ...func(tm *TransactionMaker)) *tx.Tx {
	tm := ts.NewTransactionMaker()

	for _, opt := range options {
		opt(tm)
	}
	trx := tx.NewWithdrawTx(ts.RandHeight(), tm.PubKey.ValidatorAddress(), ts.RandAccAddress(),
		tm.Amount, tm.Fee, "test withdraw-tx")
	ts.HelperSignTransaction(tm.PrvKey, trx)

	return trx
}

// GenerateTestPrecommitVote generates a precommit vote for testing purposes.
func (ts *TestSuite) GenerateTestPrecommitVote(height uint32, round int16) (*vote.Vote, *bls.ValidatorKey) {
	valKey := ts.RandValKey()
	v := vote.NewPrecommitVote(
		ts.RandHash(),
		height, round,
		valKey.Address())
	ts.HelperSignVote(valKey, v)

	return v, valKey
}

// GenerateTestPrepareVote generates a prepare vote for testing purposes.
func (ts *TestSuite) GenerateTestPrepareVote(height uint32, round int16) (*vote.Vote, *bls.ValidatorKey) {
	valKey := ts.RandValKey()
	v := vote.NewPrepareVote(
		ts.RandHash(),
		height, round,
		valKey.Address())
	ts.HelperSignVote(valKey, v)

	return v, valKey
}

// GenerateTestCommittee generates a committee for testing purposes.
// All committee members have the same power.
func (ts *TestSuite) GenerateTestCommittee(num int) (committee.Committee, []*bls.ValidatorKey) {
	if num < 4 {
		panic("the number of committee members must be at least 4")
	}
	valKeys := make([]*bls.ValidatorKey, num)
	vals := make([]*validator.Validator, num)
	for i := int32(0); i < int32(num); i++ {
		val, s := ts.GenerateTestValidator(i)
		valKeys[i] = s
		vals[i] = val

		val.UpdateLastBondingHeight(1 + uint32(i))
		val.UpdateLastSortitionHeight(1 + uint32(i))
		val.SubtractFromStake(val.Stake())
		val.AddToStake(10e9)
	}

	cmt, _ := committee.NewCommittee(vals, num, vals[0].Address())

	return cmt, valKeys
}

func (*TestSuite) HelperSignVote(valKey *bls.ValidatorKey, v *vote.Vote) {
	sig := valKey.Sign(v.SignBytes())
	v.SetSignature(sig)

	if err := v.BasicCheck(); err != nil {
		panic(err)
	}
}

func (*TestSuite) HelperSignProposal(valKey *bls.ValidatorKey, p *proposal.Proposal) {
	sig := valKey.Sign(p.SignBytes())
	p.SetSignature(sig)
}

func (*TestSuite) HelperSignTransaction(prv crypto.PrivateKey, trx *tx.Tx) {
	sig := prv.Sign(trx.SignBytes())
	trx.SetSignature(sig)
	trx.SetPublicKey(prv.PublicKey())

	if err := trx.BasicCheck(); err != nil {
		panic(err)
	}
}
