//go:build gtk

package gtkutil

import (
	"os"
	"path/filepath"
	"strings"
)

// themePrefPath returns the file that stores the user's light/dark choice,
// kept in the per-user config directory so it is shared across networks.
func themePrefPath() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		return ""
	}

	return filepath.Join(dir, "pactus", "gui_dark_mode")
}

// SaveDarkMode persists the user's dark-mode choice.
func SaveDarkMode(dark bool) {
	path := themePrefPath()
	if path == "" {
		return
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return
	}

	value := "0"
	if dark {
		value = "1"
	}

	_ = os.WriteFile(path, []byte(value), 0o644)
}

// LoadDarkMode returns the persisted dark-mode choice. ok is false when the
// user has not made a choice yet, in which case the system theme is followed.
func LoadDarkMode() (dark, ok bool) {
	path := themePrefPath()
	if path == "" {
		return false, false
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return false, false
	}

	return strings.TrimSpace(string(data)) == "1", true
}
