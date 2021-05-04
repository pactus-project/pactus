package util

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/libp2p/go-libp2p-core/peer"
)

const MaxUint64 = ^uint64(0)
const MinUint64 = 0
const MaxInt64 = int64(MaxUint64 >> 1)
const MinInt64 = -MaxInt64 - 1

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func Min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func Max64(a, b int64) int64 {
	if a < b {
		return b
	}
	return a
}

func RandInt(max int) int {
	if max <= 0 {
		return 0
	}
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max)
}

func RandInt64(max int64) int64 {
	if max <= 0 {
		return 0
	}
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Int63n(max)
}

func RandomPeerID() peer.ID {
	s := Int64ToSlice(RandInt64(MaxInt64))
	id := [34]byte{0x12, 32}
	copy(id[2:], s[:])
	return peer.ID(id[:])

}
func FingerprintPeerID(id peer.ID) string {
	pid := id.Pretty()
	return fmt.Sprintf("%s*%s", pid[:2], pid[len(pid)-6:])
}

func SetFlag(flags, mask int) int {
	flags = flags | mask
	return flags
}

func IsFlagSet(flags, mask int) bool {
	return flags&mask == mask
}
