package wallet

import (
	"testing"

	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/stretchr/testify/assert"
)

func TestActivitiesSort(t *testing.T) {
	history := history{
		Activities:   map[string][]activity{},
		Transactions: map[string]transaction{},
	}

	history.addActivity("addr-1", 20, &pactus.TransactionResponse{
		BlockTime: 2, Transaction: &pactus.TransactionInfo{
			Id: []byte{2}}})
	history.addActivity("addr-1", 40, &pactus.TransactionResponse{
		BlockTime: 4, Transaction: &pactus.TransactionInfo{
			Id: []byte{4}}})
	history.addActivity("addr-1", 30, &pactus.TransactionResponse{
		BlockTime: 3, Transaction: &pactus.TransactionInfo{
			Id: []byte{3}}})
	history.addActivity("addr-1", 10, &pactus.TransactionResponse{
		BlockTime: 1, Transaction: &pactus.TransactionInfo{
			Id: []byte{1}}})
	history.addActivity("addr-2", 50, &pactus.TransactionResponse{
		BlockTime: 5, Transaction: &pactus.TransactionInfo{
			Id: []byte{5}}})

	h := history.getAddrHistory("addr-1")
	assert.Len(t, h, 4)
	assert.Equal(t, h[0].TxID, "01")
	assert.Equal(t, h[1].TxID, "02")
	assert.Equal(t, h[2].TxID, "03")
	assert.Equal(t, h[3].TxID, "04")
}
