package address_path

import (
	"fmt"
	"strconv"
	"strings"
)

const HardenedKeyStart = uint32(0x80000000) // 2^31

type Path []uint32

func NewPath(indexes ...uint32) Path {
	p := make([]uint32, 0, len(indexes))
	p = append(p, indexes...)

	return p
}

func NewPathFromString(str string) (Path, error) {
	sub := strings.Split(str, "/")
	if sub[0] != "m" {
		return nil, ErrInvalidPath
	}
	path := []uint32{}
	for i := 1; i < len(sub); i++ {
		indexStr := sub[i]
		added := uint32(0)
		if indexStr[len(indexStr)-1] == '\'' {
			added = HardenedKeyStart
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
		if i >= HardenedKeyStart {
			builder.WriteString(fmt.Sprintf("/%d'", i-HardenedKeyStart))
		} else {
			builder.WriteString(fmt.Sprintf("/%d", i))
		}
	}
	return builder.String()
}
