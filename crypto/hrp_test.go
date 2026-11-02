package crypto_test

import (
	"sync"
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/stretchr/testify/assert"
)

func TestHRPSwitchIsRaceFree(t *testing.T) {
	t.Cleanup(crypto.ToMainnetHRP)
	crypto.ToMainnetHRP()

	addr := crypto.NewAddress(crypto.AddressTypeBLSAccount, make([]byte, 20))

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 200; i++ {
			if i%2 == 0 {
				crypto.ToTestnetHRP()
			} else {
				crypto.ToMainnetHRP()
			}
		}
	}()

	for i := 0; i < 200; i++ {
		_ = addr.String()
		_ = crypto.AddressHRP()
	}
	wg.Wait()
}

func TestHRPAccessors(t *testing.T) {
	t.Cleanup(crypto.ToMainnetHRP)

	crypto.ToMainnetHRP()
	assert.Equal(t, "pc", crypto.AddressHRP())

	crypto.ToTestnetHRP()
	assert.Equal(t, "tpc", crypto.AddressHRP())

	crypto.SetAddressHRP("pc")
	assert.Equal(t, "pc", crypto.AddressHRP())
}
