package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConf(t *testing.T) {
	conf := DefaultConfig()

	conf.EnableRelay = true
	assert.Error(t, conf.SanityCheck())

	conf.EnableRelay = false
	assert.NoError(t, conf.SanityCheck())
}

func TestValidateAddress(t *testing.T) {
	validAddresses := []string{
		"/ip4/127.0.0.1/tcp/21777/p2p/QmUv7r8T5NT3xQ2Zzvbg53uSdVccQ6LbdY2wAju14HsWyc",
		"/ip4/172.104.46.145/tcp/21777/p2p/12D3KooWNYD4bB82YZRXv6oNyYPwc5ozabx2epv75ATV3D8VD3Mq",
	}
	assert.Equal(t, true, validateAddress(validAddresses))

	inValidAddresses := []string{
		"/ip4/1wg.1/tcp/21777/p2p8VccQ8888888Wyc",
		"/ip4/172.104.46.145/tcp21777/p2p/12D3KooWNYD4bB82YZRXv6oNyYPwc5ozabx2epv75ATV3D8VD3Mq",
		"/ip3/172.104.46.145/tcp21777/p2p/12D3KooWNYD4bB82YZRXv6oNyYPwc5ozabx2epv75ATV3D8VD3Mq",
		"/ip3/172.104.46.145$tcp21777/p2p/12D3KooWNYD4bB82YZRXv6oNyYPwc5ozabx2ep$TV3D8VD3Mq",
	}
	assert.Equal(t, false, validateAddress(inValidAddresses))
}
