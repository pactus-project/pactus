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
