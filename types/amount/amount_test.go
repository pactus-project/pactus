// This file contains code modified from the btcd project,
// which is licensed under the ISC License.
//
// Original license: https://github.com/btcsuite/btcd/blob/master/LICENSE
//

package amount_test

import (
	"math"
	"strconv"
	"testing"

	"github.com/pactus-project/pactus/types/amount"
	"github.com/stretchr/testify/assert"
)

func TestAmountCreation(t *testing.T) {
	tests := []struct {
		name     string
		amount   float64
		valid    bool
		expected amount.Amount
	}{
		// Positive tests.
		{
			name:     "zero",
			amount:   0,
			valid:    true,
			expected: 0,
		},
		{
			name:     "max producible",
			amount:   42e6,
			valid:    true,
			expected: amount.MaxNanoPAC,
		},
		{
			name:     "min producible",
			amount:   -42e6,
			valid:    true,
			expected: -amount.MaxNanoPAC,
		},
		{
			name:     "exceeds max producible",
			amount:   42e6 + 8e-9,
			valid:    true,
			expected: amount.MaxNanoPAC + 8,
		},
		{
			name:     "exceeds min producible",
			amount:   -42e6 - 8e-9,
			valid:    true,
			expected: -amount.MaxNanoPAC - 8,
		},
		{
			name:     "one hundred",
			amount:   100,
			valid:    true,
			expected: 100 * amount.NanoPACPerPAC,
		},
		{
			name:     "fraction",
			amount:   0.012345678,
			valid:    true,
			expected: 12345678,
		},
		{
			name:     "rounding up",
			amount:   54.999999999999943157,
			valid:    true,
			expected: 55 * amount.NanoPACPerPAC,
		},
		{
			name:     "rounding down",
			amount:   55.000000000000056843,
			valid:    true,
			expected: 55 * amount.NanoPACPerPAC,
		},

		// Negative tests.
		{
			name:   "not-a-number",
			amount: math.NaN(),
			valid:  false,
		},
		{
			name:   "-infinity",
			amount: math.Inf(-1),
			valid:  false,
		},
		{
			name:   "+infinity",
			amount: math.Inf(1),
			valid:  false,
		},
	}

	for _, test := range tests {
		amt, err := amount.NewAmount(test.amount)
		if test.valid {
			assert.NoErrorf(t, err,
				"%v: Positive test Amount creation failed with: %v", test.name, err)
		} else {
			assert.Errorf(t, err,
				"%v: Negative test Amount creation succeeded (value %v) when should fail", test.name, amt)
		}

		assert.Equal(t, amt, test.expected,
			"%v: Created amount %v does not match expected %v", test.name, amt, test.expected)
	}
}

func TestAmountUnitConversions(t *testing.T) {
	tests := []struct {
		name      string
		amount    amount.Amount
		unit      amount.Unit
		converted float64
		str       string
	}{
		{
			name:      "MPAC",
			amount:    amount.MaxNanoPAC,
			unit:      amount.UnitMegaPAC,
			converted: 42,
			str:       "42 MPAC",
		},
		{
			name:      "kPAC",
			amount:    444_333_222_111_000,
			unit:      amount.UnitKiloPAC,
			converted: 444.333_222_111_000,
			str:       "444.333222111 kPAC",
		},
		{
			name:      "PAC",
			amount:    444_333_222_111_000,
			unit:      amount.UnitPAC,
			converted: 444_333.222_111,
			str:       "444333.222111 PAC",
		},
		{
			name:      "a thousand NanoPAC as PAC",
			amount:    1_000,
			unit:      amount.UnitPAC,
			converted: 0.000_001,
			str:       "0.000001 PAC",
		},
		{
			name:      "a single NanoPAC as PAC",
			amount:    1,
			unit:      amount.UnitPAC,
			converted: 0.000_000_001,
			str:       "0.000000001 PAC",
		},
		{
			name:      "amount with trailing zero but no decimals",
			amount:    10_000_000_000,
			unit:      amount.UnitPAC,
			converted: 10,
			str:       "10 PAC",
		},
		{
			name:      "mPAC",
			amount:    444_333_222_111_000,
			unit:      amount.UnitMilliPAC,
			converted: 444_333_222.111_000,
			str:       "444333222.111 mPAC",
		},
		{
			name:      "μPAC",
			amount:    444_333_222_111_000,
			unit:      amount.UnitMicroPAC,
			converted: 444_333_222_111.000,
			str:       "444333222111 μPAC",
		},
		{
			name:      "NanoPAC",
			amount:    444_333_222_111_000,
			unit:      amount.UnitNanoPAC,
			converted: 444_333_222_111_000,
			str:       "444333222111000 NanoPAC",
		},
		{
			name:      "non-standard unit",
			amount:    444_333_222_111_000,
			unit:      amount.Unit(-1),
			converted: 4_443_332.221_110_00,
			str:       "4443332.22111 1e-1 PAC",
		},
	}

	for _, test := range tests {
		f := test.amount.ToUnit(test.unit)
		assert.Equal(t, f, test.converted,
			"%v: converted value %v does not match expected %v", test.name, f, test.converted)

		str := test.amount.Format(test.unit)
		assert.Equal(t, str, test.str,
			"%v: format '%v' does not match expected '%v'", test.name, str, test.str)

		// Verify that Amount.ToPAC works as advertised.
		f1 := test.amount.ToUnit(amount.UnitPAC)
		f2 := test.amount.ToPAC()
		assert.Equal(t, f1, f2,
			"%v: ToPAC does not match ToUnit(AmountPAC): %v != %v", test.name, f1, f2)

		// Verify that Amount.String works as advertised.
		s1 := test.amount.Format(amount.UnitPAC)
		s2 := test.amount.String()
		assert.Equal(t, s1, s2,
			"%v: String does not match Format(AmountPac): %v != %v", test.name, s1, s2)
	}
}

func TestAmountMulF64(t *testing.T) {
	tests := []struct {
		name string
		amt  amount.Amount
		mul  float64
		res  amount.Amount
	}{
		{
			name: "Multiply 0.1 PAC by 2",
			amt:  100e6, // 0.1 PAC
			mul:  2,
			res:  200e6, // 0.2 PAC
		},
		{
			name: "Multiply 0.2 PAC by 0.02",
			amt:  200e6, // 0.2 PAC
			mul:  1.02,
			res:  204e6, // 0.204 PAC
		},
		{
			name: "Multiply 0.1 PAC by -2",
			amt:  100e6, // 0.1 PAC
			mul:  -2,
			res:  -200e6, // -0.2 PAC
		},
		{
			name: "Multiply 0.2 PAC by -0.02",
			amt:  200e6, // 0.2 PAC
			mul:  -1.02,
			res:  -204e6, // -0.204 PAC
		},
		{
			name: "Multiply -0.1 PAC by 2",
			amt:  -100e6, // -0.1 PAC
			mul:  2,
			res:  -200e6, // -0.2 PAC
		},
		{
			name: "Multiply -0.2 PAC by 0.02",
			amt:  -200e6, // -0.2 PAC
			mul:  1.02,
			res:  -204e6, // -0.204 PAC
		},
		{
			name: "Multiply -0.1 PAC by -2",
			amt:  -100e6, // -0.1 PAC
			mul:  -2,
			res:  200e6, // 0.2 PAC
		},
		{
			name: "Multiply -0.2 PAC by -0.02",
			amt:  -200e6, // -0.2 PAC
			mul:  -1.02,
			res:  204e6, // 0.204 PAC
		},
		{
			name: "Round down",
			amt:  49, // 49 NanoPACs
			mul:  0.01,
			res:  0,
		},
		{
			name: "Round up",
			amt:  50, // 50 NanoPACs
			mul:  0.01,
			res:  1, // 1 NanoPAC
		},
		{
			name: "Multiply by 0",
			amt:  1e9, // 1 PAC
			mul:  0,
			res:  0, // 0 PAC
		},
		{
			name: "Multiply 1 by 0.5",
			amt:  1, // 1 NanoPAC
			mul:  0.5,
			res:  1, // 1 NanoPAC
		},
		{
			name: "Multiply 100 by 66%",
			amt:  100, // 100 NanoPACs
			mul:  0.66,
			res:  66, // 66 NanoPACs
		},
		{
			name: "Multiply 100 by 66.6%",
			amt:  100, // 100 NanoPACs
			mul:  0.666,
			res:  67, // 67 NanoPACs
		},
		{
			name: "Multiply 100 by 2/3",
			amt:  100, // 100 NanoPACs
			mul:  2.0 / 3,
			res:  67, // 67 NanoPACs
		},
	}

	for _, test := range tests {
		a := test.amt.MulF64(test.mul)
		if a != test.res {
			t.Errorf("%v: expected %v got %v", test.name, test.res, a)
		}
	}
}

func TestCoinToChangeConversion(t *testing.T) {
	tests := []struct {
		amount  string
		PAC     float64
		NanoPac int64
		str     string
		parsErr error
	}{
		{"0", 0, 0, "0 PAC", nil},
		{"1", 1, 1000000000, "1 PAC", nil},
		{"123.123", 123.123, 123123000000, "123.123 PAC", nil},
		{"123.0123", 123.0123, 123012300000, "123.0123 PAC", nil},
		{"123.01230", 123.0123, 123012300000, "123.0123 PAC", nil},
		{"123.000123", 123.000123, 123000123000, "123.000123 PAC", nil},
		{"123.000000123", 123.000000123, 123000000123, "123.000000123 PAC", nil},
		{"-123.000000123", -123.000000123, -123000000123, "-123.000000123 PAC", nil},
		{"0123.000000123", 123.000000123, 123000000123, "123.000000123 PAC", nil},
		{"+123.000000123", 123.000000123, 123000000123, "123.000000123 PAC", nil},
		{"123.0000001234", 123.000000123, 123000000123, "123.000000123 PAC", nil},
		{"1coin", 0, 0, "0", strconv.ErrSyntax},
	}
	for _, test := range tests {
		amt, err := amount.FromString(test.amount)
		if test.parsErr == nil {
			assert.NoError(t, err)
			assert.Equal(t, amt.ToNanoPAC(), test.NanoPac)
			assert.Equal(t, amt.ToPAC(), test.PAC)
			assert.Equal(t, amt.String(), test.str)
		} else {
			assert.ErrorIs(t, err, test.parsErr)
		}
	}
}
