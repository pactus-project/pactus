package tx_test

import (
	"encoding/hex"
	"fmt"
	"strings"
	"testing"

	"github.com/fxamacker/cbor/v2"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/ed25519"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCBORMarshaling(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	tx1 := ts.GenerateTestTransferTx()
	bz, err := cbor.Marshal(tx1)
	assert.NoError(t, err)
	tx2 := new(tx.Tx)
	assert.NoError(t, cbor.Unmarshal(bz, tx2))
	assert.Equal(t, tx1.ID(), tx2.ID())

	assert.Error(t, cbor.Unmarshal([]byte{1}, tx2))
}

func TestEncodingTx(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	trx1 := ts.GenerateTestTransferTx()
	trx2 := ts.GenerateTestBondTx()
	trx3 := ts.GenerateTestUnbondTx()
	trx4 := ts.GenerateTestWithdrawTx()
	trx5 := ts.GenerateTestSortitionTx()

	tests := []*tx.Tx{trx1, trx2, trx3, trx4, trx5}
	for _, trx := range tests {
		assert.NoError(t, trx.BasicCheck())
		assert.NoError(t, trx.BasicCheck()) // double basic check

		length := trx.SerializeSize()
		for i := 0; i < length; i++ {
			w := util.NewFixedWriter(i)
			assert.Error(t, trx.Encode(w), "encode test %v failed", i)
		}
		w := util.NewFixedWriter(length)
		assert.NoError(t, trx.Encode(w))

		for i := 0; i < length; i++ {
			newTrx := new(tx.Tx)
			r := util.NewFixedReader(i, w.Bytes())
			assert.Error(t, newTrx.Decode(r), "decode test %v failed", i)
		}

		bz, err := trx.Bytes()
		assert.NoError(t, err)
		decodedTrx, err := tx.FromBytes(bz)
		assert.NoError(t, err)
		assert.Equal(t, trx.ID(), decodedTrx.ID())
	}

	_, err := tx.FromBytes([]byte{1})
	assert.Error(t, err)
}

func TestTxIDNoSignatory(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	tx1 := ts.GenerateTestTransferTx()
	tx2 := new(tx.Tx)
	*tx2 = *tx1

	tx2.SetPublicKey(nil)
	tx2.SetSignature(nil)

	require.True(t, tx1.IsSigned())
	require.False(t, tx2.IsSigned())

	require.Equal(t, tx1.ID(), tx2.ID())
	require.Equal(t, tx1.SignBytes(), tx2.SignBytes())
}

func TestBasicCheck(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	t.Run("LockTime is not defined", func(t *testing.T) {
		trx := tx.NewTransferTx(0,
			ts.RandAccAddress(), ts.RandAccAddress(), ts.RandAmount(), ts.RandFee())

		err := trx.BasicCheck()
		assert.ErrorIs(t, err, tx.BasicCheckError{
			Reason: "lock time is not defined",
		})
	})

	t.Run("Big memo, Should returns error", func(t *testing.T) {
		bigMemo := strings.Repeat("a", 65)

		trx := tx.NewTransferTx(ts.RandHeight(),
			ts.RandAccAddress(), ts.RandAccAddress(), ts.RandAmount(), ts.RandFee(), tx.WithMemo(bigMemo))

		err := trx.BasicCheck()
		assert.ErrorIs(t, err, tx.BasicCheckError{
			Reason: "memo length exceeded: 65",
		})
	})

	t.Run("Invalid payload, Should returns error", func(t *testing.T) {
		invAddr := ts.RandAccAddress()
		invAddr[0] = 4
		trx := tx.NewTransferTx(ts.RandHeight(), ts.RandAccAddress(), invAddr, 1e9, ts.RandFee())

		err := trx.BasicCheck()
		assert.ErrorIs(t, err, tx.BasicCheckError{
			Reason: "invalid payload: receiver is not an account address: " + invAddr.String(),
		})
	})

	t.Run("Invalid amount", func(t *testing.T) {
		trx := tx.NewTransferTx(ts.RandHeight(), ts.RandAccAddress(), ts.RandAccAddress(), -1, 1)

		err := trx.BasicCheck()
		assert.ErrorIs(t, err, tx.BasicCheckError{
			Reason: "invalid amount: -0.000000001 PAC",
		})
	})

	t.Run("Invalid amount", func(t *testing.T) {
		trx := tx.NewTransferTx(ts.RandHeight(), ts.RandAccAddress(), ts.RandAccAddress(), (42e15)+1, 1)

		err := trx.BasicCheck()
		assert.ErrorIs(t, err, tx.BasicCheckError{
			Reason: "invalid amount: 42000000 PAC",
		})
	})

	t.Run("Invalid fee", func(t *testing.T) {
		trx := tx.NewTransferTx(ts.RandHeight(), ts.RandAccAddress(), ts.RandAccAddress(), 1, -1)

		err := trx.BasicCheck()
		assert.ErrorIs(t, err, tx.BasicCheckError{
			Reason: "invalid fee: -0.000000001 PAC",
		})
	})

	t.Run("Invalid fee", func(t *testing.T) {
		trx := tx.NewTransferTx(ts.RandHeight(), ts.RandAccAddress(), ts.RandAccAddress(), 1, (42e15)+1)

		err := trx.BasicCheck()
		assert.ErrorIs(t, err, tx.BasicCheckError{
			Reason: "invalid fee: 42000000 PAC",
		})
	})

	t.Run("Invalid signer address", func(t *testing.T) {
		valKey := ts.RandValKey()
		trx := tx.NewTransferTx(ts.RandHeight(), ts.RandAccAddress(), ts.RandAccAddress(), 1, 1)
		sig := valKey.PrivateKey().Sign(trx.SignBytes())
		trx.SetSignature(sig)
		trx.SetPublicKey(valKey.PublicKey())

		err := trx.BasicCheck()
		assert.ErrorIs(t, err, tx.BasicCheckError{
			Reason: fmt.Sprintf("address mismatch: expected %s, got %s",
				valKey.PublicKey().AccountAddress(), trx.Payload().Signer()),
		})
	})

	t.Run("Invalid version", func(t *testing.T) {
		data := ts.DecodingHex(
			"02" + // Flags
				"02" + // Version
				"01020304" + // LockTime
				"01" + // Fee
				"00" + // Memo
				"01" + // PayloadType
				"00" + // Sender (treasury)
				"012222222222222222222222222222222222222222" + // Receiver
				"01") // Amount

		trx, err := tx.FromBytes(data)
		assert.NoError(t, err)
		err = trx.BasicCheck()
		assert.ErrorIs(t, err, tx.BasicCheckError{
			Reason: "invalid version: 2",
		})
		assert.Equal(t, len(data), trx.SerializeSize())
	})
}

func TestInvalidPayloadType(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	data := ts.DecodingHex(
		"02" + // Flags
			"01" + // Version
			"01020300" + // LockTime
			"01" + // Fee
			"00" + // Memo
			"06" + // PayloadType
			"00" + // Sender (treasury)
			"012222222222222222222222222222222222222222" + // Receiver
			"01") // Amount

	_, err := tx.FromBytes(data)
	assert.ErrorIs(t, err, tx.InvalidPayloadTypeError{
		PayloadType: payload.Type(6),
	})
}

func TestSubsidyTx(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pub, prv := ts.RandEd25519KeyPair()

	t.Run("Has signature", func(t *testing.T) {
		trx := tx.NewSubsidyTx(ts.RandHeight(), pub.AccountAddress(), 2500)
		sig := prv.Sign(trx.SignBytes())
		trx.SetSignature(sig)

		err := trx.BasicCheck()
		assert.ErrorIs(t, err, tx.BasicCheckError{
			Reason: "subsidy transaction with signatory",
		})
	})

	t.Run("Has public key", func(t *testing.T) {
		trx := tx.NewSubsidyTx(ts.RandHeight(), pub.AccountAddress(), 2500)
		trx.SetPublicKey(pub)

		err := trx.BasicCheck()
		assert.ErrorIs(t, err, tx.BasicCheckError{
			Reason: "subsidy transaction with signatory",
		})
	})

	t.Run("Strip public key", func(t *testing.T) {
		trx := tx.NewSubsidyTx(ts.RandHeight(), pub.AccountAddress(), 2500)
		trx.StripPublicKey()

		err := trx.BasicCheck()
		assert.NoError(t, err)
		assert.False(t, trx.IsPublicKeyStriped())
	})
}

func TestInvalidSignature(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	t.Run("Good", func(t *testing.T) {
		trx := ts.GenerateTestTransferTx()
		assert.NoError(t, trx.BasicCheck())
	})

	t.Run("No signature", func(t *testing.T) {
		trx := ts.GenerateTestTransferTx()
		trx.SetSignature(nil)

		err := trx.BasicCheck()
		assert.ErrorIs(t, err, tx.BasicCheckError{
			Reason: "no signature",
		})
	})

	t.Run("No public key", func(t *testing.T) {
		trx := ts.GenerateTestTransferTx()
		trx.SetPublicKey(nil)

		err := trx.BasicCheck()
		assert.ErrorIs(t, err, tx.BasicCheckError{
			Reason: "no public key",
		})
	})

	pbInv, pvInv := ts.RandBLSKeyPair()
	t.Run("Invalid signature", func(t *testing.T) {
		trx := ts.GenerateTestTransferTx()
		sig := pvInv.Sign(trx.SignBytes())
		trx.SetSignature(sig)

		err := trx.BasicCheck()
		assert.ErrorIs(t, err, tx.BasicCheckError{
			Reason: "invalid signature",
		})
	})

	t.Run("Invalid public key", func(t *testing.T) {
		trx := ts.GenerateTestTransferTx()
		trx.SetPublicKey(pbInv)

		err := trx.BasicCheck()
		assert.ErrorIs(t, err, tx.BasicCheckError{
			Reason: fmt.Sprintf("address mismatch: expected %s, got %s", pbInv.AccountAddress(), trx.Payload().Signer()),
		})
	})

	t.Run("Invalid sign Bytes", func(t *testing.T) {
		valKey := ts.RandValKey()
		trx0 := ts.GenerateTestUnbondTx(testsuite.TransactionWithBLSSigner(valKey.PrivateKey()))

		trx := tx.NewUnbondTx(trx0.LockTime(), valKey.Address(), tx.WithMemo("invalidate signature"))
		trx.SetPublicKey(trx0.PublicKey())
		trx.SetSignature(trx0.Signature())

		err := trx.BasicCheck()
		assert.ErrorIs(t, err, tx.BasicCheckError{
			Reason: "invalid signature",
		})
	})

	t.Run("Zero signature", func(t *testing.T) {
		trx := ts.GenerateTestTransferTx()
		trx.SetSignature(&bls.Signature{})

		err := trx.BasicCheck()
		assert.ErrorIs(t, err, tx.BasicCheckError{
			Reason: "invalid signature",
		})
	})

	t.Run("Zero public key", func(t *testing.T) {
		trx := ts.GenerateTestTransferTx()
		zeroPubKey := &bls.PublicKey{}
		trx.SetPublicKey(zeroPubKey)

		err := trx.BasicCheck()
		assert.ErrorIs(t, err, tx.BasicCheckError{
			Reason: fmt.Sprintf("address mismatch: expected %s, got %s",
				zeroPubKey.AccountAddress().String(), trx.Payload().Signer()),
		})
	})
}

func TestSignBytesBLS(t *testing.T) {
	data, _ := hex.DecodeString(
		"00" + // Flags
			"01" + // Version
			"01020304" + // LockTime
			"e807" + // Fee
			"0474657374" + // Memo
			"01" + // PayloadType
			"022109324a70c5bd7a0591bcd8d597cb3a06a91770" + // Sender (pc1zyyynyjnsck7h5pv3hnvdt97t8gr2j9ms0vrrp7)
			"02b00c16e60f46390455baff5a8b69ac70e67f10c8" + // Receiver (pc1zkqxpdes0gcusg4d6ladgk6dvwrn87yxgumzwn3)
			"a09c01" + // Amount
			"86d45b6532d447070cf1ee67d4a04a13f337f6a2bfd6c54419efdd4b502b529d3f3be52567d8adaf494e0edc93d4ae51" + // Signature
			"b805043a816c3213c67f365f83c6946546049f517ebe470f186b36ff53fb996ae2468b119582a7f18fe8f0bfb4e055d5" + // PublicKey
			"190601a983fb4636c36287a73d80dbb14f244f319da5eeac02ce7ee9026245ac36b9978cabd6d2cbb3c1f87e55e2fc29")

	txID, _ := hash.FromString("7ab1287fe4882918e69b9f83215378ea08f2d91e0700c2e35a73b7aae1d7bf2d")
	trx, err := tx.FromBytes(data)
	assert.NoError(t, err)
	assert.Equal(t, len(data), trx.SerializeSize())

	signBytes := data[1 : len(data)-bls.PublicKeySize-bls.SignatureSize]
	assert.Equal(t, signBytes, trx.SignBytes())
	assert.Equal(t, hash.CalcHash(signBytes), trx.ID())
	assert.Equal(t, txID, trx.ID())
	assert.Equal(t, uint32(0x04030201), trx.LockTime())
	assert.Equal(t, "test", trx.Memo())
	assert.Equal(t, amount.Amount(1000), trx.Fee())
	assert.Equal(t, amount.Amount(20000), trx.Payload().Value())
	assert.Equal(t, "pc1zyyynyjnsck7h5pv3hnvdt97t8gr2j9ms0vrrp7", trx.Payload().Signer().String())
}

func TestSignBytesEd25519(t *testing.T) {
	data, _ := hex.DecodeString(
		"00" + // Flags
			"01" + // Version
			"01020300" + // LockTime
			"e807" + // Fee (1000)
			"0474657374" + // Memo ("test")
			"01" + // PayloadType
			"037098338e0b6808119dfd4457ab806b9c2059b89b" + // Sender (pc1rwzvr8rstdqypr80ag3t6hqrtnss9nwymcxy3lr)
			"037a14ae24533816e7faaa6ed28fcdde8e55a7df21" + // Receiver (pc1r0g22ufzn8qtw0742dmfglnw73e260hep0k3yra)
			"a09c01" + // Amount (20000)
			"95794161374b22c696dabb98e93f6ca9300b22f3b904921fbf560bb72145f4fa" + // Signature
			"50ac25c7125271489b0cd230549257c93fb8c6265f2914a988ba7b81c1bc47ff" + // PublicKey
			"f027412dd59447867911035ff69742d171060a1f132ac38b95acc6e39ec0bd09")

	txID, _ := hash.FromString("34cd4656a98f7eb996e83efdc384cefbe3a9c52dca79a99245b4eacc0b0b4311")
	trx, err := tx.FromBytes(data)
	assert.NoError(t, err)
	assert.Equal(t, len(data), trx.SerializeSize())

	signBytes := data[1 : len(data)-ed25519.PublicKeySize-ed25519.SignatureSize]
	assert.Equal(t, signBytes, trx.SignBytes())
	assert.Equal(t, hash.CalcHash(signBytes), trx.ID())
	assert.Equal(t, txID, trx.ID())
	assert.Equal(t, uint32(0x00030201), trx.LockTime())
	assert.Equal(t, "test", trx.Memo())
	assert.Equal(t, amount.Amount(1000), trx.Fee())
	assert.Equal(t, amount.Amount(20000), trx.Payload().Value())
	assert.Equal(t, "pc1rwzvr8rstdqypr80ag3t6hqrtnss9nwymcxy3lr", trx.Payload().Signer().String())
}

func TestStripPublicKey(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	trx1 := ts.GenerateTestTransferTx()
	id1 := trx1.ID()
	assert.NoError(t, trx1.BasicCheck())

	trx1.StripPublicKey()
	assert.True(t, trx1.IsPublicKeyStriped())
	assert.Equal(t, id1, trx1.ID())
	assert.ErrorIs(t, trx1.BasicCheck(),
		tx.BasicCheckError{
			Reason: "no public key",
		})

	bs1, _ := trx1.Bytes()
	trx2, _ := tx.FromBytes(bs1)
	bs2, _ := trx2.Bytes()

	assert.Equal(t, bs1, bs2)
	assert.Equal(t, trx1.ID(), trx2.ID())
	assert.Nil(t, trx2.PublicKey())
}

func TestFlagNotSigned(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	trx := tx.NewTransferTx(ts.RandHeight(), ts.RandAccAddress(), ts.RandAccAddress(),
		ts.RandAmount(), ts.RandFee())
	assert.False(t, trx.IsSigned(), "FlagNotSigned should not be set for new transactions")

	trx.SetSignature(ts.RandBLSSignature())
	assert.True(t, trx.IsSigned(), "FlagNotSigned should be set for a signed transaction")

	trx.SetSignature(nil)
	assert.False(t, trx.IsSigned(), "FlagNotSigned should not be set when the signature is set to nil")
}

func TestInvalidSignerSignature(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	trx := tx.NewTransferTx(ts.RandHeight(), crypto.TreasuryAddress, ts.RandAccAddress(),
		ts.RandAmount(), ts.RandFee())
	trx.SetSignature(ts.RandBLSSignature())

	bytes, _ := trx.Bytes()
	_, err := tx.FromBytes(bytes)
	assert.ErrorIs(t, err, tx.ErrInvalidSigner)
}

func TestInvalidSignerPublicKey(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	trx := tx.NewTransferTx(ts.RandHeight(), crypto.TreasuryAddress, ts.RandAccAddress(),
		ts.RandAmount(), ts.RandFee())
	pub, _ := ts.RandBLSKeyPair()
	trx.SetSignature(ts.RandBLSSignature())
	trx.SetPublicKey(pub)

	bytes, _ := trx.Bytes()
	_, err := tx.FromBytes(bytes)
	assert.ErrorIs(t, err, tx.ErrInvalidSigner)
}

func TestIsFreeTx(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	trx1 := ts.GenerateTestTransferTx()
	trx2 := ts.GenerateTestBondTx()
	trx3 := ts.GenerateTestUnbondTx()
	trx4 := ts.GenerateTestWithdrawTx()
	trx5 := ts.GenerateTestSortitionTx()

	assert.True(t, trx1.IsTransferTx())
	assert.True(t, trx2.IsBondTx())
	assert.True(t, trx3.IsUnbondTx())
	assert.True(t, trx4.IsWithdrawTx())
	assert.True(t, trx5.IsSortitionTx())

	assert.False(t, trx1.IsFreeTx())
	assert.False(t, trx2.IsFreeTx())
	assert.True(t, trx3.IsFreeTx())
	assert.False(t, trx4.IsFreeTx())
	assert.True(t, trx5.IsFreeTx())
}
