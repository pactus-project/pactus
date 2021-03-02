package sortition

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
)

func TestVRF(t *testing.T) {
	_, pk, pv := crypto.GenerateTestKeyPair()
	signer := crypto.NewSigner(pv)
	for i := 0; i < 100; i++ {
		seed := GenerateRandomSeed()
		fmt.Printf("index is : %x \n", seed)

		vrf := NewVRF()

		//max := int64(i * 1000)
		max := int64(1 * 1e6)
		vrf.SetMax(max)
		index, proof := vrf.Evaluate(seed, signer)
		// fmt.Printf("index is : %v \n", index)

		assert.LessOrEqual(t, index, max)

		index2, result := vrf.Verify(seed, pk, proof)

		assert.Equal(t, result, true)
		assert.Equal(t, index, index2)
	}
}

func TestEntropy(t *testing.T) {
	_, _, pv := crypto.GenerateTestKeyPair()
	signer := crypto.NewSigner(pv)

	max := int64(100)
	vrf := NewVRF()
	vrf.SetMax(max)

	entropy := make([]bool, max)
	for i := int64(0); i < max; i++ {
		seed := GenerateRandomSeed()

		index, _ := vrf.Evaluate(seed, signer)
		assert.LessOrEqual(t, index, max)

		entropy[index] = true
	}

	hits := int64(0)
	for _, b := range entropy {
		if b == true {
			hits++
		}
	}

	fmt.Printf("Entropy is : %v%% \n", hits*100/max)
	assert.Greater(t, hits, int64(50))
}

func TestGetIndex(t *testing.T) {
	//  TotalStake: 1000000
	vrf := VRF{max: 1 * 1e6}

	// proof: 1719b896ec1cc66a0f44c4bf90890d988e341cb2c1a808907780af844c854291536c12fdaef9a526bb7ef80da17c0b03
	// proofH: a7b8166584387f4ea76f9caa0969bd6b0bb8df4c3bb8e87f8b6e4dad62bf3359
	//
	// 0x4e7f38846516b8a7 & 0x7fffffffffffffff = 0x4e7f38846516b8a7
	// 0x4e7f38846516b8a7*1000000/0x7fffffffffffffff=613257.46
	proof1, _ := ProofFromString("1719b896ec1cc66a0f44c4bf90890d988e341cb2c1a808907780af844c854291536c12fdaef9a526bb7ef80da17c0b03")
	assert.Equal(t, vrf.getIndex(proof1), int64(613257))

	// proof: 45180defab2daae377977bf09dcdd7d76ff4fc96d1b50cc8ac5a1601c0522fb11641c3ed0fefd4b1e1808c498d699396
	// proofH: 80212979d1de1ca4ce1258fc0be66a4453b3804e64a5ca8d95f7def2c291c7fe
	//
	// 0xa41cded179292180 & 0x7fffffffffffffff = 0x241cded179292180
	// 0x241cded179292180*1000000/0x7fffffffffffffff=282131.05419337929094808699
	proof2, _ := ProofFromString("45180defab2daae377977bf09dcdd7d76ff4fc96d1b50cc8ac5a1601c0522fb11641c3ed0fefd4b1e1808c498d699396")
	assert.Equal(t, vrf.getIndex(proof2), int64(282131))
}
