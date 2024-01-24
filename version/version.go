package version

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// These constants follow the semantic versioning 2.0.0 spec (http://semver.org/)
const (
	major           uint   = 1
	minor           uint   = 0
	patch           uint   = 0
	meta            string = "rc-0"
	protocolVersion uint   = 1
)

func Agent() string {
	return fmt.Sprintf("node=%s/node-version=v%s/protocol-version=%d/os=%s/arch=%s",
		ExecutorName(), Version(), protocolVersion, runtime.GOOS, runtime.GOARCH)
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
