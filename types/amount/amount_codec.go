package amount

import (
	"database/sql/driver"
)

// Value implements the driver.Valuer interface for SQL database operations.
// It returns the Amount as NanoPAC (int64) to avoid floating-point precision issues.
func (a Amount) Value() (driver.Value, error) {
	return int64(a), nil
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

// MarshalYAML implements the yaml.Marshaler interface.
// It serializes the Amount as a float in PAC units.
func (a Amount) MarshalYAML() (any, error) {
	// Emit PAC value (human-readable float)
	return a.ToPAC(), nil
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
//
// It accepts either:
//   - a float value in PAC format (e.g., 123.456), or
//   - an integer value in NanoPAC format (e.g., 123456000000).
func (a *Amount) UnmarshalYAML(unmarshal func(any) error) error {
	var v any
	if err := unmarshal(&v); err != nil {
		return err
	}

	switch val := v.(type) {
	case float64:
		// PAC format
		amt, err := NewAmount(val)
		if err != nil {
			return err
		}
		*a = amt

	case int:
		// NanoPAC format (small integers)
		*a = Amount(val)

	default:
		return ErrInvalidYAMLType
	}

	return nil
}
