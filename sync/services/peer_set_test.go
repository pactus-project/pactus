package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServices(t *testing.T) {
	assert.Equal(t, Services(0).String(), "")
	assert.Equal(t, Services(1).String(), "NETWORK")
	assert.Equal(t, Services(2).String(), "2")
	assert.Equal(t, Services(3).String(), "NETWORK | 2")
}
