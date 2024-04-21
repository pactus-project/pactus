package txpool

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfigCheck(t *testing.T) {
	c := DefaultConfig()
	assert.NoError(t, c.BasicCheck())

	assert.Equal(t, 600, c.transferPoolSize())
	assert.Equal(t, 100, c.bondPoolSize())
	assert.Equal(t, 100, c.unbondPoolSize())
	assert.Equal(t, 100, c.withdrawPoolSize())
	assert.Equal(t, 100, c.sortitionPoolSize())

	assert.Equal(t,
		c.transferPoolSize()+
			c.bondPoolSize()+
			c.unbondPoolSize()+
			c.withdrawPoolSize()+
			c.sortitionPoolSize(), c.MaxSize)
}

func TestInvalidConfig(t *testing.T) {
	tests := []struct {
		conf   Config
		ErrStr string
	}{
		{
			conf: Config{
				MaxSize:     0,
				MinValuePAC: 0.1,
			},
			ErrStr: "maxSize can't be less than 10",
		},

		{
			conf: Config{
				MaxSize:     9,
				MinValuePAC: 0.1,
			},
			ErrStr: "maxSize can't be less than 10",
		},

		{
			conf: Config{
				MaxSize:     100,
				MinValuePAC: 1.1,
			},
			ErrStr: "minVale can't be greater than 1 PAC",
		},

		{
			conf: Config{
				MaxSize:     100,
				MinValuePAC: 1.0,
			},
			ErrStr: "",
		},
	}

	for _, test := range tests {
		err := test.conf.BasicCheck()
		if test.ErrStr != "" {
			assert.ErrorContains(t, err, test.ErrStr)
		} else {
			assert.NoError(t, err)
		}
	}
}
