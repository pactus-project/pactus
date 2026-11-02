package ntp

import (
	"testing"
	"time"

	"github.com/beevik/ntp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClockOffsetBeforeFirstMeasurement(t *testing.T) {
	td := setup(t)
	defer td.Stop()

	offset, err := td.checker.ClockOffset()
	require.ErrorIs(t, err, ErrNotMeasured)
	assert.Zero(t, offset)
	assert.False(t, td.checker.IsOutOfSync(), "an unmeasured clock is not reported as out of sync")
}

func TestMeasure(t *testing.T) {
	t.Run("offset", func(t *testing.T) {
		td := setup(t)
		defer td.Stop()

		td.querier.res = &ntp.Response{Stratum: 1, ClockOffset: 2 * time.Second}
		td.checker.measure()

		offset, err := td.checker.ClockOffset()
		require.NoError(t, err)
		assert.Equal(t, 2*time.Second, offset)
		assert.True(t, td.checker.IsOutOfSync())
	})

	t.Run("failed measurement is not a zero offset", func(t *testing.T) {
		td := setup(t)
		defer td.Stop()

		td.querier.res = &ntp.Response{Stratum: 0}
		td.checker.measure()

		offset, err := td.checker.ClockOffset()
		require.ErrorIs(t, err, QueryError{})
		assert.Zero(t, offset)
	})
}
