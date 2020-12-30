package util

import (
	"fmt"
	"math/rand"
	"time"

	peer "github.com/libp2p/go-libp2p-peer"
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
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max)
}

func RandInt64(max int64) int64 {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Int63n(max)
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
