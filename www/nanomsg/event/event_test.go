package event

import (
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/stretchr/testify/assert"
)

func TestCreateBlockEvent(t *testing.T) {
	h, _ := hash.FromString("000102030405060708090a0b0c0d0e0f000102030405060708090a0b0c0d0e0f")
	height := uint32(0x2134)
	e := CreateBlockEvent(h, height)
	assert.Equal(t, Event{
		0x1, 0x1, 0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9,
		0xa, 0xb, 0xc, 0xd, 0xe, 0xf, 0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8,
		0x9, 0xa, 0xb, 0xc, 0xd, 0xe, 0xf, 0x34, 0x21, 0x0, 0x0,
	}, e)
}

func TestCreateNewTransactionEvent(t *testing.T) {
	h, _ := hash.FromString("000102030405060708090a0b0c0d0e0f000102030405060708090a0b0c0d0e0f")
	height := uint32(0x2134)
	e := CreateTransactionEvent(h, height)
	assert.Equal(t, Event{
		0x1, 0x2, 0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9,
		0xa, 0xb, 0xc, 0xd, 0xe, 0xf, 0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8,
		0x9, 0xa, 0xb, 0xc, 0xd, 0xe, 0xf, 0x34, 0x21, 0x0, 0x0,
	}, e)
}

func TestCreateAccountChangeEvent(t *testing.T) {
	addr, _ := crypto.AddressFromString("pc1p0hrct7eflrpw4ccrttxzs4qud2axex4dcdzdfr")
	height := uint32(0x2134)
	e := CreateAccountChangeEvent(addr, height)
	assert.Equal(t, Event{
		0x01, 0x03, 0x1, 0x7d, 0xc7, 0x85, 0xfb, 0x29, 0xf8, 0xc2, 0xea, 0xe3,
		0x3, 0x5a, 0xcc, 0x28, 0x54, 0x1c, 0x6a, 0xba, 0x6c, 0x9a, 0xad, 0x34, 0x21, 0x0, 0x0,
	}, e)
}
