package version

import (
	"fmt"

	"github.com/coreos/go-semver/semver"
)

var (
	NodeVersion Version = Version{
		Major: 1,
		Minor: 0,
		Patch: 0,
	}
	SemVersion string = "1.0.0-beta"
	GitCommit  string
)

func init() {
	sv, err := semver.NewVersion(SemVersion)
	if err == nil {
		NodeVersion = Version{
			Major: int(sv.Major),
			Minor: int(sv.Minor),
			Patch: int(sv.Patch),
		}
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
