package txpool

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfigCheck(t *testing.T) {
	c := DefaultConfig()
	assert.NoError(t, c.SanityCheck())

	assert.Equal(t,
		c.sendPoolSize()+
		c.bondPoolSize()+
		c.unbondPoolSize()+
		c.withdrawPoolSize()+
		c.sortitionPoolSize(), c.MaxSize)

	c.MaxSize = 0
	assert.Error(t, c.SanityCheck())
}
