package autostart

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnableDisable(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	entry := New([]string{"-working-dir=/Users/user/pactus"})
	assert.False(t, entry.IsEnabled())

	require.NoError(t, enable("/Applications/pactus-gui.app/Contents/MacOS/pactus-gui-bin", entry.args))
	assert.True(t, entry.IsEnabled())

	data, err := os.ReadFile(filepath.Join(home, "Library", "LaunchAgents", launchAgentLabel+".plist"))
	require.NoError(t, err)
	assert.Contains(t, string(data), "<string>/usr/bin/open</string>")
	assert.Contains(t, string(data), "<string>/Applications/pactus-gui.app</string>")
	assert.Contains(t, string(data), "<string>-working-dir=/Users/user/pactus</string>")
	assert.Contains(t, string(data), "<key>RunAtLoad</key>")

	require.NoError(t, entry.SetEnabled(false))
	assert.False(t, entry.IsEnabled())

	// Disabling twice is not an error.
	require.NoError(t, entry.SetEnabled(false))
}

func TestProgramArguments(t *testing.T) {
	bundleExe := "/Applications/pactus-gui.app/Contents/MacOS/pactus-gui-bin"

	assert.Equal(t,
		[]string{"/usr/bin/open", "-a", "/Applications/pactus-gui.app"},
		programArguments(bundleExe, []string{}))
	assert.Equal(t,
		[]string{"/usr/bin/open", "-a", "/Applications/pactus-gui.app", "--args", "-testnet=true"},
		programArguments(bundleExe, []string{"-testnet=true"}))
	assert.Equal(t,
		[]string{"/usr/local/bin/pactus-gui", "-testnet=true"},
		programArguments("/usr/local/bin/pactus-gui", []string{"-testnet=true"}))
}

func TestLaunchAgentEscapesArguments(t *testing.T) {
	data := string(launchAgent([]string{"/tmp/a&b", "<x>"}))

	assert.Contains(t, data, "<string>/tmp/a&amp;b</string>")
	assert.Contains(t, data, "<string>&lt;x&gt;</string>")
}
