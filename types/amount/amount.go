// This file contains code modified from the btcd project,
// which is licensed under the ISC License.
//
// Original license: https://github.com/btcsuite/btcd/blob/master/LICENSE
//

package amount

import (
	"errors"
	"math"
	"strconv"
)

const (
	// NanoPACPerPAC is the number of NanoPAC in one PAC.
	NanoPACPerPAC = 1e9

	// MaxNanoPAC is the maximum transaction amount allowed in NanoPAC.
	MaxNanoPAC = 42e6 * NanoPACPerPAC
)

// Unit describes a method of converting an Amount to something
// other than the base unit of a Pactus.  The value of the Unit
// is the exponent component of the decadic multiple to convert from
// an amount in Pactus to an amount counted in units.
type Unit int

// These constants define various units used when describing a Pactus
// monetary amount.
const (
	UnitMegaPAC  Unit = 6
	UnitKiloPAC  Unit = 3
	UnitPAC      Unit = 0
	UnitMilliPAC Unit = -3
	UnitMicroPAC Unit = -6
	UnitNanoPAC  Unit = -9
)

// String returns the unit as a string.  For recognized units, the SI
// prefix is used, or "NanoPAC" for the base unit.  For all unrecognized
// units, "1eN PAC" is returned, where N is the AmountUnit.
func (u Unit) String() string {
	switch u {
	case UnitMegaPAC:
		return "MPAC"
	case UnitKiloPAC:
		return "kPAC"
	case UnitPAC:
		return "PAC"
	case UnitMilliPAC:
		return "mPAC"
	case UnitMicroPAC:
		return "μPAC"
	case UnitNanoPAC:
		return "NanoPAC"
	default:
		return "1e" + strconv.FormatInt(int64(u), 10) + " PAC"
	}
}

// Amount represents the atomic unit in Pactus blockchain.
// Each unit equals to 1e-9 of a PAC.
type Amount int64

// round converts a floating point number, which may or may not be representable
// as an integer, to the Amount integer type by rounding to the nearest integer.
// This is performed by adding or subtracting 0.5 depending on the sign, and
// relying on integer truncation to round the value to the nearest Amount.
func round(f float64) Amount {
	if f < 0 {
		return Amount(f - 0.5)
	}

	return Amount(f + 0.5)
}

// NewAmount creates an Amount from a floating-point value representing
// an amount in PAC.  NewAmount returns an error if f is NaN or +-Infinity,
// but it does not check whether the amount is within the total amount of PAC
// producible, as it may not refer to an amount at a single moment in time.
//
// NewAmount is specifically for converting PAC to NanoPAC.
// For creating a new Amount with an int64 value which denotes a quantity of NanoPAC,
// do a simple type conversion from type int64 to Amount.
func NewAmount(f float64) (Amount, error) {
	// The amount is only considered invalid if it cannot be represented
	// as an integer type. This may happen if f is NaN or +-Infinity.
	switch {
	case math.IsNaN(f),
		math.IsInf(f, 1),
		math.IsInf(f, -1):
		return 0, errors.New("invalid PAC amount")
	}

	return round(f * float64(NanoPACPerPAC)), nil
}

// FromString parses a string representing a value in PAC.
// It then uses NewAmount to create an Amount based on the parsed
// floating-point value.
// If the parsing of the string fails, it returns an error.
func FromString(str string) (Amount, error) {
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, err
	}

	return NewAmount(f)
}

// ToUnit converts a monetary amount counted in Pactus base units to a
// floating-point value representing an amount of Pactus (PAC).
func (a Amount) ToUnit(u Unit) float64 {
	return float64(a) / math.Pow10(int(u+9))
}

// ToPAC is equivalent to calling ToUnit with AmountPAC.
func (a Amount) ToPAC() float64 {
	return a.ToUnit(UnitPAC)
}

// ToNanoPAC is equivalent to calling ToUnit with AmountNanoPAC.
// It returns the amount of NanoPAC or atomic unit as a 64-bit integer.
func (a Amount) ToNanoPAC() int64 {
	return int64(a)
}

// Format formats a monetary amount counted in Pactus base units as a
// string for a given unit.  The conversion will succeed for any unit,
// however, known units will be formatted with an appended label describing
// the units with SI notation, and "NanoPAC" for the base unit.
func (a Amount) Format(u Unit) string {
	units := " " + u.String()
	formatted := strconv.FormatFloat(a.ToUnit(u), 'f', -int(u+9), 64)

	return formatted + units
}

// String is the equivalent of calling Format with AmountPAC.
func (a Amount) String() string {
	return a.Format(UnitPAC)
}

// MulF64 multiplies an Amount by a floating point value.
func (a Amount) MulF64(f float64) Amount {
	return round(float64(a) * f)
}
