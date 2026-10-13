package status

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatus(t *testing.T) {
	tests := []struct {
		status           Status
		name             string
		isUnknown        bool
		isKnown          bool
		isDisconnected   bool
		isConnected      bool
		isConnectedKnown bool
		isBanned         bool
	}{
		{StatusBanned, "banned", false, false, false, false, false, true},
		{StatusUnknown, "unknown", true, false, false, false, false, false},
		{StatusDisconnected, "disconnected", false, false, true, false, false, false},
		{StatusConnected, "connected", false, false, false, true, true, false},
		{StatusKnown, "known", false, true, false, false, true, false},
		{Status(100), "invalid", false, false, false, false, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.name, tt.status.String())
			assert.Equal(t, tt.isUnknown, tt.status.IsUnknown())
			assert.Equal(t, tt.isKnown, tt.status.IsKnown())
			assert.Equal(t, tt.isDisconnected, tt.status.IsDisconnected())
			assert.Equal(t, tt.isConnected, tt.status.IsConnected())
			assert.Equal(t, tt.isConnectedKnown, tt.status.IsConnectedOrKnown())
			assert.Equal(t, tt.isBanned, tt.status.IsBanned())
		})
	}
}
