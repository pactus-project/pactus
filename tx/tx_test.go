package tx

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/sortition"
	"github.com/zarbchain/zarb-go/tx/payload"
)

func TestJSONMarshaling(t *testing.T) {
	tx, _ := GenerateTestSendTx()
	_, err := tx.MarshalJSON()
	require.NoError(t, err)
}

func TestSendEncodingTx(t *testing.T) {
	tx, _ := GenerateTestSendTx()
	bz, err := tx.Encode()
	require.NoError(t, err)
	var tx2 Tx
	require.NoError(t, tx2.Decode(bz))
	require.Equal(t, tx.ID(), tx2.ID())
}

func TestBondEncodingTx(t *testing.T) {
	tx, _ := GenerateTestBondTx()
	bz, err := tx.MarshalCBOR()
	require.NoError(t, err)
	var tx2 Tx
	require.NoError(t, tx2.UnmarshalCBOR(bz))
	require.Equal(t, tx.ID(), tx2.ID())
}

func TestEncodingTxNoSig(t *testing.T) {
	tx, _ := GenerateTestSendTx()
	bz, _ := tx.MarshalCBOR()
	fmt.Printf("%x\n", bz)
	tx.data.Signature = nil
	tx.data.PublicKey = nil
	bz, err := tx.MarshalCBOR()
	require.NoError(t, err)
	var tx2 Tx
	require.NoError(t, tx2.UnmarshalCBOR(bz))
	require.Equal(t, tx.ID(), tx2.ID())
}

func TestTxSanityCheck(t *testing.T) {
	t.Run("Invalid sequence", func(t *testing.T) {
		tx, signer := GenerateTestSendTx()
		tx.data.Sequence = -1
		signer.SignMsg(tx)
		assert.Error(t, tx.SanityCheck())
	})

	t.Run("Invalid payload type", func(t *testing.T) {
		tx, signer := GenerateTestSendTx()
		tx.data.Type = 2
		signer.SignMsg(tx)
		assert.Error(t, tx.SanityCheck())
	})

	t.Run("Transaction ID should be same for signed and unsigned transactions", func(t *testing.T) {
		tx, _ := GenerateTestSendTx()
		id1 := tx.ID()
		sb1 := tx.SignBytes()
		tx.data.PublicKey = nil
		tx.data.Signature = nil
		id2 := tx.ID()
		sb2 := tx.SignBytes()
		assert.Equal(t, id1, id2)
		assert.Equal(t, sb1, sb2)
	})
}

func TestSubsidyTx(t *testing.T) {
	pub, prv := bls.GenerateTestKeyPair()
	t.Run("Invalid fee", func(t *testing.T) {
		stamp := hash.GenerateTestStamp()
		trx := NewMintbaseTx(stamp, 88, pub.Address(), 2500, "subsidy")
		assert.True(t, trx.IsMintbaseTx())
		trx.data.Fee = 1
		assert.Error(t, trx.SanityCheck())
	})

	t.Run("Has signature", func(t *testing.T) {
		stamp := hash.GenerateTestStamp()
		trx := NewMintbaseTx(stamp, 88, pub.Address(), 2500, "subsidy")
		sig := prv.Sign(trx.SignBytes())
		trx.SetSignature(sig)
		assert.Error(t, trx.SanityCheck())
	})

	t.Run("Has public key", func(t *testing.T) {
		stamp := hash.GenerateTestStamp()
		trx := NewMintbaseTx(stamp, 88, pub.Address(), 2500, "subsidy")
		trx.SetPublicKey(pub)
		assert.Error(t, trx.SanityCheck())
	})
}

func TestInvalidSignature(t *testing.T) {
	t.Run("Good", func(t *testing.T) {
		tx, _ := GenerateTestSendTx()
		assert.NoError(t, tx.SanityCheck())
	})

	t.Run("No signature", func(t *testing.T) {
		tx, _ := GenerateTestSendTx()
		tx.data.Signature = nil
		assert.Error(t, tx.SanityCheck())
	})

	t.Run("No public key", func(t *testing.T) {
		tx, _ := GenerateTestSendTx()
		tx.data.PublicKey = nil
		assert.Error(t, tx.SanityCheck())
	})

	pbInv, pvInv := bls.GenerateTestKeyPair()
	t.Run("Invalid signature", func(t *testing.T) {
		tx, _ := GenerateTestSendTx()
		sig := pvInv.Sign(tx.SignBytes())
		tx.SetSignature(sig)
		assert.Error(t, tx.SanityCheck())
	})

	t.Run("Invalid public key", func(t *testing.T) {
		tx, _ := GenerateTestSendTx()
		tx.SetPublicKey(pbInv)
		assert.Error(t, tx.SanityCheck())

	})

	t.Run("Invalid sign Bytes", func(t *testing.T) {
		tx, signer := GenerateTestSendTx()
		var tx2 = new(Tx)
		tx2.data.Memo = "Hello"
		tx.SetSignature(signer.SignData(tx2.SignBytes()))
		assert.Error(t, tx.SanityCheck())
	})
}

func TestSendSanityCheck(t *testing.T) {
	invAddr := crypto.GenerateTestAddress()
	t.Run("Ok", func(t *testing.T) {
		trx, _ := GenerateTestSendTx()
		assert.NoError(t, trx.SanityCheck())
	})

	t.Run("Invalid amount", func(t *testing.T) {
		trx, signer := GenerateTestSendTx()
		pld := trx.data.Payload.(*payload.SendPayload)
		pld.Amount = -1
		signer.SignMsg(trx)
		assert.Error(t, trx.SanityCheck())
	})

	t.Run("Invalid fee", func(t *testing.T) {
		trx, signer := GenerateTestSendTx()
		trx.data.Fee = 0
		signer.SignMsg(trx)
		assert.Error(t, trx.SanityCheck())
	})

	t.Run("Invalid sender", func(t *testing.T) {
		trx, signer := GenerateTestSendTx()
		pld := trx.data.Payload.(*payload.SendPayload)
		pld.Sender = invAddr
		signer.SignMsg(trx)
		assert.Error(t, trx.SanityCheck())
	})

	t.Run("Invalid receiver", func(t *testing.T) {
		trx, signer := GenerateTestSendTx()
		pld := trx.data.Payload.(*payload.SendPayload)
		pld.Receiver = crypto.TreasuryAddress
		signer.SignMsg(trx)
		assert.Error(t, trx.SanityCheck())
	})
}

func TestBondSanityCheck(t *testing.T) {
	invAddr := crypto.GenerateTestAddress()
	t.Run("Ok", func(t *testing.T) {
		trx, _ := GenerateTestBondTx()
		assert.NoError(t, trx.SanityCheck())
	})

	t.Run("Invalid stake", func(t *testing.T) {
		trx, signer := GenerateTestBondTx()
		pld := trx.data.Payload.(*payload.BondPayload)
		pld.Stake = -1
		signer.SignMsg(trx)
		assert.Error(t, trx.SanityCheck())
	})

	t.Run("Invalid sender", func(t *testing.T) {
		trx, signer := GenerateTestBondTx()
		pld := trx.data.Payload.(*payload.BondPayload)
		pld.Sender = invAddr
		signer.SignMsg(trx)
		assert.Error(t, trx.SanityCheck())
	})
}

func TestSortitionSanityCheck(t *testing.T) {
	invAddr := crypto.GenerateTestAddress()
	t.Run("Ok", func(t *testing.T) {
		trx, _ := GenerateTestSortitionTx()
		assert.NoError(t, trx.SanityCheck())
	})

	t.Run("Invalid address", func(t *testing.T) {
		trx, signer := GenerateTestSortitionTx()
		pld := trx.data.Payload.(*payload.SortitionPayload)
		pld.Address = invAddr
		signer.SignMsg(trx)
		assert.Error(t, trx.SanityCheck())
	})

	t.Run("Invalid fee", func(t *testing.T) {
		trx, signer := GenerateTestSortitionTx()
		trx.data.Fee = 1
		signer.SignMsg(trx)
		assert.Error(t, trx.SanityCheck())
	})

}

func TestSendDecodingAndHash(t *testing.T) {
	d, _ := hex.DecodeString("a901010244e4f59ccd03186e041903e80501065833a3015501d75c059a4157d78f9b86741164037392de0fa53102550194f782f332649a4234b79216277e0b1594836313031903e8076c746573742073656e642d7478145860a4de42541ddeebfa6c4c8f008d2a64e6a2c8069096a5ad2fd807089a2f3ca8b71554365a01a2a3d5eee73f814b2aaeee0a49496e9222bc5cb4e9ffec219b4dca5091844ac1752286a524ca89928187ea60d0bdd6f10047d06f204bac5c215967155830b1c1b312df0ac1877c8daeb35eaf53c5008fb1de9654c698bab851b73d8730204c5c93c13c7d5d6b29ee439d1bdb7118")
	s, _ := hex.DecodeString("a701010244e4f59ccd03186e041903e80501065833a3015501d75c059a4157d78f9b86741164037392de0fa53102550194f782f332649a4234b79216277e0b1594836313031903e8076c746573742073656e642d7478")
	h, _ := hash.FromString("2177a8040c435ee56601f1fb2c4b9fd97adc581a153e8e8ffed80ff3b0743e57")
	var trx Tx
	err := trx.UnmarshalCBOR(d)
	assert.NoError(t, err)
	d2, _ := trx.MarshalCBOR()
	assert.Equal(t, d, d2)
	assert.Equal(t, trx.SignBytes(), s)
	assert.Equal(t, trx.ID(), h)
	assert.Equal(t, trx.PayloadType(), payload.PayloadTypeSend)
}

func TestBondDecodingAndHash(t *testing.T) {
	d, _ := hex.DecodeString("a90101024483c5d2eb03186e041903e8050206587fa30155010a59687ee8bfc4f2784a46cbc6426676aa5d456f025860a8b2e5841d1ae408ac460fa97350457603588d619db0cb515f933387745107317a1f2d25b24ace665010adf0c21310030b685c773c59a19092cb5780ba2b755aceee2b1478fe64c0a807d7d271ba0c3cc559df951773928bcb61cc8bd67332a9031903e8076c7465737420626f6e642d7478145860b6c76ad57056d913059710eb74d6dce8e446d782cd2d30fdc774be910333cfd10bd954eba4c08bcefdf5bc11857735e00504188cc1bf5e08e5e11ee7688d55096c326af829da00b449311ca1b754cb0575a1da45b3f2737150e24e4810ea077c1558308e3d6c3ee4a39f7c3d7ce3dcd170442604b9c203e70d836886533a43fec733cf3b5068664dbae0e965bdbfcdb25c9c0a")
	s, _ := hex.DecodeString("a70101024483c5d2eb03186e041903e8050206587fa30155010a59687ee8bfc4f2784a46cbc6426676aa5d456f025860a8b2e5841d1ae408ac460fa97350457603588d619db0cb515f933387745107317a1f2d25b24ace665010adf0c21310030b685c773c59a19092cb5780ba2b755aceee2b1478fe64c0a807d7d271ba0c3cc559df951773928bcb61cc8bd67332a9031903e8076c7465737420626f6e642d7478")
	h, _ := hash.FromString("ff433a1704dbbb35caaeaffc8a000b4cea721974471a9edbe94a2607865246f8")
	var trx Tx
	err := trx.UnmarshalCBOR(d)
	assert.NoError(t, err)
	d2, _ := trx.MarshalCBOR()
	assert.Equal(t, d, d2)
	assert.Equal(t, trx.SignBytes(), s)
	assert.Equal(t, trx.ID(), h)
	assert.Equal(t, trx.PayloadType(), payload.PayloadTypeBond)
}

func TestSortitionDecodingAndHash(t *testing.T) {
	d, _ := hex.DecodeString("a801010244d712c7f503186e0400050306584ba2015501ccb2131a6465585355e952ed8fe1760b4c2dc3620258306c35ee9c9b5827cffc0623fcc312febf3745eb4a351517a66dbe88c0f11dd25e83456c0ee57b14553f466747d76a7b3e145860a48b996abd267241980aa238b4ef4da7ce39896694b78710c493ef50eb0a8730672a4857a77532c65118e95f08f1671c06cab68e66f390d2b19aa3f7ad471662f30a87a146392e96b5636eaeb444fcfd320f7a46c767a1e7000cf452b212d468155830a8e667bbc8d9a53934c314ecc5c3ae3e2aab6684c7ddd7df60b2d4d27a37fa62d806fccf7d665385e6cd1c4e425be94d")
	s, _ := hex.DecodeString("a601010244d712c7f503186e0400050306584ba2015501ccb2131a6465585355e952ed8fe1760b4c2dc3620258306c35ee9c9b5827cffc0623fcc312febf3745eb4a351517a66dbe88c0f11dd25e83456c0ee57b14553f466747d76a7b3e")
	h, _ := hash.FromString("a7525eea0411c6c60d83b90b79777c4192a8648bdd7a7bda0c1ea710c6a426c3")
	var trx Tx
	err := trx.UnmarshalCBOR(d)
	assert.NoError(t, err)
	d2, _ := trx.MarshalCBOR()
	assert.Equal(t, d, d2)
	assert.Equal(t, trx.SignBytes(), s)
	assert.Equal(t, trx.ID(), h)
	assert.Equal(t, trx.PayloadType(), payload.PayloadTypeSortition)
}

func TestSendSignBytes(t *testing.T) {
	stamp := hash.GenerateTestStamp()
	signer := bls.GenerateTestSigner()
	addr := crypto.GenerateTestAddress()

	trx1 := NewSendTx(stamp, 1, signer.Address(), addr, 100, 10, "test send-tx")
	signer.SignMsg(trx1)

	trx2 := NewSendTx(stamp, 1, signer.Address(), addr, 100, 10, "test send-tx")
	trx3 := NewSendTx(stamp, 2, signer.Address(), addr, 100, 10, "test send-tx")

	assert.Equal(t, trx1.SignBytes(), trx2.SignBytes())
	assert.NotEqual(t, trx1.SignBytes(), trx3.SignBytes())
}

func TestBondSignBytes(t *testing.T) {
	stamp := hash.GenerateTestStamp()
	signer := bls.GenerateTestSigner()
	pub, _ := bls.GenerateTestKeyPair()

	trx1 := NewBondTx(stamp, 1, signer.Address(), pub, 100, 100, "test bond-tx")
	signer.SignMsg(trx1)

	trx2 := NewBondTx(stamp, 1, signer.Address(), pub, 100, 100, "test bond-tx")
	trx3 := NewBondTx(stamp, 2, signer.Address(), pub, 100, 100, "test bond-tx")

	assert.Equal(t, trx1.SignBytes(), trx2.SignBytes())
	assert.NotEqual(t, trx1.SignBytes(), trx3.SignBytes())
	assert.True(t, trx1.IsBondTx())
}

func TestUnbondSignBytes(t *testing.T) {
	stamp := hash.GenerateTestStamp()
	signer := bls.GenerateTestSigner()

	trx1 := NewUnbondTx(stamp, 1, signer.Address(), "test unbond-tx")
	signer.SignMsg(trx1)

	trx2 := NewUnbondTx(stamp, 1, signer.Address(), "test unbond-tx")
	trx3 := NewUnbondTx(stamp, 2, signer.Address(), "test unbond-tx")

	assert.Equal(t, trx1.SignBytes(), trx2.SignBytes())
	assert.NotEqual(t, trx1.SignBytes(), trx3.SignBytes())
	assert.True(t, trx1.IsUnbondTx())

}
func TestWithdrawSignBytes(t *testing.T) {
	stamp := hash.GenerateTestStamp()
	signer := bls.GenerateTestSigner()
	addr := crypto.GenerateTestAddress()

	trx1 := NewWithdrawTx(stamp, 1, signer.Address(), addr, 1000, 1000, "test unbond-tx")
	signer.SignMsg(trx1)

	trx2 := NewWithdrawTx(stamp, 1, signer.Address(), addr, 1000, 1000, "test unbond-tx")
	trx3 := NewWithdrawTx(stamp, 2, signer.Address(), addr, 1000, 1000, "test unbond-tx")

	assert.Equal(t, trx1.SignBytes(), trx2.SignBytes())
	assert.NotEqual(t, trx1.SignBytes(), trx3.SignBytes())
	assert.True(t, trx1.IsWithdrawTx())

}

func TestSortitionSignBytes(t *testing.T) {
	stamp := hash.GenerateTestStamp()
	signer := bls.GenerateTestSigner()
	proof := sortition.GenerateRandomProof()

	trx1 := NewSortitionTx(stamp, 1, signer.Address(), proof)
	signer.SignMsg(trx1)

	trx2 := NewSortitionTx(stamp, 1, signer.Address(), proof)
	trx3 := NewSortitionTx(stamp, 2, signer.Address(), proof)

	assert.Equal(t, trx1.SignBytes(), trx2.SignBytes())
	assert.NotEqual(t, trx1.SignBytes(), trx3.SignBytes())
	assert.True(t, trx1.IsSortitionTx())
}
