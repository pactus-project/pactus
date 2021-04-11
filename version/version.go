package version

import (
	"fmt"
)

var (
	NodeVersion Version
	GitCommit   string
)

func init() {
	NodeVersion = Version{
		Major: 1,
		Minor: 0,
		Patch: 0,
	}
}

type Version struct {
	Major int `cbor:"1,keyasint"`
	Minor int `cbor:"2,keyasint"`
	Patch int `cbor:"3,keyasint"`
}

func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d-%s",
		v.Major, v.Minor, v.Patch, GitCommit)
}

func (v Version) MarshalText() ([]byte, error) {
	str := fmt.Sprintf("%d.%d.%d",
		v.Major, v.Minor, v.Patch)
	return []byte(str), nil
}
