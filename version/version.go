package version

import (
	"fmt"
	"runtime"
)

// These constants follow the semantic versioning 2.0.0 spec (http://semver.org/)
const (
	major uint   = 0
	minor uint   = 18
	patch uint   = 0
	meta  string = "beta"
)

func Agent() string {
	return fmt.Sprintf("node=pactus/version=%s/os=%s/architecture=%s", Version(), runtime.GOOS, runtime.GOARCH)
}

func Version() string {
	version := fmt.Sprintf("%d.%d.%d", major, minor, patch)

	if meta != "" {
		version = fmt.Sprintf("%s-%s", version, meta)
	}

	return version
}
