//go:build !linux && !darwin && !windows

package autostart

import "os"

func executable() (string, error) {
	return os.Executable()
}

func enable(string, []string) error {
	return ErrUnsupported
}

func disable() error {
	return ErrUnsupported
}

func isEnabled() bool {
	return false
}
