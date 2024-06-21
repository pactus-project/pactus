package ntp

import (
	"fmt"
	"testing"
	"time"

	"github.com/beevik/ntp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testQuerier struct {
	err error
	res *ntp.Response
}

func (q *testQuerier) Query(_ string) (*ntp.Response, error) {
	return q.res, q.err
}

type testData struct {
	checker *Checker
	querier *testQuerier
}

func setup(t *testing.T) *testData {
	t.Helper()

	querier := &testQuerier{}
	opts := []CheckerOption{
		WithInterval(100 * time.Millisecond),
		WithThreshold(1 * time.Second),
		WithQuerier(querier),
	}

	checker := NewNtpChecker(opts...)

	td := &testData{
		checker: checker,
		querier: querier,
	}

	return td
}

func (td *testData) Stop() {
	td.checker.Stop()
}

func TestNTPChecker(t *testing.T) {
	t.Run("Offset less than one second", func(t *testing.T) {
		td := setup(t)
		defer td.Stop()

		ntpOffset := 100 * time.Millisecond
		td.querier.res = &ntp.Response{
			Stratum:     1,
			ClockOffset: ntpOffset,
		}

		go td.checker.Start()

		require.Eventually(t, func() bool {
			offset, _ := td.checker.ClockOffset()

			return offset == ntpOffset
		}, 5*time.Second, 100*time.Millisecond)

		offset, err := td.checker.ClockOffset()
		assert.NoError(t, err)
		assert.Equal(t, ntpOffset, offset)
		assert.False(t, td.checker.IsOutOfSync())
	})

	t.Run("Offset more than one second", func(t *testing.T) {
		td := setup(t)
		defer td.Stop()

		ntpOffset := 2 * time.Second
		td.querier.res = &ntp.Response{
			Stratum:     1,
			ClockOffset: ntpOffset,
		}

		go td.checker.Start()

		require.Eventually(t, func() bool {
			offset, _ := td.checker.ClockOffset()

			return offset == ntpOffset
		}, 5*time.Second, 100*time.Millisecond)

		offset, err := td.checker.ClockOffset()
		assert.NoError(t, err)
		assert.Equal(t, ntpOffset, offset)
		assert.True(t, td.checker.IsOutOfSync())
	})

	t.Run("Query error", func(t *testing.T) {
		td := setup(t)
		defer td.Stop()

		td.querier.err = fmt.Errorf("unable to query")

		go td.checker.Start()

		require.Eventually(t, func() bool {
			_, err := td.checker.ClockOffset()

			return err != nil
		}, 5*time.Second, 100*time.Millisecond)

		offset, err := td.checker.ClockOffset()
		assert.Error(t, err)
		assert.Zero(t, offset)
		assert.True(t, td.checker.IsOutOfSync())
	})

	t.Run("Validation error", func(t *testing.T) {
		td := setup(t)
		defer td.Stop()

		td.querier.err = nil
		td.querier.res = &ntp.Response{
			Stratum: 0,
		}

		go td.checker.Start()

		require.Eventually(t, func() bool {
			_, err := td.checker.ClockOffset()

			return err != nil
		}, 5*time.Second, 100*time.Millisecond)

		offset, err := td.checker.ClockOffset()
		assert.Error(t, err)
		assert.Zero(t, offset)
		assert.True(t, td.checker.IsOutOfSync())
	})
}
