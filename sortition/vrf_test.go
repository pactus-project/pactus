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
		h := crypto.GenerateTestHash()
		vrf := NewVRF(signer)

		max := int64(i + 1*1000)
		vrf.SetMax(max)
		index, proof := vrf.Evaluate(h)
		//fmt.Printf("index is : %v \n", index)

		assert.Equal(t, index <= max, true)

		index2, result := vrf.Verify(h, pk, proof)

		assert.Equal(t, result, true)
		assert.Equal(t, index, index2)
	}
}

func TestEntropy(t *testing.T) {
	_, _, pv := crypto.GenerateTestKeyPair()
	signer := crypto.NewSigner(pv)

	entropy := make([]bool, 100)
	for i := 0; i < 100; i++ {
		h := crypto.GenerateTestHash()

		vrf := NewVRF(signer)

		max := int64(100)
		vrf.SetMax(max)
		index, _ := vrf.Evaluate(h)

		entropy[index] = true
	}

	hits := 0
	for _, b := range entropy {
		if b == true {
			hits++
		}
	}

	fmt.Printf("Entropy is : %v \n", hits)
	assert.Greater(t, hits, 50)
}
