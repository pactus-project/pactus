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
		//nolint:gosec
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
		//nolint:gosec
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

// RandHeight returns a random number between [1, 1,000,000] for block height.
func (ts *TestSuite) RandHeight() uint32 {
	return ts.RandUint32NonZero(1e6)
}

// RandRound returns a random number between [0, 10) for block round.
func (ts *TestSuite) RandRound() int16 {
	return ts.RandInt16(10)
}

// RandBytes returns a slice of random bytes of the given length.
func (ts *TestSuite) RandBytes(len int) []byte {
	buf := make([]byte, len)
	_, err := ts.Rand.Read(buf)
	if err != nil {
		panic(err)
	}
	return buf
}

// RandString generates a random string of the given length.
func (ts *TestSuite) RandString(len int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, len)
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

// RandSigner generates a random signer for testing purposes.
func (ts *TestSuite) RandSigner() crypto.Signer {
	_, prv := ts.RandBLSKeyPair()
	return crypto.NewSigner(prv)
}

// RandBLSKeyPair generates a random BLS key pair for testing purposes.
func (ts *TestSuite) RandBLSKeyPair() (*bls.PublicKey, *bls.PrivateKey) {
	buf := make([]byte, bls.PrivateKeySize)
	_, err := ts.Rand.Read(buf)
	if err != nil {
		panic(err)
	}
	prv, _ := bls.PrivateKeyFromBytes(buf)
	pub := prv.PublicKey().(*bls.PublicKey)

	return pub, prv
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

// RandStamp generates a random stamp for testing purposes.
func (ts *TestSuite) RandStamp() hash.Stamp {
	return ts.RandHash().Stamp()
}

// RandAddress generates a random address for testing purposes.
func (ts *TestSuite) RandAddress() crypto.Address {
	data := make([]byte, 20)
	_, err := ts.Rand.Read(data)
	if err != nil {
		panic(err)
	}
	data = append([]byte{1}, data...)
	var addr crypto.Address
	copy(addr[:], data[:])
	return addr
}

// RandSeed generates a random VerifiableSeed for testing purposes.
func (ts *TestSuite) RandSeed() sortition.VerifiableSeed {
	h := ts.RandHash()
	signer := ts.RandSigner()
	sig := signer.SignData(h.Bytes())
	seed, _ := sortition.VerifiableSeedFromBytes(sig.Bytes())
	return seed
}

// RandProof generates a random Proof for testing purposes.
func (ts *TestSuite) RandProof() sortition.Proof {
	sig := ts.RandSigner().SignData(ts.RandHash().Bytes())
	proof, _ := sortition.ProofFromBytes(sig.Bytes())
	return proof
}

// RandPeerID returns a random peer ID.
func (ts *TestSuite) RandPeerID() peer.ID {
	s := ts.RandBytes(32)
	id := [34]byte{0x12, 32}
	copy(id[2:], s[:])
	return peer.ID(id[:])
}

// GenerateTestAccount generates an account for testing purposes.
func (ts *TestSuite) GenerateTestAccount(number int32) (*account.Account, crypto.Signer) {
	signer := ts.RandSigner()
	acc := account.NewAccount(number)
	acc.AddToBalance(ts.RandInt64(100 * 1e14))
	for i := 0; i < ts.RandInt(10); i++ {
		acc.IncSequence()
	}
	return acc, signer
}

// GenerateTestValidator generates a validator for testing purposes.
func (ts *TestSuite) GenerateTestValidator(number int32) (*validator.Validator, crypto.Signer) {
	pub, pv := ts.RandBLSKeyPair()
	val := validator.NewValidator(pub, number)
	val.AddToStake(ts.RandInt64(100 * 1e9))
	for i := 0; i < ts.RandInt(10); i++ {
		val.IncSequence()
	}
	return val, crypto.NewSigner(pv)
}

// GenerateTestBlockWithTime generates a block at the give time for testing purposes.
func (ts *TestSuite) GenerateTestBlockWithTime(proposer *crypto.Address, time time.Time) *block.Block {
	if proposer == nil {
		addr := ts.RandAddress()
		proposer = &addr
	}
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

	cert := ts.GenerateTestCertificate()
	header := block.NewHeader(1, time,
		ts.RandHash(),
		ts.RandHash(),
		ts.RandSeed(),
		*proposer)

	return block.NewBlock(header, cert, txs)
}

// GenerateTestBlock generates a block for testing purposes.
func (ts *TestSuite) GenerateTestBlock(proposer *crypto.Address) *block.Block {
	return ts.GenerateTestBlockWithTime(proposer, util.Now())
}

// GenerateTestCertificate generates a certificate for testing purposes.
func (ts *TestSuite) GenerateTestCertificate() *certificate.Certificate {
	sig := ts.RandBLSSignature()

	c1 := ts.RandInt32NonZero(10)
	c2 := ts.RandInt32NonZero(10) + 10
	c3 := ts.RandInt32NonZero(10) + 20
	c4 := ts.RandInt32NonZero(10) + 30
	return certificate.NewCertificate(
		ts.RandHeight(),
		ts.RandRound(),
		[]int32{c1, c2, c3, c4},
		[]int32{c2},
		sig)
}

// GenerateTestProposal generates a proposal for testing purposes.
func (ts *TestSuite) GenerateTestProposal(height uint32, round int16) (*proposal.Proposal, crypto.Signer) {
	signer := ts.RandSigner()
	addr := signer.Address()
	b := ts.GenerateTestBlock(&addr)
	p := proposal.NewProposal(height, round, b)
	signer.SignMsg(p)
	return p, signer
}

// GenerateTestTransferTx generates a transfer transaction for testing purposes.
func (ts *TestSuite) GenerateTestTransferTx() (*tx.Tx, crypto.Signer) {
	stamp := ts.RandStamp()
	s := ts.RandSigner()
	pub, _ := ts.RandBLSKeyPair()
	tx := tx.NewTransferTx(stamp, ts.RandInt32(1000), s.Address(), pub.Address(),
		ts.RandInt64(1000*1e10), ts.RandInt64(1*1e10), "test send-tx")
	s.SignMsg(tx)
	return tx, s
}

// GenerateTestBondTx generates a bond transaction for testing purposes.
func (ts *TestSuite) GenerateTestBondTx() (*tx.Tx, crypto.Signer) {
	stamp := ts.RandStamp()
	s := ts.RandSigner()
	pub, _ := ts.RandBLSKeyPair()
	tx := tx.NewBondTx(stamp, ts.RandInt32(1000), s.Address(), pub.Address(),
		pub, ts.RandInt64(1000*1e10), ts.RandInt64(1*1e10), "test bond-tx")
	s.SignMsg(tx)
	return tx, s
}

// GenerateTestSortitionTx generates a sortition transaction for testing purposes.
func (ts *TestSuite) GenerateTestSortitionTx() (*tx.Tx, crypto.Signer) {
	stamp := ts.RandStamp()
	s := ts.RandSigner()
	proof := ts.RandProof()
	tx := tx.NewSortitionTx(stamp, ts.RandInt32(1000), s.Address(), proof)
	s.SignMsg(tx)
	return tx, s
}

// GenerateTestUnbondTx generates an unbond transaction for testing purposes.
func (ts *TestSuite) GenerateTestUnbondTx() (*tx.Tx, crypto.Signer) {
	stamp := ts.RandStamp()
	s := ts.RandSigner()
	tx := tx.NewUnbondTx(stamp, ts.RandInt32(1000), s.Address(), "test unbond-tx")
	s.SignMsg(tx)
	return tx, s
}

// GenerateTestWithdrawTx generates a withdraw transaction for testing purposes.
func (ts *TestSuite) GenerateTestWithdrawTx() (*tx.Tx, crypto.Signer) {
	stamp := ts.RandStamp()
	s := ts.RandSigner()
	tx := tx.NewWithdrawTx(stamp, ts.RandInt32(1000), s.Address(), ts.RandAddress(),
		ts.RandInt64(1000*1e10), ts.RandInt64(1*1e10), "test withdraw-tx")
	s.SignMsg(tx)
	return tx, s
}

// GenerateTestPrecommitVote generates a precommit vote for testing purposes.
func (ts *TestSuite) GenerateTestPrecommitVote(height uint32, round int16) (*vote.Vote, crypto.Signer) {
	s := ts.RandSigner()
	v := vote.NewPrecommitVote(
		ts.RandHash(),
		height, round,
		s.Address())
	s.SignMsg(v)

	return v, s
}

// GenerateTestPrepareVote generates a prepare vote for testing purposes.
func (ts *TestSuite) GenerateTestPrepareVote(height uint32, round int16) (*vote.Vote, crypto.Signer) {
	s := ts.RandSigner()
	v := vote.NewPrepareVote(
		ts.RandHash(),
		height, round,
		s.Address())
	s.SignMsg(v)

	return v, s
}

// GenerateTestCommittee generates a committee for testing purposes.
// All committee members have the same power.
func (ts *TestSuite) GenerateTestCommittee(num int) (committee.Committee, []crypto.Signer) {
	signers := make([]crypto.Signer, num)
	vals := make([]*validator.Validator, num)
	for i := int32(0); i < int32(num); i++ {
		val, s := ts.GenerateTestValidator(i)
		signers[i] = s
		vals[i] = val

		val.UpdateLastBondingHeight(1 + uint32(i))
		val.UpdateLastSortitionHeight(1 + uint32(i))
		val.SubtractFromStake(val.Stake())
		val.AddToStake(10 * 1e9)
	}

	committee, _ := committee.NewCommittee(vals, num, vals[0].Address())
	return committee, signers
}
