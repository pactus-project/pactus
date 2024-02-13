package wallet2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDB(t *testing.T) {
	someDB, err := newDB(":memory:")

	assert.Nil(t, err)
	assert.NotNil(t, someDB)
}
