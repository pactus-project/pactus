package simplemerkle

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/stretchr/testify/assert"
)

func strToHash(str string) hash.Hash {
	h := hash.CalcHash([]byte(str))

	return h
}

func TestMerkleTree(t *testing.T) {
	slices := [][]byte{[]byte("a"), []byte("b"), []byte("c")}
	hashes := []hash.Hash{
		strToHash("a"),
		strToHash("b"),
		strToHash("c"),
	}

	tree1 := NewTreeFromHashes(hashes)
	if tree1.Root().String() != "e6061997a9011668bcf216020aaad9cc7f5f34d5b6f78f1e63ef6257c1aa1f37" {
		t.Error("invalid merkle root")
	}

	tree2 := NewTreeFromSlices(slices)
	assert.Equal(t, tree1.Root(), tree2.Root())

	fmt.Println(tree2.ToString())
}

func TestMerkleTreeDepth2(t *testing.T) {
	slices := [][]byte{[]byte("a"), []byte("b"), []byte("c")}
	hashes := []hash.Hash{
		strToHash("a"),
		strToHash("b"),
		strToHash("c"),
	}

	tree1 := NewTreeFromHashes(hashes)
	if tree1.Root().String() != "e6061997a9011668bcf216020aaad9cc7f5f34d5b6f78f1e63ef6257c1aa1f37" {
		t.Error("invalid merkle root")
	}

	tree2 := NewTreeFromSlices(slices)
	assert.Equal(t, tree1.Root(), tree2.Root())

	fmt.Println(tree2.ToString())
}

func TestMerkleTree_Bitcoin_Block100000(t *testing.T) {
	hasher = func(data []byte) hash.Hash {
		first := sha256.Sum256(data)
		second := sha256.Sum256(first[:])
		h, _ := hash.FromBytes(second[:])

		return h
	}

	// Block 100000 in bitcoin
	root, _ := hash.FromString("6657A9252AACD5C0B2940996ECFF952228C3067CC38D4885EFB5A4AC4247E9F3")
	hash1, _ := hash.FromString("876DD0A3EF4A2816FFD1C12AB649825A958B0FF3BB3D6F3E1250F13DDBF0148C")
	hash2, _ := hash.FromString("C40297F730DD7B5A99567EB8D27B78758F607507C52292D02D4031895B52F2FF")
	hash3, _ := hash.FromString("C46E239AB7D28E2C019B6D66AD8FAE98A56EF1F21AEECB94D1B1718186F05963")
	hash4, _ := hash.FromString("1D0CB83721529A062D9675B98D6E5C587E4A770FC84ED00ABC5A5DE04568A6E9")

	hashes := []hash.Hash{
		hash1,
		hash2,
		hash3,
		hash4,
	}

	tree := NewTreeFromHashes(hashes)
	assert.Equal(t, root, tree.Root())
	assert.Equal(t, 2, tree.Depth())
}

func TestMerkleTree_Bitcoin_Block153342(t *testing.T) {
	hasher = func(data []byte) hash.Hash {
		first := sha256.Sum256(data)
		second := sha256.Sum256(first[:])
		h, _ := hash.FromBytes(second[:])

		return h
	}

	// Block 153342 in bitcoin
	wantMerkle, _ := hex.DecodeString("dd8ee246e19ec5c77ddd46c1138e8af6a272da4dbb6500ea74a79c0bf89e2c07")
	hash1, _ := hash.FromString("216404816ca6261f9206d471d0403ba49bda4264719d879819fbda9849781e62")
	hash2, _ := hash.FromString("56f2602c15cb0b8e0b38e54b2961a2e541a7febbe852516cd425aa5fb72c5578")
	hash3, _ := hash.FromString("0d065da59871386321c2c9b2e4b6482426bcce88600ab7f55f0d27b9916a9e0c")
	hash4, _ := hash.FromString("1129038c38783f4c4241e54d9d702965b305b8d1e54c091fdd9f9df21240586e")
	hash5, _ := hash.FromString("81461f9e0e093dad14d0c5fb3978431a321bf61de33512d6cc344edb86f359f3")
	hash6, _ := hash.FromString("22140f4b15d76ff27d657a731fdc3040487c22ee3577c6522239d9cfbe0292ad")
	hash7, _ := hash.FromString("0fa273bce5137a0dbffac068ebb6f1ebe64e6be2b00cdae5a967edeb0cd96b93")
	hash8, _ := hash.FromString("cab481631e7f2f7d864a65d23c34bd357f46ecba60bb8117f55ed43232aa75e5")
	hash9, _ := hash.FromString("dffea4c267fa6949111fed23b15977d5e2efa82fefd9cd5ac81e38518d2c2bef")
	hash10, _ := hash.FromString("ed9f4ee5e07a47a7026725173de32efa7372243117be1aa7f60a650aef075475")
	hash11, _ := hash.FromString("8822c80afa3eb84bc3603509b8b6deeee37cf771ca7b49d3dd73294e05f7b29f")
	hash12, _ := hash.FromString("23ad44934167cc712b358f2a097b7316ca2b3c2f34472017273969e7c7e5cdb4")
	hash13, _ := hash.FromString("c1dc3762c6a57757a9aa895b8229613d96f272f79d14c9854132b980eaa2a2c4")

	root, _ := hash.FromBytes(wantMerkle)

	hashes := []hash.Hash{
		hash1, hash2, hash3, hash4, hash5, hash6, hash7, hash8, hash9, hash10, hash11, hash12, hash13,
	}

	tree := NewTreeFromHashes(hashes)
	assert.Equal(t, root, tree.Root())
	assert.Equal(t, 4, tree.Depth())
	fmt.Println(tree.ToString())
	assert.Contains(t, tree.ToString(), root.String())

	right, _ := hash.FromString("4a3ee07bb7baf6dfa265fa5c85a8955c8e79ddab0f70657a14df5744a103e24d")
	left, _ := hash.FromString("114799e25e6dc376d65fd5406516919e1e619b89316be91ea064a69400472d1e")

	root2 := HashMerkleBranches(&left, &right)
	assert.Equal(t, root, *root2)
}
