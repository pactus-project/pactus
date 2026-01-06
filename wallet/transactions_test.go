package wallet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransactionsListOptions(t *testing.T) {
	tests := []struct {
		opts     []ListTransactionsOption
		expCount int
		expSkip  int
	}{
		{
			opts: []ListTransactionsOption{
				WithCount(-1),
				WithSkip(-1),
			},
			expCount: 10, expSkip: 0,
		},
		{
			opts: []ListTransactionsOption{
				WithCount(0),
				WithSkip(0),
			},
			expCount: 10, expSkip: 0,
		},
		{
			opts: []ListTransactionsOption{
				WithCount(5),
				WithSkip(2),
			},
			expCount: 5, expSkip: 2,
		},
	}

	for _, tt := range tests {
		cfg := defaultListTransactionsConfig
		for _, opt := range tt.opts {
			opt(&cfg)
		}

		assert.Equal(t, tt.expCount, cfg.count)
		assert.Equal(t, tt.expSkip, cfg.skip)
	}
}
