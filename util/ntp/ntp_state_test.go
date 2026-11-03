package ntp

import (
	"errors"
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

type perServerQuerier struct {
	responses map[string]*ntp.Response
	queried   []string
}

func (q *perServerQuerier) Query(server string) (*ntp.Response, error) {
	q.queried = append(q.queried, server)

	res, ok := q.responses[server]
	if !ok {
		return nil, errors.New("unable to query")
	}

	return res, nil
}

func TestClockOffsetIsTheMedianOfSeveralPools(t *testing.T) {
	response := func(offset time.Duration) *ntp.Response {
		return &ntp.Response{Stratum: 1, ClockOffset: offset}
	}

	t.Run("an outlier is ignored", func(t *testing.T) {
		querier := &perServerQuerier{responses: map[string]*ntp.Response{
			_pools[0]: response(time.Hour),
			_pools[1]: response(100 * time.Millisecond),
			_pools[2]: response(200 * time.Millisecond),
			_pools[3]: response(300 * time.Millisecond),
		}}
		checker := NewNtpChecker(t.Context(), WithQuerier(querier))
		defer checker.Stop()

		checker.measure()

		offset, err := checker.ClockOffset()
		require.NoError(t, err)
		assert.Equal(t, 200*time.Millisecond, offset)
		assert.Len(t, querier.queried, maxSamples)
		assert.False(t, checker.IsOutOfSync())
	})

	t.Run("failed pools are skipped", func(t *testing.T) {
		querier := &perServerQuerier{responses: map[string]*ntp.Response{
			_pools[2]: response(100 * time.Millisecond),
			_pools[4]: response(300 * time.Millisecond),
		}}
		checker := NewNtpChecker(t.Context(), WithQuerier(querier))
		defer checker.Stop()

		checker.measure()

		offset, err := checker.ClockOffset()
		require.NoError(t, err)
		assert.Equal(t, 200*time.Millisecond, offset)
		assert.Len(t, querier.queried, len(_pools))
	})

	t.Run("a single answer is enough", func(t *testing.T) {
		querier := &perServerQuerier{responses: map[string]*ntp.Response{
			_pools[5]: response(150 * time.Millisecond),
		}}
		checker := NewNtpChecker(t.Context(), WithQuerier(querier))
		defer checker.Stop()

		checker.measure()

		offset, err := checker.ClockOffset()
		require.NoError(t, err)
		assert.Equal(t, 150*time.Millisecond, offset)
	})
}
