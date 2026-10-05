package types_test

import (
	"testing"

	"github.com/pactus-project/pactus/wallet/types"
	"github.com/stretchr/testify/assert"
)

func TestIsLegacyDriver(t *testing.T) {
	assert.True(t, types.IsLegacyDriver(types.DriverLegacyJSON))
	assert.False(t, types.IsLegacyDriver(types.DriverSQLite))
	assert.False(t, types.IsLegacyDriver(""))
}
