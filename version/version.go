package version

import (
	"fmt"
)

// These constants follow the semantic versioning 2.0.0 spec (http://semver.org/)
const (
	major uint   = 0
	minor uint   = 15
	patch uint   = 0
	meta  string = "beta"
)

var build string

func Agent() string {
	return fmt.Sprintf("pactus/%s", Version())
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
