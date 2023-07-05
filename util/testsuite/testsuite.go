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

// NewTestSuiteForSeed creates a new TestSuite with the given seed.
func NewTestSuiteForSeed(seed int64) *TestSuite {
	return &TestSuite{
		Seed: seed,
		// nolint:gosec
		Rand: rand.New(rand.NewSource(seed)),
	}
}

// NewTestSuite creates a new TestSuite by generating new seed.
func NewTestSuite(t *testing.T) *TestSuite {
	seed := time.Now().UTC().UnixNano()
	t.Logf("%v seed is %v", t.Name(), seed)
	return &TestSuite{
		Seed: seed,
		// nolint:gosec
		Rand: rand.New(rand.NewSource(seed)),
	}
}

// RandInt16 returns a random int16 between 0 and max.
func (ts *TestSuite) RandInt16(max int16) int16 {
	return int16(ts.RandUint64(uint64(max)))
}

// RandUint16 returns a random uint16 between 0 and max.
func (ts *TestSuite) RandUint16(max uint16) uint16 {
	return uint16(ts.RandUint64(uint64(max)))
}

// RandInt32 returns a random int32 between 0 and max.
func (ts *TestSuite) RandInt32(max int32) int32 {
	return int32(ts.RandUint64(uint64(max)))
}

// RandUint32 returns a random uint32 between 0 and max.
func (ts *TestSuite) RandUint32(max uint32) uint32 {
	return uint32(ts.RandUint64(uint64(max)))
}

// RandInt64 returns a random int64 between 0 and max.
func (ts *TestSuite) RandInt64(max int64) int64 {
	return ts.Rand.Int63n(max)
}

// RandUint64 returns a random uint64 between 0 and max.
func (ts *TestSuite) RandUint64(max uint64) uint64 {
	return uint64(ts.RandInt64(int64(max)))
}

// RandInt returns a random int between 0 and max.
func (ts *TestSuite) RandInt(max int) int {
	return int(ts.RandInt64(int64(max))) + 1
}

// RandInt16NonZero returns a random int16 between 1 and max+1.
func (ts *TestSuite) RandInt16NonZero(max int16) int16 {
	return ts.RandInt16(max) + 1
}

// RandUint16NonZero returns a random uint16 between 1 and max+1.
func (ts *TestSuite) RandUint16NonZero(max uint16) uint16 {
	return ts.RandUint16(max) + 1
}

// RandInt32NonZero returns a random int32 between 1 and max+1.
func (ts *TestSuite) RandInt32NonZero(max int32) int32 {
	return ts.RandInt32(max) + 1
}

// RandUint32NonZero returns a random uint32 between 1 and max+1.
func (ts *TestSuite) RandUint32NonZero(max uint32) uint32 {
	return ts.RandUint32(max) + 1
}

// RandInt64NonZero returns a random int64 between 1 and max+1.
func (ts *TestSuite) RandInt64NonZero(max int64) int64 {
	return ts.RandInt64(max) + 1
}

// RandUint64NonZero returns a random uint64 between 1 and max+1.
func (ts *TestSuite) RandUint64NonZero(max uint64) uint64 {
	return ts.RandUint64(max) + 1
}

// RandIntNonZero returns a random int between 1 and max+1.
func (ts *TestSuite) RandIntNonZero(max int) int {
	return ts.RandInt(max) + 1
}

// RandomBytes returns a slice of random bytes of the given length.
func (ts *TestSuite) RandomBytes(len int) []byte {
	buf := make([]byte, len)
	_, err := ts.Rand.Read(buf)
	if err != nil {
		panic(err)
	}
	return buf
}

// RandomString generates a random string of the given length.
func (ts *TestSuite) RandomString(len int) string {
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

// RandomSigner generates a random signer for testing.
func (ts *TestSuite) RandomSigner() crypto.Signer {
	_, prv := ts.RandomBLSKeyPair()
	return crypto.NewSigner(prv)
}

// RandomBLSKeyPair generates a random BLS key pair for testing.
func (ts *TestSuite) RandomBLSKeyPair() (*bls.PublicKey, *bls.PrivateKey) {
	buf := make([]byte, bls.PrivateKeySize)
	_, err := ts.Rand.Read(buf)
	if err != nil {
		panic(err)
	}
	prv, _ := bls.PrivateKeyFromBytes(buf)
	pub := prv.PublicKey().(*bls.PublicKey)

	return pub, prv
}

// RandomHash generates a random hash for testing.
func (ts *TestSuite) RandomHash() hash.Hash {
	return hash.CalcHash(util.Int64ToSlice(ts.RandInt64(util.MaxInt64)))
}

// RandomStamp generates a random stamp for testing.
func (ts *TestSuite) RandomStamp() hash.Stamp {
	return ts.RandomHash().Stamp()
}

// RandomAddress generates a random address for testing.
func (ts *TestSuite) RandomAddress() crypto.Address {
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

// RandomSeed generates a random VerifiableSeed for testing.
func (ts *TestSuite) RandomSeed() sortition.VerifiableSeed {
	h := ts.RandomHash()
	signer := ts.RandomSigner()
	sig := signer.SignData(h.Bytes())
	seed, _ := sortition.VerifiableSeedFromBytes(sig.Bytes())
	return seed
}

// RandomProof generates a random Proof for testing.
func (ts *TestSuite) RandomProof() sortition.Proof {
	sig := ts.RandomSigner().SignData(ts.RandomHash().Bytes())
	proof, _ := sortition.ProofFromBytes(sig.Bytes())
	return proof
}

// RandomPeerID returns a random peer ID.
func (ts *TestSuite) RandomPeerID() peer.ID {
	s := ts.RandomBytes(32)
	id := [34]byte{0x12, 32}
	copy(id[2:], s[:])
	return peer.ID(id[:])
}

// GenerateTestAccount generates an account for testing purposes.
func (ts *TestSuite) GenerateTestAccount(number int32) (*account.Account, crypto.Signer) {
	signer := ts.RandomSigner()
	acc := account.NewAccount(number)
	acc.AddToBalance(ts.RandInt64(100 * 1e14))
	for i := 0; i < ts.RandInt(10); i++ {
		acc.IncSequence()
	}
	return acc, signer
}

// GenerateTestValidator generates a validator for testing purposes.
func (ts *TestSuite) GenerateTestValidator(number int32) (*validator.Validator, crypto.Signer) {
	pub, pv := ts.RandomBLSKeyPair()
	val := validator.NewValidator(pub, number)
	val.AddToStake(ts.RandInt64(100 * 1e9))
	for i := 0; i < ts.RandInt(10); i++ {
		val.IncSequence()
	}
	return val, crypto.NewSigner(pv)
}

// GenerateTestBlock generates a block vote for testing.
func (ts *TestSuite) GenerateTestBlock(proposer *crypto.Address, prevBlockHash *hash.Hash) *block.Block {
	if proposer == nil {
		addr := ts.RandomAddress()
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

	if prevBlockHash == nil {
		h := ts.RandomHash()
		prevBlockHash = &h
	}
	cert := ts.GenerateTestCertificate(*prevBlockHash)
	if prevBlockHash.IsUndef() {
		cert = nil
	}
	sortitionSeed := ts.RandomSeed()
	header := block.NewHeader(1, util.Now(),
		ts.RandomHash(),
		*prevBlockHash,
		sortitionSeed,
		*proposer)

	return block.NewBlock(header, cert, txs)
}

// GenerateTestCertificate generates a certificate for testing.
func (ts *TestSuite) GenerateTestCertificate(blockHash hash.Hash) *block.Certificate {
	_, priv2 := ts.RandomBLSKeyPair()
	_, priv3 := ts.RandomBLSKeyPair()
	_, priv4 := ts.RandomBLSKeyPair()

	sigs := []*bls.Signature{
		priv2.Sign(blockHash.Bytes()).(*bls.Signature),
		priv3.Sign(blockHash.Bytes()).(*bls.Signature),
		priv4.Sign(blockHash.Bytes()).(*bls.Signature),
	}
	sig := bls.SignatureAggregate(sigs)

	c1 := ts.RandInt32(10)
	c2 := ts.RandInt32(10) + 10
	c3 := ts.RandInt32(10) + 20
	c4 := ts.RandInt32(10) + 30
	return block.NewCertificate(
		ts.RandInt16(10),
		[]int32{c1, c2, c3, c4},
		[]int32{c2},
		sig)
}

// GenerateTestProposal generates a proposal for testing.
func (ts *TestSuite) GenerateTestProposal(height uint32, round int16) (*proposal.Proposal, crypto.Signer) {
	signer := ts.RandomSigner()
	addr := signer.Address()
	b := ts.GenerateTestBlock(&addr, nil)
	p := proposal.NewProposal(height, round, b)
	signer.SignMsg(p)
	return p, signer
}

// GenerateTestTransferTx generates a transfer transaction for testing.
func (ts *TestSuite) GenerateTestTransferTx() (*tx.Tx, crypto.Signer) {
	stamp := ts.RandomStamp()
	s := ts.RandomSigner()
	pub, _ := ts.RandomBLSKeyPair()
	tx := tx.NewTransferTx(stamp, ts.RandInt32(1000), s.Address(), pub.Address(),
		ts.RandInt64(1000*1e10), ts.RandInt64(1*1e10), "test send-tx")
	s.SignMsg(tx)
	return tx, s
}

// GenerateTestBondTx generates a bond transaction for testing.
func (ts *TestSuite) GenerateTestBondTx() (*tx.Tx, crypto.Signer) {
	stamp := ts.RandomStamp()
	s := ts.RandomSigner()
	pub, _ := ts.RandomBLSKeyPair()
	tx := tx.NewBondTx(stamp, ts.RandInt32(1000), s.Address(), pub.Address(),
		pub, ts.RandInt64(1000*1e10), ts.RandInt64(1*1e10), "test bond-tx")
	s.SignMsg(tx)
	return tx, s
}

// GenerateTestSortitionTx generates a sortition transaction for testing.
func (ts *TestSuite) GenerateTestSortitionTx() (*tx.Tx, crypto.Signer) {
	stamp := ts.RandomStamp()
	s := ts.RandomSigner()
	proof := ts.RandomProof()
	tx := tx.NewSortitionTx(stamp, ts.RandInt32(1000), s.Address(), proof)
	s.SignMsg(tx)
	return tx, s
}

// GenerateTestUnbondTx generates an unbond transaction for testing.
func (ts *TestSuite) GenerateTestUnbondTx() (*tx.Tx, crypto.Signer) {
	stamp := ts.RandomStamp()
	s := ts.RandomSigner()
	tx := tx.NewUnbondTx(stamp, ts.RandInt32(1000), s.Address(), "test unbond-tx")
	s.SignMsg(tx)
	return tx, s
}

// GenerateTestWithdrawTx generates a withdraw transaction for testing.
func (ts *TestSuite) GenerateTestWithdrawTx() (*tx.Tx, crypto.Signer) {
	stamp := ts.RandomStamp()
	s := ts.RandomSigner()
	tx := tx.NewWithdrawTx(stamp, ts.RandInt32(1000), s.Address(), ts.RandomAddress(),
		ts.RandInt64(1000*1e10), ts.RandInt64(1*1e10), "test withdraw-tx")
	s.SignMsg(tx)
	return tx, s
}

// GenerateTestPrecommitVote generates a precommit vote for testing.
func (ts *TestSuite) GenerateTestPrecommitVote(height uint32, round int16) (*vote.Vote, crypto.Signer) {
	s := ts.RandomSigner()
	v := vote.NewVote(
		vote.VoteTypePrecommit,
		height, round,
		ts.RandomHash(),
		s.Address())
	s.SignMsg(v)

	return v, s
}

// GenerateTestPrepareVote generates a prepare vote for testing.
func (ts *TestSuite) GenerateTestPrepareVote(height uint32, round int16) (*vote.Vote, crypto.Signer) {
	s := ts.RandomSigner()
	v := vote.NewVote(
		vote.VoteTypePrepare,
		height, round,
		ts.RandomHash(),
		s.Address())
	s.SignMsg(v)

	return v, s
}

// GenerateTestChangeProposerVote generates a proposer-change vote for testing.
func (ts *TestSuite) GenerateTestChangeProposerVote(height uint32, round int16) (*vote.Vote, crypto.Signer) {
	s := ts.RandomSigner()
	v := vote.NewVote(
		vote.VoteTypeChangeProposer,
		height, round,
		ts.RandomHash(),
		s.Address())
	s.SignMsg(v)

	return v, s
}

// GenerateTestCommittee generates a committee for testing purposes.
// All committee members have the same power.
func (ts *TestSuite) GenerateTestCommittee(num int) (committee.Committee, []crypto.Signer) {
	signers := make([]crypto.Signer, num)
	vals := make([]*validator.Validator, num)
	h1 := ts.RandUint32(100000)
	for i := int32(0); i < int32(num); i++ {
		val, s := ts.GenerateTestValidator(i)
		signers[i] = s
		vals[i] = val

		val.UpdateLastBondingHeight(h1 + uint32(i))
		val.UpdateLastJoinedHeight(h1 + 100 + uint32(i))
		val.SubtractFromStake(val.Stake())
		val.AddToStake(10 * 1e9)
	}

	committee, _ := committee.NewCommittee(vals, num, vals[0].Address())
	return committee, signers
}
