package amount

import (
	"database/sql/driver"
	"encoding/json"
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

// MarshalJSON implements the json.Marshaler interface.
// It serializes the Amount as a string in PAC format.
func (a Amount) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.Format())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It expects a string value in PAC format.
func (a *Amount) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	amt, err := FromString(str)
	if err != nil {
		return err
	}

	*a = amt

	return nil
}

// MarshalYAML implements the yaml.Marshaler interface.
// It serializes the Amount as a string in PAC format.
func (a Amount) MarshalYAML() (any, error) {
	return a.Format(), nil
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
// It expects a string value in PAC format.
func (a *Amount) UnmarshalYAML(unmarshal func(any) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return err
	}

	amt, err := FromString(str)
	if err != nil {
		return err
	}

	*a = amt

	return nil
}
