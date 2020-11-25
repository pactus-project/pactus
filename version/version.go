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
		Major: 0,
		Minor: 9,
		Patch: 0,
	}
}

type Version struct {
	Major int
	Minor int
	Patch int
}

func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d-%s",
		v.Major, v.Minor, v.Patch, GitCommit)
}
