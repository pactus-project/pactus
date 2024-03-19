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
	return amount.Amount(ts.RandInt64(1000e9))
}

// func (ts *TestSuite) RandAmountInPAC() float64 {
// 	return util.ChangeToCoin(ts.RandAmount())
// }

// RandBytes returns a slice of random bytes of the given length.
func (ts *TestSuite) RandBytes(length int) []byte {
	buf := make([]byte, length)
	_, err := ts.Rand.Read(buf)
	if err != nil {
		panic(err)
	}

	return buf
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
func (ts *TestSuite) DecodingHex(in string) []byte {
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

// GenerateTestBlockWithProposer generates a block with the give proposer address for testing purposes.
func (ts *TestSuite) GenerateTestBlockWithProposer(height uint32, proposer crypto.Address,
) (*block.Block, *certificate.Certificate) {
	return ts.generateTestBlock(height, proposer, util.Now())
}

// GenerateTestBlockWithTime generates a block with the given time for testing purposes.
func (ts *TestSuite) GenerateTestBlockWithTime(height uint32, tme time.Time,
) (*block.Block, *certificate.Certificate) {
	return ts.generateTestBlock(height, ts.RandValAddress(), tme)
}

// GenerateTestBlock generates a block for testing purposes.
func (ts *TestSuite) GenerateTestBlock(height uint32) (*block.Block, *certificate.Certificate) {
	return ts.generateTestBlock(height, ts.RandValAddress(), util.Now())
}

func (ts *TestSuite) generateTestBlock(height uint32, proposer crypto.Address, tme time.Time,
) (*block.Block, *certificate.Certificate) {
	txs := block.NewTxs()
	tx1, _ := ts.GenerateTestTransferTx()
	tx2, _ := ts.GenerateTestSortitionTx()
	tx3, _ := ts.GenerateTestBondTx()
	tx4, _ := ts.GenerateTestUnbondTx()
	tx5, _ := ts.GenerateTestWithdrawTx()

	txs.Append(tx1)
	txs.Append(tx2)
	txs.Append(tx3)
	txs.Append(tx4)
	txs.Append(tx5)

	var prevCert *certificate.Certificate
	prevBlockHash := ts.RandHash()
	if height == 1 {
		prevCert = nil
		prevBlockHash = hash.UndefHash
	} else {
		prevCert = ts.GenerateTestCertificate(height - 1)
	}
	blockCert := ts.GenerateTestCertificate(height)
	header := block.NewHeader(1, tme,
		ts.RandHash(),
		prevBlockHash,
		ts.RandSeed(),
		proposer)

	blk := block.NewBlock(header, prevCert, txs)

	err := blk.BasicCheck()
	if err != nil {
		panic(err)
	}

	return blk, blockCert
}

// GenerateTestCertificate generates a certificate for testing purposes.
func (ts *TestSuite) GenerateTestCertificate(height uint32) *certificate.Certificate {
	sig := ts.RandBLSSignature()

	c1 := ts.RandInt32NonZero(10)
	c2 := ts.RandInt32NonZero(10) + 10
	c3 := ts.RandInt32NonZero(10) + 20
	c4 := ts.RandInt32NonZero(10) + 30
	cert := certificate.NewCertificate(
		height,
		ts.RandRound(),
		[]int32{c1, c2, c3, c4},
		[]int32{c2},
		sig)

	err := cert.BasicCheck()
	if err != nil {
		panic(err)
	}

	return cert
}

// GenerateTestProposal generates a proposal for testing purposes.
func (ts *TestSuite) GenerateTestProposal(height uint32, round int16) (*proposal.Proposal, *bls.ValidatorKey) {
	valKey := ts.RandValKey()
	blk, _ := ts.GenerateTestBlockWithProposer(height, valKey.Address())
	prop := proposal.NewProposal(height, round, blk)
	ts.HelperSignProposal(valKey, prop)

	return prop, valKey
}

// GenerateTestTransferTx generates a transfer transaction for testing purposes.
func (ts *TestSuite) GenerateTestTransferTx() (*tx.Tx, *bls.PrivateKey) {
	pub, prv := ts.RandBLSKeyPair()
	trx := tx.NewTransferTx(ts.RandHeight(), pub.AccountAddress(), ts.RandAccAddress(),
		ts.RandAmount(), ts.RandAmount(), "test send-tx")
	ts.HelperSignTransaction(prv, trx)

	return trx, prv
}

// GenerateTestBondTx generates a bond transaction for testing purposes.
func (ts *TestSuite) GenerateTestBondTx() (*tx.Tx, *bls.PrivateKey) {
	pub, prv := ts.RandBLSKeyPair()
	trx := tx.NewBondTx(ts.RandHeight(), pub.AccountAddress(), ts.RandValAddress(),
		nil, ts.RandAmount(), ts.RandAmount(), "test bond-tx")
	ts.HelperSignTransaction(prv, trx)

	return trx, prv
}

// GenerateTestSortitionTx generates a sortition transaction for testing purposes.
func (ts *TestSuite) GenerateTestSortitionTx() (*tx.Tx, *bls.PrivateKey) {
	pub, prv := ts.RandBLSKeyPair()
	proof := ts.RandProof()
	trx := tx.NewSortitionTx(ts.RandHeight(), pub.ValidatorAddress(), proof)
	ts.HelperSignTransaction(prv, trx)

	return trx, prv
}

// GenerateTestUnbondTx generates an unbond transaction for testing purposes.
func (ts *TestSuite) GenerateTestUnbondTx() (*tx.Tx, *bls.PrivateKey) {
	pub, prv := ts.RandBLSKeyPair()
	trx := tx.NewUnbondTx(ts.RandHeight(), pub.ValidatorAddress(), "test unbond-tx")
	ts.HelperSignTransaction(prv, trx)

	return trx, prv
}

// GenerateTestWithdrawTx generates a withdraw transaction for testing purposes.
func (ts *TestSuite) GenerateTestWithdrawTx() (*tx.Tx, *bls.PrivateKey) {
	pub, prv := ts.RandBLSKeyPair()
	trx := tx.NewWithdrawTx(ts.RandHeight(), pub.ValidatorAddress(), ts.RandAccAddress(),
		ts.RandAmount(), ts.RandAmount(), "test withdraw-tx")
	ts.HelperSignTransaction(prv, trx)

	return trx, prv
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
	valKeys := make([]*bls.ValidatorKey, num)
	vals := make([]*validator.Validator, num)
	for i := int32(0); i < int32(num); i++ {
		val, s := ts.GenerateTestValidator(i)
		valKeys[i] = s
		vals[i] = val

		val.UpdateLastBondingHeight(1 + uint32(i))
		val.UpdateLastSortitionHeight(1 + uint32(i))
		val.SubtractFromStake(val.Stake())
		val.AddToStake(10 * 1e9)
	}

	cmt, _ := committee.NewCommittee(vals, num, vals[0].Address())

	return cmt, valKeys
}

func (ts *TestSuite) HelperSignVote(valKey *bls.ValidatorKey, v *vote.Vote) {
	sig := valKey.Sign(v.SignBytes())
	v.SetSignature(sig)

	if err := v.BasicCheck(); err != nil {
		panic(err)
	}
}

func (ts *TestSuite) HelperSignProposal(valKey *bls.ValidatorKey, p *proposal.Proposal) {
	sig := valKey.Sign(p.SignBytes())
	p.SetSignature(sig)

	if err := p.BasicCheck(); err != nil {
		panic(err)
	}
}

func (ts *TestSuite) HelperSignTransaction(prv crypto.PrivateKey, trx *tx.Tx) {
	sig := prv.Sign(trx.SignBytes())
	trx.SetSignature(sig)
	trx.SetPublicKey(prv.PublicKey())

	if err := trx.BasicCheck(); err != nil {
		panic(err)
	}
}
