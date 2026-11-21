package autostart

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnableDisable(t *testing.T) {
	configDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", configDir)
	t.Setenv("APPIMAGE", "/opt/Pactus GUI.AppImage")
	t.Setenv("FLATPAK_ID", "")

	entry := New([]string{"-working-dir=/home/user/pactus"})
	assert.False(t, entry.IsEnabled())

	require.NoError(t, entry.SetEnabled(true))
	assert.True(t, entry.IsEnabled())

	data, err := os.ReadFile(filepath.Join(configDir, "autostart", "pactus-gui.desktop"))
	require.NoError(t, err)
	assert.Contains(t, string(data), `Exec="/opt/Pactus GUI.AppImage" "-working-dir=/home/user/pactus"`+"\n")

	require.NoError(t, entry.SetEnabled(false))
	assert.False(t, entry.IsEnabled())

	// Disabling twice is not an error.
	require.NoError(t, entry.SetEnabled(false))
}

func TestEnableInFlatpak(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	// Inside a Flatpak, XDG_CONFIG_HOME is the private folder of the app.
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".var", "app", "org.pactus.PactusGui", "config"))
	t.Setenv("FLATPAK_ID", "org.pactus.PactusGui")
	t.Setenv("APPIMAGE", "")

	entry := New([]string{"-working-dir=/home/user/pactus"})
	require.NoError(t, entry.SetEnabled(true))
	assert.True(t, entry.IsEnabled())

	data, err := os.ReadFile(filepath.Join(home, ".config", "autostart", "pactus-gui.desktop"))
	require.NoError(t, err)
	assert.Contains(t, string(data),
		`Exec="flatpak" "run" "org.pactus.PactusGui" "-working-dir=/home/user/pactus"`+"\n")

	require.NoError(t, entry.SetEnabled(false))
	assert.False(t, entry.IsEnabled())
}

func TestQuoteExecArg(t *testing.T) {
	tests := []struct {
		arg  string
		want string
	}{
		{"/usr/bin/pactus-gui", `"/usr/bin/pactus-gui"`},
		{"/home/a b/pactus", `"/home/a b/pactus"`},
		{`say "hi"`, `"say \\"hi\\""`},
		{"$HOME", `"\\$HOME"`},
		{"100%", `"100%%"`},
		{`a\b`, `"a\\\\b"`},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, quoteExecArg(tt.arg), "arg: %s", tt.arg)
	}
}
