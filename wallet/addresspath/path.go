package addresspath

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pactus-project/pactus/crypto"
	"golang.org/x/exp/constraints"
)

type Purpose uint32

const (
	PurposeBLS12381         Purpose = 12381
	PurposeBIP44            Purpose = 44
	PurposeImportPrivateKey Purpose = 65535
)

type CoinType uint32

const (
	CoinTypePactusMainnet CoinType = 21888
	CoinTypePactusTestnet CoinType = 21777
)

const hardenedKeyStart = uint32(0x80000000) // 2^31

// Harden hardens the integer value 'i' by adding 0x80000000 (2^31) to it.
// This function does not check if 'i' is already hardened.
func Harden[T constraints.Integer](i T) uint32 {
	return uint32(i) + hardenedKeyStart
}

// UnHarden unhardens the integer value 'i' by subtracting 0x80000000 (2^31) from it.
// This function does not check if 'i' is already non-hardened.
func UnHarden[T constraints.Integer](i T) uint32 {
	return uint32(i) - hardenedKeyStart
}

type Path []uint32

func NewPath(indexes ...uint32) Path {
	p := make([]uint32, 0, len(indexes))
	p = append(p, indexes...)

	return p
}

func FromString(str string) (Path, error) {
	sub := strings.Split(str, "/")
	if sub[0] != "m" {
		return nil, ErrInvalidPath
	}

	// TODO: check the path should exactly 4 levels.

	var path []uint32
	for i := 1; i < len(sub); i++ {
		indexStr := sub[i]
		added := uint32(0)
		if indexStr[len(indexStr)-1] == '\'' {
			added = hardenedKeyStart
			indexStr = indexStr[:len(indexStr)-1]
		}
		val, err := strconv.ParseInt(indexStr, 10, 32)
		if err != nil {
			return nil, err
		}
		path = append(path, uint32(val)+added)
	}

	return path, nil
}

func (p Path) String() string {
	var builder strings.Builder
	builder.WriteString("m")
	for _, i := range p {
		if i >= hardenedKeyStart {
			fmt.Fprintf(&builder, "/%d'", i-hardenedKeyStart)
		} else {
			fmt.Fprintf(&builder, "/%d", i)
		}
	}

	return builder.String()
}

// TODO: we can add IsBLSPurpose or IsImportedPurpose functions

func (p Path) Purpose() Purpose {
	return Purpose(UnHarden(p[0]))
}

func (p Path) CoinType() CoinType {
	return CoinType(UnHarden(p[len(p)-3]))
}

func (p Path) AddressType() crypto.AddressType {
	return crypto.AddressType(UnHarden(p[len(p)-2]))
}

func (p Path) AddressIndex() uint32 {
	return p[len(p)-1]
}
