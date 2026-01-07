package amount_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/util"
	"github.com/stretchr/testify/assert"
	"go.yaml.in/yaml/v2"
)

func TestAmount_Codecs(t *testing.T) {
	tests := []struct {
		name    string
		input   amount.Amount
		expJSON string
		expYAML string
	}{
		{
			name:    "positive amount",
			input:   amount.Amount(123456000000),
			expJSON: "123456000000",
			expYAML: "123.456\n",
		},
		{
			name:    "negative amount",
			input:   amount.Amount(-123456000000),
			expJSON: "-123456000000",
			expYAML: "-123.456\n",
		},
		{
			name:    "zero amount",
			input:   amount.Amount(0),
			expJSON: "0",
			expYAML: "0\n",
		},
		{
			name:    "large amount",
			input:   amount.Amount(123456789000000),
			expJSON: "123456789000000",
			expYAML: "123456.789\n",
		},
		{
			name:    "largest amount",
			input:   amount.Amount(42000000e9),
			expJSON: "42000000000000000",
			expYAML: "4.2e+07\n",
		},
		{
			name:    "smallest unit amount",
			input:   amount.Amount(1),
			expJSON: "1",
			expYAML: "1e-09\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run("JSON marshal and unmarshal", func(t *testing.T) {
				data, err := json.Marshal(tt.input)
				assert.NoError(t, err)
				assert.Equal(t, tt.expJSON, string(data))

				var decoded amount.Amount
				err = json.Unmarshal(data, &decoded)
				assert.NoError(t, err)
				assert.Equal(t, tt.input.ToNanoPAC(), decoded.ToNanoPAC())
			})

			t.Run("YAML marshal and unmarshal", func(t *testing.T) {
				data, err := yaml.Marshal(tt.input)
				assert.NoError(t, err)
				assert.Equal(t, tt.expYAML, string(data))

				var decoded amount.Amount
				err = yaml.Unmarshal(data, &decoded)
				assert.NoError(t, err)
				assert.Equal(t, tt.input.ToNanoPAC(), decoded.ToNanoPAC())
			})

			t.Run("SQL value and scan round trip", func(t *testing.T) {
				val, err := tt.input.Value()
				assert.NoError(t, err)
				assert.IsType(t, int64(0), val)

				var scanned amount.Amount
				err = scanned.Scan(val)
				assert.NoError(t, err)
				assert.Equal(t, tt.input.ToNanoPAC(), scanned.ToNanoPAC())
			})
		})
	}
}

func TestAmount_SQLDriver(t *testing.T) {
	t.Run("scan nil value returns error", func(t *testing.T) {
		var amt amount.Amount
		err := amt.Scan(nil)
		assert.ErrorIs(t, err, amount.ErrInvalidSQLType)
	})

	t.Run("scan unsupported string type returns error", func(t *testing.T) {
		var amt amount.Amount
		err := amt.Scan("123")
		assert.ErrorIs(t, err, amount.ErrInvalidSQLType)
	})

	t.Run("scan unsupported float type returns error", func(t *testing.T) {
		var amt amount.Amount
		err := amt.Scan(123.456)
		assert.ErrorIs(t, err, amount.ErrInvalidSQLType)
	})

	t.Run("value and scan round trip preserves amount", func(t *testing.T) {
		original := amount.Amount(util.RandInt64(1000e9))
		val, err := original.Value()
		assert.NoError(t, err)

		var scanned amount.Amount
		err = scanned.Scan(val)
		assert.NoError(t, err)
		assert.Equal(t, original.ToNanoPAC(), scanned.ToNanoPAC())
	})
}

func TestAmount_JSONMarshaling(t *testing.T) {
	t.Run("unmarshal numeric JSON value", func(t *testing.T) {
		var amt amount.Amount
		err := json.Unmarshal([]byte("123456"), &amt)
		assert.NoError(t, err)
		assert.Equal(t, amount.Amount(123456), amt)
	})

	t.Run("unmarshal fractional numeric JSON returns error", func(t *testing.T) {
		var amt amount.Amount
		err := json.Unmarshal([]byte("1234.567"), &amt)
		assert.Error(t, err)
	})

	t.Run("unmarshal nil input returns error", func(t *testing.T) {
		var amt amount.Amount
		err := json.Unmarshal(nil, &amt)
		assert.Error(t, err)
	})

	t.Run("unmarshal invalid JSON returns error", func(t *testing.T) {
		var amt amount.Amount
		err := json.Unmarshal([]byte(`"invalid"`), &amt)
		assert.Error(t, err)
	})
}

func TestAmount_YAMLMarshaling(t *testing.T) {
	t.Run("unmarshal numeric YAML value", func(t *testing.T) {
		var amt amount.Amount
		err := yaml.Unmarshal([]byte("123456\n"), &amt)
		assert.NoError(t, err)
		assert.Equal(t, amount.Amount(123456), amt)
	})

	t.Run("unmarshal fractional numeric YAML returns error", func(t *testing.T) {
		var amt amount.Amount
		err := yaml.Unmarshal([]byte("1234.567"), &amt)
		assert.NoError(t, err)
		assert.Equal(t, amount.Amount(1234567000000), amt)
	})

	t.Run("unmarshal nil input results in zero value", func(t *testing.T) {
		var amt amount.Amount
		err := yaml.Unmarshal(nil, &amt)
		assert.NoError(t, err)
		assert.Equal(t, amount.Amount(0), amt)
	})

	t.Run("unmarshal invalid input returns error", func(t *testing.T) {
		var amt amount.Amount
		err := yaml.Unmarshal([]byte("invalid"), &amt)
		assert.Error(t, err)
	})

	t.Run("unmarshal infinite value returns error", func(t *testing.T) {
		var amt amount.Amount
		err := yaml.Unmarshal([]byte(".inf"), &amt)
		assert.Error(t, err)
	})

	t.Run("unmarshal error", func(t *testing.T) {
		var amt amount.Amount
		err := amt.UnmarshalYAML(func(any) error {
			return errors.New("error")
		})
		assert.Error(t, err)
	})
}
