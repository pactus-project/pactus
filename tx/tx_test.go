package tx

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/zarbchain/zarb-go/tx/payload"

	"github.com/zarbchain/zarb-go/crypto"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestSendEncodingTx(t *testing.T) {
	tx, _ := GenerateTestSendTx()

	bz, err := tx.MarshalCBOR()
	require.NoError(t, err)
	var tx2 Tx
	err = tx2.UnmarshalCBOR(bz)

	bz2, _ := tx2.MarshalCBOR()

	fmt.Printf("%x\n", bz)

	require.NoError(t, err)
	require.Equal(t, bz, bz2)
	require.Equal(t, tx.Hash(), tx2.Hash())
}

func TestBondEncodingTx(t *testing.T) {
	tx, _ := GenerateTestBondTx()

	bz, err := tx.MarshalCBOR()
	require.NoError(t, err)
	var tx2 Tx
	err = tx2.UnmarshalCBOR(bz)

	bz2, _ := tx2.MarshalCBOR()

	fmt.Printf("%x\n", bz)

	require.NoError(t, err)
	require.Equal(t, bz, bz2)
	require.Equal(t, tx.Hash(), tx2.Hash())
}

func TestEncodingTxNoMemo(t *testing.T) {
	tx, _ := GenerateTestSendTx()
	tx.data.Memo = ""

	bz, err := tx.MarshalCBOR()
	require.NoError(t, err)
	var tx2 Tx
	err = tx2.UnmarshalCBOR(bz)

	bz2, _ := tx2.MarshalCBOR()

	require.NoError(t, err)
	require.Equal(t, bz, bz2)
	require.Equal(t, tx.Hash(), tx2.Hash())
}

func TestEncodingTxNoSig(t *testing.T) {
	tx, _ := GenerateTestSendTx()
	tx.SetPublicKey(nil)
	tx.SetSignature(nil)

	bz, err := tx.MarshalCBOR()
	require.NoError(t, err)
	var tx2 Tx
	err = tx2.UnmarshalCBOR(bz)

	bz2, _ := tx2.MarshalCBOR()

	fmt.Printf("%x\n", bz)
	fmt.Printf("%x\n", bz2)

	require.NoError(t, err)
	require.Equal(t, bz, bz2)
	require.Equal(t, tx.Hash(), tx2.Hash())
}

func TestTxSanityCheck(t *testing.T) {

	t.Run("Invalid version", func(t *testing.T) {
		tx, priv := GenerateTestSendTx()
		tx.data.Version = 2
		sig := priv.Sign(tx.SignBytes())
		tx.data.Signature = sig
		assert.Error(t, tx.SanityCheck())
	})

	t.Run("Invalid sequence", func(t *testing.T) {
		tx, priv := GenerateTestSendTx()
		tx.data.Sequence = -1
		sig := priv.Sign(tx.SignBytes())
		tx.data.Signature = sig
		assert.Error(t, tx.SanityCheck())
	})

	t.Run("Invalid payload type", func(t *testing.T) {
		tx, priv := GenerateTestSendTx()
		tx.data.Type = 2
		sig := priv.Sign(tx.SignBytes())
		tx.data.Signature = sig
		assert.Error(t, tx.SanityCheck())
	})
}

func TestMintbaseTx(t *testing.T) {
	a, pub, priv := crypto.GenerateTestKeyPair()
	trx := NewMintbaseTx(crypto.GenerateTestHash(), 111, a, 1111, "mintbase")

	trx.data.Fee = 1
	assert.Error(t, trx.SanityCheck())

	sig := priv.Sign(trx.SignBytes())
	trx.SetSignature(sig)
	assert.Error(t, trx.SanityCheck())

	trx.SetPublicKey(&pub)
	assert.Error(t, trx.SanityCheck())
}

func TestInvalidSignature(t *testing.T) {
	tx, pv := GenerateTestSendTx()
	assert.NoError(t, tx.SanityCheck())

	fmt.Printf("%x\n", tx.SignBytes())

	tx.SetSignature(nil)
	assert.Error(t, tx.SanityCheck())

	tx.SetPublicKey(nil)
	assert.Error(t, tx.SanityCheck())

	_, pbInv, pvInv := crypto.GenerateTestKeyPair()
	tx.SetPublicKey(&pbInv)
	assert.Error(t, tx.SanityCheck())

	sig := pvInv.Sign(tx.SignBytes())
	tx.SetSignature(sig)
	assert.Error(t, tx.SanityCheck())

	// Invalid sign Bytes
	var tx2 = new(Tx)
	tx2.data.Memo = "Hack me"
	sig = pv.Sign(tx2.SignBytes())
	pb := pv.PublicKey()
	tx.SetPublicKey(&pb)
	tx.SetSignature(sig)
	assert.Error(t, tx.SanityCheck())
}

func TestSendSanityCheck(t *testing.T) {
	invAddr, _, _ := crypto.GenerateTestKeyPair()
	t.Run("Ok", func(t *testing.T) {
		trx, _ := GenerateTestSendTx()
		assert.NoError(t, trx.SanityCheck())
	})

	t.Run("Invalid amount", func(t *testing.T) {
		trx, priv := GenerateTestSendTx()
		pld := trx.data.Payload.(*payload.SendPayload)
		pld.Amount = -1
		trx.SetSignature(priv.Sign(trx.SignBytes()))
		assert.Error(t, trx.SanityCheck())
	})

	t.Run("Invalid sender", func(t *testing.T) {
		trx, priv := GenerateTestSendTx()
		pld := trx.data.Payload.(*payload.SendPayload)
		pld.Sender = invAddr
		trx.SetSignature(priv.Sign(trx.SignBytes()))
		assert.Error(t, trx.SanityCheck())
	})

	t.Run("Invalid receiver", func(t *testing.T) {
		trx, priv := GenerateTestSendTx()
		pld := trx.data.Payload.(*payload.SendPayload)
		pld.Receiver = crypto.MintbaseAddress
		trx.SetSignature(priv.Sign(trx.SignBytes()))
		assert.Error(t, trx.SanityCheck())
	})
}

func TestBondSanityCheck(t *testing.T) {
	invAddr, _, _ := crypto.GenerateTestKeyPair()
	t.Run("Ok", func(t *testing.T) {
		trx, _ := GenerateTestBondTx()
		assert.NoError(t, trx.SanityCheck())
	})

	t.Run("Invalid stake", func(t *testing.T) {
		trx, priv := GenerateTestBondTx()
		pld := trx.data.Payload.(*payload.BondPayload)
		pld.Stake = -1
		trx.SetSignature(priv.Sign(trx.SignBytes()))
		assert.Error(t, trx.SanityCheck())
	})

	t.Run("Invalid bonder", func(t *testing.T) {
		trx, priv := GenerateTestBondTx()
		pld := trx.data.Payload.(*payload.BondPayload)
		pld.Bonder = invAddr
		trx.SetSignature(priv.Sign(trx.SignBytes()))
		assert.Error(t, trx.SanityCheck())
	})
}

func TestSortitionSanityCheck(t *testing.T) {
	invAddr, _, _ := crypto.GenerateTestKeyPair()
	t.Run("Ok", func(t *testing.T) {
		trx, _ := GenerateTestSortitionTx()
		assert.NoError(t, trx.SanityCheck())
	})

	t.Run("Invalid address", func(t *testing.T) {
		trx, priv := GenerateTestSortitionTx()
		pld := trx.data.Payload.(*payload.SortitionPayload)
		pld.Address = invAddr
		trx.SetSignature(priv.Sign(trx.SignBytes()))
		assert.Error(t, trx.SanityCheck())
	})

	t.Run("Invalid proof", func(t *testing.T) {
		trx, priv := GenerateTestSortitionTx()
		pld := trx.data.Payload.(*payload.SortitionPayload)
		pld.Proof = nil
		trx.SetSignature(priv.Sign(trx.SignBytes()))
		assert.Error(t, trx.SanityCheck())
	})

}

func TestSendDecodingAndHash(t *testing.T) {
	d, _ := hex.DecodeString("a9010102582008f7d9c21fdaa4a4147e60a0f3933c850b0c0d9af6b2a308c0a7b5639a7e49d603186e040a050106a301548dfaf698d3889b13251529ff971277305fbf1f440254bdd1540a13d82c38e5b4dfbd7b5b2bcab5fb5f290318640767746573742074781458603846ed5d519e51f6dd63e552ac410c531d5436c726475f6f8fb51c1133b07e32bd3bc4c674359546a1145cb1935a3c0621fc5329c6039707445e472a73d857d8eff832b971838b21c53baa090d90b02f6d2c5a1e9358e46f4f4955ff737c08991558309e6b99c60cf5ccb551efc793ec2bedd66070bde8fbddeb8305c5f670a21532304775637f9e622d5212f34c9479d12d11")
	s, _ := hex.DecodeString("a7010102582008f7d9c21fdaa4a4147e60a0f3933c850b0c0d9af6b2a308c0a7b5639a7e49d603186e040a050106a301548dfaf698d3889b13251529ff971277305fbf1f440254bdd1540a13d82c38e5b4dfbd7b5b2bcab5fb5f29031864076774657374207478")
	h, _ := crypto.HashFromString("a89c828d859a00677a5f8d58425c257a9ae4567b86d8029981272a5e42c3b90b")
	var trx Tx
	err := trx.UnmarshalCBOR(d)
	assert.NoError(t, err)
	d2, _ := trx.MarshalCBOR()
	assert.Equal(t, d, d2)
	assert.Equal(t, trx.SignBytes(), s)
	assert.Equal(t, trx.Hash(), h)
}

func TestBondDecodingAndHash(t *testing.T) {
	d, _ := hex.DecodeString("a90101025820401a78337d2715db2a69916bbedbbb0a44336915fa1e5938b5b9cf350c340bb603186e0400050206a30154ebd6d0afdbe2ac5968224763f1c2dbb99fedda53025860c0ac961b453008e5cf521afab81fe303dc957de4b385bf472108e27b054f14e2da340f288928472b443a716eb795800cd8b885791827debda0a7909c741286cebd324db9546bab8af03830a9f91792389303b3746d76fc57148830184c3e0506031864076c7465737420626f6e642d74781458603b42c95589d73085882cf972c80b8c652de6e89ae4268a37a6c2d3206019bf9b706e048a7555e77f926895f6dccdb10f42efcaf0008006a0a15b8c58d704a9a0f277c66d7c6f6e8be1efb50b509c32b3b4364b2b7feb7eda106164c7ad18408a1558305a53c9e788896943587b594a2c2b930f2f9945abdf12e1c16fae8d235f74663a2f28d1b8fa442af0781f13f7de5a1c99")
	s, _ := hex.DecodeString("a70101025820401a78337d2715db2a69916bbedbbb0a44336915fa1e5938b5b9cf350c340bb603186e0400050206a30154ebd6d0afdbe2ac5968224763f1c2dbb99fedda53025860c0ac961b453008e5cf521afab81fe303dc957de4b385bf472108e27b054f14e2da340f288928472b443a716eb795800cd8b885791827debda0a7909c741286cebd324db9546bab8af03830a9f91792389303b3746d76fc57148830184c3e0506031864076c7465737420626f6e642d7478")
	h, _ := crypto.HashFromString("a428e98f0df5dff7187fce047f9c463e9e7a013de155d2e2cfdfb3acdb286e2d")
	var trx Tx
	err := trx.UnmarshalCBOR(d)
	assert.NoError(t, err)
	d2, _ := trx.MarshalCBOR()
	assert.Equal(t, d, d2)
	assert.Equal(t, trx.SignBytes(), s)
	assert.Equal(t, trx.Hash(), h)
}
