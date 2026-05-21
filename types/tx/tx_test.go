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
	"github.com/pactus-project/pactus/crypto/secp256k1"
	"github.com/pactus-project/pactus/types"
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
	require.NoError(t, err)
	tx2 := new(tx.Tx)
	require.NoError(t, cbor.Unmarshal(bz, tx2))
	assert.Equal(t, tx1.ID(), tx2.ID())

	require.Error(t, cbor.Unmarshal([]byte{1}, tx2))
}

func TestEncodingTx(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	trx1 := ts.GenerateTestTransferTx()
	trx2 := ts.GenerateTestSubsidyTx()
	trx3 := ts.GenerateTestBatchTransferTx()
	trx4 := ts.GenerateTestBondTx()
	trx5 := ts.GenerateTestUnbondTx()
	trx6 := ts.GenerateTestWithdrawTx()
	trx7 := ts.GenerateTestSortitionTx()

	tests := []*tx.Tx{trx1, trx2, trx3, trx4, trx5, trx6, trx7}
	for _, trx := range tests {
		require.NoError(t, trx.BasicCheck())
		require.NoError(t, trx.BasicCheck()) // double basic check

		length := trx.SerializeSize()
		for i := 0; i < length; i++ {
			w := util.NewFixedWriter(i)
			require.Error(t, trx.Encode(w), "encode test %v failed", i)
		}
		w := util.NewFixedWriter(length)
		require.NoError(t, trx.Encode(w))

		for i := 0; i < length; i++ {
			newTrx := new(tx.Tx)
			r := util.NewFixedReader(i, w.Bytes())
			require.Error(t, newTrx.Decode(r), "decode test %v failed", i)
		}

		bz, err := trx.Bytes()
		require.NoError(t, err)
		decodedTrx, err := tx.FromBytes(bz)
		require.NoError(t, err)
		assert.Equal(t, trx.ID(), decodedTrx.ID())
	}

	_, err := tx.FromString("badcow")
	require.Error(t, err)
}

func TestTxIDNoSignatory(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	trx := ts.GenerateTestTransferTx()
	id1 := trx.ID()

	trx.SetSignature(nil)
	trx.SetPublicKey(nil)
	id2 := trx.ID()

	require.Equal(t, id1, id2)
}

func TestBasicCheck(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	t.Run("LockTime is not defined", func(t *testing.T) {
		trx := tx.NewTransferTx(0,
			ts.RandAccAddress(), ts.RandAccAddress(), ts.RandAmount(), ts.RandFee())

		err := trx.BasicCheck()
		require.ErrorIs(t, err, tx.BasicCheckError{
			Reason: "lock time is not defined",
		})
	})

	t.Run("Big memo, Should returns error", func(t *testing.T) {
		bigMemo := strings.Repeat("a", 65)

		trx := tx.NewTransferTx(ts.RandHeight(),
			ts.RandAccAddress(), ts.RandAccAddress(), ts.RandAmount(), ts.RandFee(), tx.WithMemo(bigMemo))

		err := trx.BasicCheck()
		require.ErrorIs(t, err, tx.BasicCheckError{
			Reason: "memo length exceeded: 65",
		})
	})

	t.Run("Invalid payload, Should returns error", func(t *testing.T) {
		invAddr := ts.RandAccAddress()
		invAddr[0] = 5
		trx := tx.NewTransferTx(ts.RandHeight(), ts.RandAccAddress(), invAddr, 1e9, ts.RandFee())

		err := trx.BasicCheck()
		require.ErrorIs(t, err, tx.BasicCheckError{
			Reason: "invalid payload: receiver is not an account address: " + invAddr.String(),
		})
	})

	t.Run("Invalid amount", func(t *testing.T) {
		trx := tx.NewTransferTx(ts.RandHeight(), ts.RandAccAddress(), ts.RandAccAddress(), -1, 1)

		err := trx.BasicCheck()
		require.ErrorIs(t, err, tx.BasicCheckError{
			Reason: "invalid amount: -0.000000001 PAC",
		})
	})

	t.Run("Invalid amount", func(t *testing.T) {
		trx := tx.NewTransferTx(ts.RandHeight(), ts.RandAccAddress(), ts.RandAccAddress(), (42e15)+1, 1)

		err := trx.BasicCheck()
		require.ErrorIs(t, err, tx.BasicCheckError{
			Reason: "invalid amount: 42,000,000 PAC",
		})
	})

	t.Run("Invalid signer address", func(t *testing.T) {
		valKey := ts.RandValKey()
		trx := tx.NewTransferTx(ts.RandHeight(), ts.RandAccAddress(), ts.RandAccAddress(), 1, 1)
		sig := valKey.PrivateKey().Sign(trx.SignBytes())
		trx.SetSignature(sig)
		trx.SetPublicKey(valKey.PublicKey())

		err := trx.BasicCheck()
		require.ErrorIs(t, err, tx.BasicCheckError{
			Reason: fmt.Sprintf("address mismatch: expected %s, got %s",
				valKey.PublicKey().AccountAddress(), trx.Payload().Signer()),
		})
	})

	t.Run("Invalid version", func(t *testing.T) {
		str := "02" + // Flags
			"02" + // Version
			"01020304" + // LockTime
			"01" + // Fee
			"00" + // Memo
			"01" + // PayloadType
			"00" + // Sender (treasury)
			"012222222222222222222222222222222222222222" + // Receiver
			"01" // Amount

		trx, err := tx.FromString(str)
		require.NoError(t, err)
		err = trx.BasicCheck()
		require.ErrorIs(t, err, tx.BasicCheckError{
			Reason: "invalid version: 2",
		})
	})
}

func TestInvalidPayloadType(t *testing.T) {
	str := "02" + // Flags
		"01" + // Version
		"01020300" + // LockTime
		"01" + // Fee
		"00" + // Memo
		"07" + // PayloadType
		"00" + // Sender (treasury)
		"012222222222222222222222222222222222222222" + // Receiver
		"01" // Amount

	_, err := tx.FromString(str)
	require.ErrorIs(t, err, tx.InvalidPayloadTypeError{
		PayloadType: payload.Type(7),
	})
}

func TestSubsidyTx(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pub, prv := ts.RandEd25519KeyPair()

	t.Run("Has signature", func(t *testing.T) {
		trx := ts.GenerateTestSubsidyTx()
		sig := prv.Sign(trx.SignBytes())
		trx.SetSignature(sig)

		err := trx.BasicCheck()
		require.ErrorIs(t, err, tx.BasicCheckError{
			Reason: "subsidy transaction with signatory",
		})
	})

	t.Run("Has public key", func(t *testing.T) {
		trx := ts.GenerateTestSubsidyTx()
		trx.SetPublicKey(pub)

		err := trx.BasicCheck()
		require.ErrorIs(t, err, tx.BasicCheckError{
			Reason: "subsidy transaction with signatory",
		})
	})

	t.Run("Strip public key", func(t *testing.T) {
		trx := ts.GenerateTestSubsidyTx()
		trx.StripPublicKey()

		err := trx.BasicCheck()
		require.NoError(t, err)
		assert.False(t, trx.IsPublicKeyStriped())
	})
}

func TestInvalidSignature(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	t.Run("Good", func(t *testing.T) {
		trx := ts.GenerateTestTransferTx()
		require.NoError(t, trx.BasicCheck())
	})

	t.Run("No signature", func(t *testing.T) {
		trx := ts.GenerateTestTransferTx()
		trx.SetSignature(nil)

		err := trx.BasicCheck()
		require.ErrorIs(t, err, tx.BasicCheckError{
			Reason: "no signature",
		})
	})

	t.Run("No public key", func(t *testing.T) {
		trx := ts.GenerateTestTransferTx()
		trx.SetPublicKey(nil)

		err := trx.BasicCheck()
		require.ErrorIs(t, err, tx.BasicCheckError{
			Reason: "no public key",
		})
	})

	pbInv, pvInv := ts.RandBLSKeyPair()
	t.Run("Invalid signature", func(t *testing.T) {
		trx := ts.GenerateTestTransferTx()
		sig := pvInv.Sign(trx.SignBytes())
		trx.SetSignature(sig)

		err := trx.BasicCheck()
		require.ErrorIs(t, err, tx.BasicCheckError{
			Reason: "invalid signature",
		})
	})

	t.Run("Invalid public key", func(t *testing.T) {
		trx := ts.GenerateTestTransferTx()
		trx.SetPublicKey(pbInv)

		err := trx.BasicCheck()
		require.ErrorIs(t, err, tx.BasicCheckError{
			Reason: fmt.Sprintf("address mismatch: expected %s, got %s", pbInv.AccountAddress(), trx.Payload().Signer()),
		})
	})

	t.Run("Invalid sign Bytes", func(t *testing.T) {
		valKey := ts.RandValKey()
		trx0 := ts.GenerateTestUnbondTx(testsuite.TransactionWithSigner(valKey.PrivateKey()))

		trx := tx.NewUnbondTx(trx0.LockTime(), valKey.Address(), tx.WithMemo("invalidate signature"))
		trx.SetPublicKey(trx0.PublicKey())
		trx.SetSignature(trx0.Signature())

		err := trx.BasicCheck()
		require.ErrorIs(t, err, tx.BasicCheckError{
			Reason: "invalid signature",
		})
	})

	t.Run("Zero signature", func(t *testing.T) {
		trx := ts.GenerateTestTransferTx()
		trx.SetSignature(&bls.Signature{})

		err := trx.BasicCheck()
		require.ErrorIs(t, err, tx.BasicCheckError{
			Reason: "invalid signature",
		})
	})

	t.Run("Zero public key", func(t *testing.T) {
		trx := ts.GenerateTestTransferTx()
		zeroPubKey := &bls.PublicKey{}
		trx.SetPublicKey(zeroPubKey)

		err := trx.BasicCheck()
		require.ErrorIs(t, err, tx.BasicCheckError{
			Reason: fmt.Sprintf("address mismatch: expected %s, got %s",
				zeroPubKey.AccountAddress().String(), trx.Payload().Signer()),
		})
	})
}

func TestSignBytesBLS(t *testing.T) {
	data, _ := hex.DecodeString(
		"00" + // Flags
			"01" + // Version
			"01020300" + // LockTime
			"e807" + // Fee
			"0474657374" + // Memo ("test")
			"01" + // PayloadType
			"02da21c67cdd20af88deeb89b95ac5ce771b8be228" + // Sender (pc1zmgsuvlxayzhc3hht3xu443wwwudchc3gwtenu2)
			"02ef6058d7e8863bad5f149d31684c12d67560f257" + // Receiver (pc1zaas934lgsca66hc5n5cksnqj6e6kpujhvtvtlj)
			"a09c01" + // Amount
			"97793eda4277d44baa090af28fe649beb9666546b1cb56302befdff61ee4afa2d7d4635019fefc7d5a68a50430213beb" + // Signature
			"94a8a0d79f2db62b84123eac8fd057c2d212c8a314468420730cf99bf19fa718807ee6e94a7e17681b085d970d6fb89c" + // PublicKey
			"190e5209ee2b3138176c3964fd1b0187f5446ba8d692d691e508aa790cf6fbcee02b61ead97bae7b65067b3b485bb740")

	txID, _ := hash.FromString("a57ae2db8d361504b6b97803cdf2bc27561fde536f49683cdbcbb23ceaded0ae")
	trx, err := tx.FromBytes(data)
	require.NoError(t, err)
	require.NoError(t, trx.BasicCheck())
	assert.Equal(t, len(data), trx.SerializeSize())

	signBytes := data[1 : len(data)-bls.PublicKeySize-bls.SignatureSize]
	assert.Equal(t, signBytes, trx.SignBytes())
	assert.Equal(t, hash.CalcHash(signBytes), trx.ID())
	assert.Equal(t, txID, trx.ID())
	assert.Equal(t, types.Height(0x30201), trx.LockTime())
	assert.Equal(t, "test", trx.Memo())
	assert.Equal(t, amount.Amount(1000), trx.Fee())
	assert.Equal(t, amount.Amount(20000), trx.Payload().Value())
	assert.Equal(t, "pc1zmgsuvlxayzhc3hht3xu443wwwudchc3gwtenu2", trx.Payload().Signer().String())
}

func TestSignBytesEd25519(t *testing.T) {
	data, _ := hex.DecodeString(
		"00" + // Flags
			"01" + // Version
			"01020300" + // LockTime
			"e807" + // Fee (1000)
			"0474657374" + // Memo ("test")
			"01" + // PayloadType
			"03706ce1d84f2c03c543dfecd6adb96bfa653fe560" + // Sender (pc1rwpkwrkz09spu2s7lant2mwttlfjnletqkqnlpq)
			"03088868ebc9ec0528a46000bc7dae285eb65a0ef0" + // Receiver (pc1rpzyx367faszj3frqqz78mt3gt6m95rhsvavwd2)
			"a09c01" + // Amount (20000)
			"109a480d296f61e793b64a9ffcf1ce4a5a535c882edfa9c52aa0af8003e39cee" + // Signature
			"da02d5c5ce69c135514bbd1f4c6c42798fe2b9615e56ebef93ad2923b33f170e" + // PublicKey
			"ade295295be1368479ecf0345c13ea0db1248d4dffe39eac627319954d39b264")

	txID, _ := hash.FromString("5f8ea503759a69a6850a2cbdafb60ce769cdf85b38f14848d4120c5c984a33fd")
	trx, err := tx.FromBytes(data)
	require.NoError(t, err)
	require.NoError(t, trx.BasicCheck())
	assert.Equal(t, len(data), trx.SerializeSize())

	signBytes := data[1 : len(data)-ed25519.PublicKeySize-ed25519.SignatureSize]
	assert.Equal(t, signBytes, trx.SignBytes())
	assert.Equal(t, hash.CalcHash(signBytes), trx.ID())
	assert.Equal(t, txID, trx.ID())
	assert.Equal(t, types.Height(0x00030201), trx.LockTime())
	assert.Equal(t, "test", trx.Memo())
	assert.Equal(t, amount.Amount(1000), trx.Fee())
	assert.Equal(t, amount.Amount(20000), trx.Payload().Value())
	assert.Equal(t, "pc1rwpkwrkz09spu2s7lant2mwttlfjnletqkqnlpq", trx.Payload().Signer().String())
}

func TestSignBytesSecp256k1(t *testing.T) {
	data, _ := hex.DecodeString(
		"00" + // Flags
			"01" + // Version
			"01020300" + // LockTime
			"e807" + // Fee (1000)
			"0474657374" + // Memo ("test")
			"01" + // PayloadType
			"042ff976b92d83950db9e1b2fe01599a685939a59b" + // Sender (pc1y9luhdwfdsw2smw0pktlqzkv6dpvnnfvmcu4qmp)
			"04847648ef3454f9d947cef12c76e57e8c92ee3962" + // Receiver (pc1ys3my3me52nuaj37w7yk8det73jfwuwtz7klsty)
			"a09c01" + // Amount (20000)
			"aa25da26f0ba3bc95b32f4fc5565440c41485437e61a095652a2d828ad94a254" +
			"20c355314a7d020bbbadc443d51b1f6a8397bee5827ecace1785aa4d0c8684f7" + // Signature
			"0336bd9ebe4b518fec7f424ee2e6dc600f4cb22f10f7c4e4e466af374bfe39b870") // PublicKey

	txID, _ := hash.FromString("ac3bb060ed6fd4c12eb9f272e55b59231956c3276ee1a9287c1b05fe7928e2b5")
	trx, err := tx.FromBytes(data)
	require.NoError(t, trx.BasicCheck())
	require.NoError(t, err)
	assert.Equal(t, len(data), trx.SerializeSize())

	signBytes := data[1 : len(data)-secp256k1.PublicKeySize-secp256k1.SignatureSize]
	assert.Equal(t, signBytes, trx.SignBytes())
	assert.Equal(t, hash.CalcHash(signBytes), trx.ID())
	assert.Equal(t, txID, trx.ID())
	assert.Equal(t, types.Height(0x00030201), trx.LockTime())
	assert.Equal(t, "test", trx.Memo())
	assert.Equal(t, amount.Amount(1000), trx.Fee())
	assert.Equal(t, amount.Amount(20000), trx.Payload().Value())
	assert.Equal(t, "pc1y9luhdwfdsw2smw0pktlqzkv6dpvnnfvmcu4qmp", trx.Payload().Signer().String())
}

func TestStripPublicKey(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	trx1 := ts.GenerateTestTransferTx()
	id1 := trx1.ID()
	require.NoError(t, trx1.BasicCheck())

	trx1.StripPublicKey()
	assert.True(t, trx1.IsPublicKeyStriped())
	assert.Equal(t, id1, trx1.ID())
	require.ErrorIs(t, trx1.BasicCheck(),
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
	require.ErrorIs(t, err, tx.ErrInvalidSigner)
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
	require.ErrorIs(t, err, tx.ErrInvalidSigner)
}

func TestIsFreeTx(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	trx1 := ts.GenerateTestTransferTx()
	trx2 := ts.GenerateTestBatchTransferTx()
	trx3 := ts.GenerateTestBondTx()
	trx4 := ts.GenerateTestUnbondTx()
	trx5 := ts.GenerateTestWithdrawTx()
	trx6 := ts.GenerateTestSortitionTx()

	assert.True(t, trx1.IsTransferTx())
	assert.True(t, trx2.IsBatchTransferTx())
	assert.True(t, trx3.IsBondTx())
	assert.True(t, trx4.IsUnbondTx())
	assert.True(t, trx5.IsWithdrawTx())
	assert.True(t, trx6.IsSortitionTx())

	assert.False(t, trx1.IsFreeTx())
	assert.False(t, trx2.IsFreeTx())
	assert.False(t, trx3.IsFreeTx())
	assert.True(t, trx4.IsFreeTx())
	assert.False(t, trx5.IsFreeTx())
	assert.True(t, trx6.IsFreeTx())
}

func TestCheckFee(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	tests := []struct {
		name        string
		trx         *tx.Tx
		expectedErr error
	}{
		{
			name: "Negative fee",
			trx: ts.GenerateTestTransferTx(
				testsuite.TransactionWithFee(-1)),
			expectedErr: tx.BasicCheckError{Reason: "invalid fee: -0.000000001 PAC"},
		},
		{
			name: "Big fee",
			trx: ts.GenerateTestTransferTx(
				testsuite.TransactionWithFee(42e15 + 1)),
			expectedErr: tx.BasicCheckError{Reason: "invalid fee: 42,000,000 PAC"},
		},
		{
			name:        "Subsidy transaction with fee",
			trx:         tx.NewTransferTx(ts.RandHeight(), crypto.TreasuryAddress, ts.RandAccAddress(), ts.RandAmount(), 1),
			expectedErr: tx.BasicCheckError{Reason: "invalid fee: 0.000000001 PAC"},
		},
		{
			name:        "Subsidy transaction with zero fee",
			trx:         tx.NewTransferTx(ts.RandHeight(), crypto.TreasuryAddress, ts.RandAccAddress(), ts.RandAmount(), 0),
			expectedErr: nil,
		},
		{
			name: "Transfer transaction with zero fee",
			trx: ts.GenerateTestTransferTx(
				testsuite.TransactionWithFee(0)),
			expectedErr: nil,
		},
		{
			name: "Transfer transaction with non-zero fee",
			trx: ts.GenerateTestTransferTx(
				testsuite.TransactionWithFee(1)),
			expectedErr: nil,
		},
		{
			name: "Unbond transaction with zero fee",
			trx: ts.GenerateTestUnbondTx(
				testsuite.TransactionWithFee(0)),
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.trx.BasicCheck()
			require.ErrorIs(t, err, tt.expectedErr)
		})
	}
}
