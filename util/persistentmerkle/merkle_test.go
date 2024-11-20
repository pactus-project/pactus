package persistentmerkle

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNodeID(t *testing.T) {
	assert.Equal(t, uint32(0x00000000), nodeID(0, 0))
	assert.Equal(t, uint32(0x01000000), nodeID(0, 1))
	assert.Equal(t, uint32(0x00000001), nodeID(1, 0))
	assert.Equal(t, uint32(0x01000001), nodeID(1, 1))
	assert.Equal(t, uint32(0xffffffff), nodeID(0xffffff, 0xff))
	assert.Equal(t, uint32(0x00ffffff), nodeID(0xffffff, 0x00))
	assert.Equal(t, uint32(0x77ff00ff), nodeID(0xff00ff, 0x77))
}

func TestCalculateHeight(t *testing.T) {
	tree := New()

	tree.recalculateHeight(0)
	assert.Equal(t, int32(0), tree.maxHeight)

	tree.recalculateHeight(1)
	assert.Equal(t, int32(1), tree.maxHeight)

	tree.recalculateHeight(2)
	assert.Equal(t, int32(2), tree.maxHeight)

	tree.recalculateHeight(4)
	assert.Equal(t, int32(3), tree.maxHeight)

	tree.recalculateHeight(5)
	assert.Equal(t, int32(4), tree.maxHeight)

	tree.recalculateHeight(8)
	assert.Equal(t, int32(4), tree.maxHeight)

	tree.recalculateHeight(9)
	assert.Equal(t, int32(5), tree.maxHeight)
}

func TestMerkleTree(t *testing.T) {
	tree := New()

	data := []string{
		"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
		"N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
	}
	roots := []string{
		"a3a0081351bb785d0758ddf68076a95ffd3f10b88bbc9911e9fea4d793c06414",
		"f8c780ed425c674e458b1e42abc894d144af4660476b63adc900e5ca72ef3a7d",
		"fba254b63bc4f71e560d3d94bacaa555b57b0d15073bdde95d028176cc702952",
		"62edcdc731948255463d6a8b7809bf26de2bdc1e2b351e76bf0adb5a2b863f5a",
		"c0c5e5d5389593188ef48b73f808571f88252491696cb481d8d587e3d59a282b",
		"0ff43eeca0576185050fea37c7cfc81219ceca762f9e29bc77c57a651efc6ae6",
		"e071cda6628d87d405c82c9948795b99b5afabb01c5a93b4970b4f97b866b32b",
		"f4b351a896639e122c81c5819d278fc67392abc87592daeaa1bddda9bec57186",
		"0086ad0c0b55b087d82733f8feaff319bc9dd7bcfbc6ddd3f8322e4dc9bd492e",
		"b0b18c74adc184b14f6dc6546062aad6475f4b167334bc4a4a537c5b78a22b44",
		"fbb6818ed0eb4b31ae3a273498e2cd26d4dbdb44d09ea71a115958f31687ebf5",
		"8e7cc64ef18bf44612520a06e9cfa2b6117d2587a1832ec79365b463656b81ea",
		"7d1fa0a14c9d83fde2638c7f404041a667365c1ada96e59c42ff2573a4b1bd83",
		"dcab86df00b139cec708c3cfe2c29cc4adf47ba3b9981059e3f62403d26635bb",
		"6d44b6badb9521bdafc83b84bfd4ce3dbc2eae270ad22b676c84e356fab81c9c",
		"41cfad12a1964a1e5f534dd3b247f6b9e80fb6e64772ad9ff7f7f8cee375b94c",
		"673fe9e5a964a3d2fd2fea41be5a4bc8e6fcc958422f6636d26c25156902f0ea",
		"880e8f2691c0ebc2292b8e03733ad96796a4f9792d1a604a28e575fbfcc56b45",
		"a1796d57ebf723c58b8518c6c70ee602f772afe6bc3b127a74804c0944957162",
		"3618d980b28f432b383d8658c0d8ee8c087c206df59a1d535c7512c85fa745c3",
		"e05db943b9ef0f90c90989d2bd29e74ece5cc99882d7ef75832b575cd14f8c88",
		"4c490979a140a3d4a4661ccb0804c08aabb3ee15fc941d473f866fa04cef0613",
		"414e11e90180869501a2bf1762be6e3e5c34eaa2b9eb7510b7a77b7bf32f08d2",
		"e0e435e50a8bb83d1610b4288f11a4d313db0b47a61929078990508185243bbc",
		"7f8da04ec98800942cb7b9bb2fe795b51cdb91b2ff48b3451d4e74da0414b6ef",
		"1d1f6d1c6e1afb8a6c3e6b93d708cebbfcaab5fc8d27c7d44ff5454d8906e3e9",
	}

	for i, d := range data {
		tree.SetData(int32(i), []byte(d))
		expected, _ := hex.DecodeString(roots[i])
		assert.Equal(t, expected, tree.Root().Bytes(), "Root %d not matched", i)
	}

	// Modifying some data blocks
	tree.SetData(0, []byte("a"))
	tree.SetData(21, []byte("v"))
	expected, _ := hex.DecodeString("ec4446ea16b8f82083cc2d727b8f9e7b9c318e35bb37295a2e87064393572800")
	assert.Equal(t, expected, tree.Root().Bytes())
}
