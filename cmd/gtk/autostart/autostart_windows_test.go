package autostart

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/windows/registry"
)

func TestEnableDisable(t *testing.T) {
	// Use a temporary key, so the test does not touch the real Run key.
	runKeyPath = `Software\PactusAutostartTest\` + t.Name()
	t.Cleanup(func() {
		_ = registry.DeleteKey(registry.CURRENT_USER, runKeyPath)
		_ = registry.DeleteKey(registry.CURRENT_USER, `Software\PactusAutostartTest`)
	})

	entry := New([]string{`-working-dir=C:\Users\user\pactus`})
	assert.False(t, entry.IsEnabled())

	require.NoError(t, enable(`C:\Program Files\Pactus\pactus-gui.exe`, entry.args))
	assert.True(t, entry.IsEnabled())

	key, err := registry.OpenKey(registry.CURRENT_USER, runKeyPath, registry.QUERY_VALUE)
	require.NoError(t, err)
	value, _, err := key.GetStringValue(runValueName)
	_ = key.Close()
	require.NoError(t, err)
	assert.Equal(t, `"C:\Program Files\Pactus\pactus-gui.exe" -working-dir=C:\Users\user\pactus`, value)

	require.NoError(t, entry.SetEnabled(false))
	assert.False(t, entry.IsEnabled())

	// Disabling twice is not an error.
	require.NoError(t, entry.SetEnabled(false))
}
