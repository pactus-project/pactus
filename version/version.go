package version

import (
	"fmt"
	"os"
	"path/filepath"
)

// These constants follow the semantic versioning 2.0.0 spec (http://semver.org/)
const (
	major           uint   = 0
	minor           uint   = 20
	patch           uint   = 0
	meta            string = "beta"
	protocolVersion uint   = 1
)

func Agent() string {
	return fmt.Sprintf("pactus/%s", Version())
}

func Version() string {
	version := fmt.Sprintf("%d.%d.%d", major, minor, patch)
	if meta != "" {
		version = fmt.Sprintf("%s-%s", version, meta)
	}

	return version
}

func ExecutorName() string {
	executorName, _ := os.Executable()
	return filepath.Base(executorName)
}
