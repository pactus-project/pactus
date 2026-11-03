package session

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatsStringOnAValue(t *testing.T) {
	stats := Stats{Total: 4, Open: 1, Completed: 2, Uncompleted: 1}

	// A value, not only a pointer, formats with String.
	assert.Equal(t, "total: 4, open: 1, completed: 2, uncompleted: 1", stats.String())
}
