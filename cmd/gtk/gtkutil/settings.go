//go:build gtk

package gtkutil

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
)

// ThemeMode is the user's appearance choice.
type ThemeMode string

const (
	// ThemeSystem follows the operating system light/dark preference.
	ThemeSystem ThemeMode = "system"
	// ThemeLight always uses the light appearance.
	ThemeLight ThemeMode = "light"
	// ThemeDark always uses the dark appearance.
	ThemeDark ThemeMode = "dark"
)

// ApplyThemeMode sets the appearance through the libadwaita style manager.
// With ThemeSystem, the appearance also changes when the operating system
// preference changes while the application is running.
func ApplyThemeMode(mode ThemeMode) {
	scheme := adw.ColorSchemeDefault
	switch mode {
	case ThemeLight:
		scheme = adw.ColorSchemeForceLight
	case ThemeDark:
		scheme = adw.ColorSchemeForceDark
	case ThemeSystem:
	}
	adw.StyleManagerGetDefault().SetColorScheme(scheme)
}

// ApplySavedTheme applies the persisted theme choice. It must run before the
// first window is shown, so the startup assistant, the splash screen and the
// wallet password dialog use the chosen theme too.
func ApplySavedTheme() {
	ApplyThemeMode(LoadThemeMode())
}

// configPath returns a file in the per-user Pactus config directory, which is
// shared across networks. It returns an empty string when the directory is
// unknown.
func configPath(name string) string {
	dir, err := os.UserConfigDir()
	if err != nil {
		return ""
	}

	return filepath.Join(dir, "pactus", name)
}

// themePrefPath returns the file that stores the user's theme choice.
func themePrefPath() string {
	return configPath("gui_theme")
}

// legacyDarkModePath returns the file used by older versions, which stored
// only a light ("0") or dark ("1") choice.
func legacyDarkModePath() string {
	return configPath("gui_dark_mode")
}

// SaveThemeMode persists the user's theme choice.
func SaveThemeMode(mode ThemeMode) {
	path := themePrefPath()
	if path == "" {
		return
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		return
	}

	if err := os.WriteFile(path, []byte(mode), 0o600); err != nil {
		return
	}

	// The new file is now the only source of truth.
	_ = os.Remove(legacyDarkModePath())
}

// LoadThemeMode returns the persisted theme choice. It falls back to the
// choice saved by older versions, and to ThemeSystem when there is none.
func LoadThemeMode() ThemeMode {
	if path := themePrefPath(); path != "" {
		if data, err := os.ReadFile(path); err == nil {
			switch mode := ThemeMode(strings.TrimSpace(string(data))); mode {
			case ThemeSystem, ThemeLight, ThemeDark:
				return mode
			}
		}
	}

	if path := legacyDarkModePath(); path != "" {
		if data, err := os.ReadFile(path); err == nil {
			if strings.TrimSpace(string(data)) == "1" {
				return ThemeDark
			}

			return ThemeLight
		}
	}

	return ThemeSystem
}
