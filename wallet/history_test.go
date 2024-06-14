package wallet_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHistory(t *testing.T) {
	td := setup(t)
	defer td.Close()

	history := td.wallet.GetHistory(td.RandAccAddress().String())
	assert.Empty(t, history)
}
