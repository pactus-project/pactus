package amount_test

import (
	"math/rand"
	"testing"

	"github.com/pactus-project/pactus/types/amount"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// An amount is parsed to the exact number of NanoPAC, also when it has more digits
// than a floating-point value can hold.
func TestFromStringIsExact(t *testing.T) {
	tests := []struct {
		str     string
		nanoPAC int64
	}{
		{"7000000.001", 7_000_000_001_000_000},
		{"10000000.05", 10_000_000_050_000_000},
		{"41999999.999999999", 41_999_999_999_999_999},
		{"42000000", 42_000_000_000_000_000},
		{"-41999999.999999999", -41_999_999_999_999_999},
		{"999999999.999999999", 999_999_999_999_999_999},
		{"0.000000001", 1},
		{"-0.000000001", -1},
		{"000.100", 100_000_000},
		// The digits after the ninth decimal are rounded to the nearest NanoPAC.
		{"0.0000000004", 0},
		{"0.0000000005", 1},
		{"0.00000000149999", 1},
		{"0.0000000015", 2},
		{"-0.0000000015", -2},
		{"41999999.9999999994", 41_999_999_999_999_999},
		{"41999999.9999999995", 42_000_000_000_000_000},
	}

	for _, tt := range tests {
		amt, err := amount.FromString(tt.str)
		require.NoError(t, err, tt.str)
		assert.Equal(t, tt.nanoPAC, amt.ToNanoPAC(), tt.str)
	}
}

func TestFromStringBounds(t *testing.T) {
	for _, str := range []string{"1000000000", "-1000000000", "999999999.9999999995", "9999999999.5"} {
		_, err := amount.FromString(str)
		require.Error(t, err, str)
	}
}

// An amount that is formatted without delimiters is parsed back to the same amount.
func TestFormatAndParseBack(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))

	for range 100_000 {
		amt := amount.Amount(rnd.Int63n(int64(amount.MaxNanoPAC) + 1))
		if rnd.Intn(2) == 0 {
			amt = -amt
		}

		parsed, err := amount.FromString(amt.Format(amount.WithUnit(amount.UnitPAC)))
		require.NoError(t, err)
		require.Equal(t, amt, parsed, "%d NanoPAC is formatted as %s", int64(amt), amt.Format())
	}
}

// An amount is written with all its digits, in each unit, with and without delimiters.
func TestFormatIsExact(t *testing.T) {
	amt := amount.Amount(33_006_791_947_779_278)

	assert.Equal(t, "33006791.947779278", amt.Format())
	assert.Equal(t, "33,006,791.947779278 PAC", amt.String())
	assert.Equal(t, "-33,006,791.947779278 PAC", (-amt).String())
	assert.Equal(t, "33006791947779278 NanoPAC", amt.Format(amount.WithUnit(amount.UnitNanoPAC)))
	assert.Equal(t, "33006791947.779278 mPAC", amt.Format(amount.WithUnit(amount.UnitMilliPAC)))
	assert.Equal(t, "33.006791947779278 MPAC", amt.Format(amount.WithUnit(amount.UnitMegaPAC)))
	assert.Equal(t, "33,006,791,947,779,278 NanoPAC",
		amt.Format(amount.WithUnit(amount.UnitNanoPAC), amount.WithDelimiters()))

	assert.Equal(t, "0 PAC", amount.Amount(0).String())
	assert.Equal(t, "0.000000001 PAC", amount.Amount(1).String())
	assert.Equal(t, "-0.000000001 PAC", amount.Amount(-1).String())
	assert.Equal(t, "0.000000000001 MPAC", amount.Amount(1000).Format(amount.WithUnit(amount.UnitMegaPAC)))
	assert.Equal(t, "42,000,000 PAC", amount.Amount(amount.MaxNanoPAC).String())
}
