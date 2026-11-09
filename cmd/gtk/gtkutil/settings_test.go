//go:build gtk

package gtkutil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// useTempConfigDir points os.UserConfigDir to a temporary directory on every
// platform and returns the Pactus config directory inside it.
func useTempConfigDir(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	t.Setenv("AppData", dir)         // Windows
	t.Setenv("XDG_CONFIG_HOME", dir) // Linux and other Unix
	t.Setenv("HOME", dir)            // macOS

	configDir, err := os.UserConfigDir()
	require.NoError(t, err)

	return filepath.Join(configDir, "pactus")
}

func writeFile(t *testing.T, path, data string) {
	t.Helper()

	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o750))
	require.NoError(t, os.WriteFile(path, []byte(data), 0o600))
}

func TestLoadThemeMode_DefaultIsSystem(t *testing.T) {
	useTempConfigDir(t)

	assert.Equal(t, ThemeSystem, LoadThemeMode())
}

func TestThemeMode_SaveAndLoad(t *testing.T) {
	useTempConfigDir(t)

	for _, mode := range []ThemeMode{ThemeLight, ThemeDark, ThemeSystem} {
		SaveThemeMode(mode)
		assert.Equal(t, mode, LoadThemeMode())
	}
}

func TestLoadThemeMode_LegacyChoice(t *testing.T) {
	tests := []struct {
		data string
		want ThemeMode
	}{
		{"1", ThemeDark},
		{"1\n", ThemeDark},
		{"0", ThemeLight},
	}

	for _, tt := range tests {
		dir := useTempConfigDir(t)
		writeFile(t, filepath.Join(dir, "gui_dark_mode"), tt.data)

		assert.Equal(t, tt.want, LoadThemeMode(), "legacy value %q", tt.data)
	}
}

func TestSaveThemeMode_RemovesLegacyFile(t *testing.T) {
	dir := useTempConfigDir(t)
	legacy := filepath.Join(dir, "gui_dark_mode")
	writeFile(t, legacy, "1")

	SaveThemeMode(ThemeLight)

	assert.NoFileExists(t, legacy)
	assert.Equal(t, ThemeLight, LoadThemeMode())
}

func TestLoadThemeMode_InvalidValueFallsBack(t *testing.T) {
	dir := useTempConfigDir(t)
	writeFile(t, filepath.Join(dir, "gui_theme"), "purple")

	assert.Equal(t, ThemeSystem, LoadThemeMode())

	writeFile(t, filepath.Join(dir, "gui_dark_mode"), "1")
	assert.Equal(t, ThemeDark, LoadThemeMode())
}
