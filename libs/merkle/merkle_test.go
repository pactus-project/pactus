package simpleMerkle

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
)

func strToHash(str string) crypto.Hash {
	h := crypto.HashH([]byte(str))
	return h
}

func TestMerkleTree(t *testing.T) {
	slices := [][]byte{[]byte("a"), []byte("b"), []byte("c")}
	hashes := []crypto.Hash{
		strToHash("a"),
		strToHash("b"),
		strToHash("c"),
	}

	tree1 := NewTreeFromHashes(hashes)
	if tree1.Root().String() != "e6061997a9011668bcf216020aaad9cc7f5f34d5b6f78f1e63ef6257c1aa1f37" {
		t.Errorf("Invalid merkle root")
	}

	tree2 := NewTreeFromSlices(slices)
	assert.Equal(t, tree1.Root(), tree2.Root())
}

func TestMerkleTree_Bitcoin_Block100000(t *testing.T) {
	hasher = func(data []byte) crypto.Hash {
		first := sha256.Sum256(data)
		second := sha256.Sum256(first[:])
		h, _ := crypto.HashFromRawBytes(second[:])
		return h

	}

	// Block 100000 in bitcoin
	wantMerkle, _ := hex.DecodeString("6657A9252AACD5C0B2940996ECFF952228C3067CC38D4885EFB5A4AC4247E9F3")
	txHash1, _ := hex.DecodeString("876DD0A3EF4A2816FFD1C12AB649825A958B0FF3BB3D6F3E1250F13DDBF0148C")
	txHash2, _ := hex.DecodeString("C40297F730DD7B5A99567EB8D27B78758F607507C52292D02D4031895B52F2FF")
	txHash3, _ := hex.DecodeString("C46E239AB7D28E2C019B6D66AD8FAE98A56EF1F21AEECB94D1B1718186F05963")
	txHash4, _ := hex.DecodeString("1D0CB83721529A062D9675B98D6E5C587E4A770FC84ED00ABC5A5DE04568A6E9")

	root, _ := crypto.HashFromRawBytes(wantMerkle)
	h1, _ := crypto.HashFromRawBytes(txHash1)
	h2, _ := crypto.HashFromRawBytes(txHash2)
	h3, _ := crypto.HashFromRawBytes(txHash3)
	h4, _ := crypto.HashFromRawBytes(txHash4)

	hashes := []crypto.Hash{
		h1,
		h2,
		h3,
		h4,
	}

	tree := NewTreeFromHashes(hashes)
	assert.True(t, tree.Root().EqualsTo(root))
}

func TestMerkleTree_Bitcoin_Block113345(t *testing.T) {

	hasher = func(data []byte) crypto.Hash {
		first := sha256.Sum256(data)
		second := sha256.Sum256(first[:])
		h, _ := crypto.HashFromRawBytes(second[:])
		return h

	}

	// Block 113345 in bitcoin
	wantMerkle, _ := hex.DecodeString("31E613DEC2B7D9E78F9FD6E08071B768C5E5FC5E14BE0CDDE728EA19F1EAE3F2")
	txHash1, _ := hex.DecodeString("650E05E757D130F0778300724AF5734D9A57E25072780D2CD0F89D8EC1118FEF")
	txHash2, _ := hex.DecodeString("7D21A7D8984607E94D4E1A298CCDD750331A397BF38623108D496D474206F373")
	txHash3, _ := hex.DecodeString("AF9AF7364A893B67952D87EAB3172AA896287825734AEA96F797F3A0DF1BF1D8")
	txHash4, _ := hex.DecodeString("F85ABB36FCF0BC9C08590273200E0DD63BCD079706973FAB3227565255920F32")
	txHash5, _ := hex.DecodeString("F5BEA40AC2A3DACCF5B2A74B4CCC35EDAEE7203A7F2FF5A512EEEF638E197E32")

	root, _ := crypto.HashFromRawBytes(wantMerkle)
	h1, _ := crypto.HashFromRawBytes(txHash1)
	h2, _ := crypto.HashFromRawBytes(txHash2)
	h3, _ := crypto.HashFromRawBytes(txHash3)
	h4, _ := crypto.HashFromRawBytes(txHash4)
	h5, _ := crypto.HashFromRawBytes(txHash5)

	hashes := []crypto.Hash{
		h1,
		h2,
		h3,
		h4,
		h5,
	}

	tree := NewTreeFromHashes(hashes)
	assert.True(t, tree.Root().EqualsTo(root))
}
