package version

import (
	"fmt"

	"github.com/coreos/go-semver/semver"
)

// These constants follow the semantic versioning 2.0.0 spec (http://semver.org/)
var (
	major uint   = 1
	minor uint   = 0
	patch uint   = 0
	meta  string = "beta"
)

var build string
var semVer string

func init() {
	sv, err := semver.NewVersion(semVer)
	if err == nil {
		meta = sv.Metadata
		major = uint(sv.Major)
		minor = uint(sv.Minor)
		patch = uint(sv.Patch)
	}

}

func Version() string {
	version := fmt.Sprintf("%d.%d.%d", major, minor, patch)

	if meta != "" {
		version = fmt.Sprintf("%s-%s", version, meta)
	}

	if build != "" {
		version = fmt.Sprintf("%s+%s", version, build)
	}

	return version
}
