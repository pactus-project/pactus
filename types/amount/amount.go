// This file contains code modified from the btcd project,
// which is licensed under the ISC License.
//
// Original license: https://github.com/btcsuite/btcd/blob/master/LICENSE
//

package amount

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/encoding"
)

const (
	// NanoPACPerPAC is the number of NanoPAC in one PAC.
	NanoPACPerPAC = 1e9

	// MaxNanoPAC is the maximum transaction amount allowed in NanoPAC.
	MaxNanoPAC = 42e6 * NanoPACPerPAC
)

// Unit describes a method of converting an Amount to something
// other than the base unit of PAC.  The value of the Unit
// is the exponent component of the decimal multiple to convert from
// an amount in PAC to an amount counted in units.
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
// units, "1eN PAC" is returned, where N is the Unit.
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
// Each unit equals 1e-9 of a PAC.
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
// an amount in PAC.  NewAmount returns an error if f is NaN, +-Infinity or too large to be
// converted safely, but it does not check whether the amount is within the total amount of PAC
// producible, as it may not refer to an amount at a single moment in time.
//
// NewAmount is specifically for converting PAC to NanoPAC.
// For creating a new Amount with an int64 value which denotes a quantity of NanoPAC,
// do a simple type conversion from type int64 to Amount.
func NewAmount(pac float64) (Amount, error) {
	// The amount is only considered invalid if it cannot be represented
	// as an integer type. This may happen if f is NaN or +-Infinity.
	switch {
	case math.IsNaN(pac),
		math.IsInf(pac, 1),
		math.IsInf(pac, -1):
		return 0, errors.New("invalid PAC amount")
	}

	nanoPAC := pac * float64(NanoPACPerPAC)

	// Converting a float64 outside the int64 range is implementation-defined
	// (saturates on arm64, wraps to MinInt64 on amd64). maxSafeNanoPAC is far
	// above MaxNanoPAC (the real total supply) but far enough below the
	// int64/float64 range limits to make that conversion always well-defined.
	const maxSafeNanoPAC = 1e18
	if nanoPAC >= maxSafeNanoPAC || nanoPAC <= -maxSafeNanoPAC {
		return 0, errors.New("PAC amount out of range")
	}

	return round(nanoPAC), nil
}

// FromString parses a string representing a value in PAC.
// It then uses NewAmount to create an Amount based on the parsed
// floating-point value.
// If the parsing of the string fails, it returns an error.
//
// The value is a decimal number, like "1", "-0.5" or "12.25 PAC".
// The other forms that a floating-point parser accepts are rejected,
// because they are easy to misread in an amount: an exponent ("1e3"),
// a digit separator ("1_0"), a hexadecimal number, "inf" or "nan".
func FromString(str string) (Amount, error) {
	str = strings.TrimSpace(str)
	str = strings.TrimSpace(strings.TrimSuffix(str, "PAC"))
	if !isDecimal(str) {
		return 0, fmt.Errorf("invalid PAC amount %q: %w", str, strconv.ErrSyntax)
	}

	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, err
	}

	return NewAmount(f)
}

// isDecimal reports whether str is an optional sign, digits and an optional fraction.
func isDecimal(str string) bool {
	if strings.HasPrefix(str, "+") || strings.HasPrefix(str, "-") {
		str = str[1:]
	}
	if strings.Count(str, ".") > 1 {
		return false
	}

	digits := strings.Replace(str, ".", "", 1)
	if digits == "" {
		return false
	}
	for _, char := range digits {
		if char < '0' || char > '9' {
			return false
		}
	}

	return true
}

// ToUnit converts a monetary amount counted in NanoPAC (base units) to a
// floating-point value representing an amount in the specified unit.
func (a Amount) ToUnit(u Unit) float64 {
	return float64(a) / math.Pow10(int(u+9))
}

// ToPAC is equivalent to calling ToUnit with UnitPAC.
func (a Amount) ToPAC() float64 {
	return a.ToUnit(UnitPAC)
}

// ToNanoPAC returns the amount of NanoPAC (atomic units) as a 64-bit integer.
func (a Amount) ToNanoPAC() int64 {
	return int64(a)
}

// formatOptions holds options for formatting amounts.
type formatOptions struct {
	unit           Unit
	WithUnit       bool
	withDelimiters bool
}

// FormatOption is a function that configures formatting options.
type FormatOption func(*formatOptions)

// WithUnit sets the unit for formatting.
func WithUnit(u Unit) FormatOption {
	return func(opts *formatOptions) {
		opts.unit = u
		opts.WithUnit = true
	}
}

// WithDelimiters enables delimiter formatting (e.g., "1,234.56").
func WithDelimiters() FormatOption {
	return func(opts *formatOptions) {
		opts.withDelimiters = true
	}
}

// Format formats a monetary amount counted in Pactus base units as a
// string for a given unit.  The conversion will succeed for any unit,
// however, known units will be formatted with an appended label describing
// the units with SI notation, and "NanoPAC" for the base unit.
//
// Format accepts options:
//   - WithUnit(u Unit): sets the unit for formatting (defaults to no unit)
//   - WithDelimiters(): enables comma delimiters in the formatted number (defaults to false)
func (a Amount) Format(opts ...FormatOption) string {
	options := &formatOptions{}

	// Apply options
	for _, opt := range opts {
		opt(options)
	}

	precision := -int(options.unit + 9)
	value := a.ToUnit(options.unit)

	var formatted string
	if options.withDelimiters {
		formatted = util.FormatFloatWithDelimiters(value, precision)
	} else {
		formatted = strconv.FormatFloat(value, 'f', precision, 64)
	}

	if options.WithUnit {
		formatted += " " + options.unit.String()
	}

	return formatted
}

// String is the equivalent of calling Format with UnitPAC.
func (a Amount) String() string {
	return a.Format(WithUnit(UnitPAC), WithDelimiters())
}

// MulF64 multiplies an Amount by a floating point value.
func (a Amount) MulF64(f float64) Amount {
	return round(float64(a) * f)
}

// Scan implements the sql.Scanner interface for SQL database operations.
// It accepts int64 values representing NanoPAC and converts them to Amount.
func (a *Amount) Scan(src any) error {
	switch v := src.(type) {
	case int64:
		*a = Amount(v)

		return nil
	default:
		return ErrInvalidSQLType
	}
}

// Value implements the driver.Valuer interface for SQL database operations.
// It returns the Amount as NanoPAC (int64) to avoid floating-point precision issues.
func (a Amount) Value() (driver.Value, error) {
	return int64(a), nil
}

// Encode writes the raw bytes of the Amount to the provided writer.
func (a Amount) Encode(w io.Writer) error {
	return encoding.WriteVarInt(w, uint64(a))
}

// Decode reads the raw bytes of the Amount from the provided reader and initializes the Amount.
func (a *Amount) Decode(r io.Reader) error {
	val, err := encoding.ReadVarInt(r)
	if err != nil {
		return err
	}

	*a = Amount(val)

	return nil
}
