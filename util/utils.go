package util

import (
	crand "crypto/rand"
	"fmt"
	"math/big"

	"github.com/libp2p/go-libp2p-core/peer"
)

const MaxUint64 = ^uint64(0)
const MinUint64 = 0
const MaxInt64 = int64(MaxUint64 >> 1)
const MinInt64 = -MaxInt64 - 1

/// Max returns the biggest of two numbers
func Max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

/// Min returns the smallest of two numbers
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

/// Min64 returns the smallest of two numbers
func Min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

/// Max64 returns the biggest of two numbers
func Max64(a, b int64) int64 {
	if a < b {
		return b
	}
	return a
}

/// RandInt returns a random value in [0, max)
/// If max set to zero or negative, the max will set to Big.Int
func RandInt(max int) int {
	rnd := RandInt64(int64(max))
	return int(rnd)
}

/// RandInt64 returns a uniform random value in [0, max)
/// If max set to zero or negative, the max will set to Big.Int64
func RandInt64(max int64) int64 {
	if max <= 0 {
		max = MaxInt64
	}
	bigMax := big.NewInt(max)
	bigRnd, _ := crand.Int(crand.Reader, bigMax)
	return bigRnd.Int64()
}

/// RandomPeerID returns a random peer ID
func RandomPeerID() peer.ID {
	s := RandomSlice(32)
	id := [34]byte{0x12, 32}
	copy(id[2:], s[:])
	return peer.ID(id[:])

}

/// FingerprintPeerID returns a pretty and short string for the given peer ID
func FingerprintPeerID(id peer.ID) string {
	pid := id.Pretty()
	return fmt.Sprintf("%s*%s", pid[:2], pid[len(pid)-6:])
}

/// SetFlag applies mask to the flags
func SetFlag(flags, mask int) int {
	flags = flags | mask
	return flags
}

/// UnsetFlag removes mask from the flags
func UnsetFlag(flags, mask int) int {
	flags = flags & ^mask
	return flags
}

/// IsFlagSet checks if flags is set or not
func IsFlagSet(flags, mask int) bool {
	return flags&mask == mask
}
