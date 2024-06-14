package wallet_test

// func TestActivitiesSort(t *testing.T) {
// 	historyMock := history{
// 		Activities:   map[string][]activity{},
// 		Transactions: map[string]transaction{},
// 	}

// 	historyMock.addActivity("addr-1", 20, &pactus.GetTransactionResponse{
// 		BlockTime: 2, Transaction: &pactus.TransactionInfo{
// 			Id: []byte{2},
// 		},
// 	})
// 	historyMock.addActivity("addr-1", 40, &pactus.GetTransactionResponse{
// 		BlockTime: 4, Transaction: &pactus.TransactionInfo{
// 			Id: []byte{4},
// 		},
// 	})
// 	historyMock.addActivity("addr-1", 30, &pactus.GetTransactionResponse{
// 		BlockTime: 3, Transaction: &pactus.TransactionInfo{
// 			Id: []byte{3},
// 		},
// 	})
// 	historyMock.addActivity("addr-1", 10, &pactus.GetTransactionResponse{
// 		BlockTime: 1, Transaction: &pactus.TransactionInfo{
// 			Id: []byte{1},
// 		},
// 	})
// 	historyMock.addActivity("addr-2", 50, &pactus.GetTransactionResponse{
// 		BlockTime: 5, Transaction: &pactus.TransactionInfo{
// 			Id: []byte{5},
// 		},
// 	})

// 	h := historyMock.getAddrHistory("addr-1")
// 	assert.Len(t, h, 4)
// 	assert.Equal(t, h[0].TxID, "01")
// 	assert.Equal(t, h[1].TxID, "02")
// 	assert.Equal(t, h[2].TxID, "03")
// 	assert.Equal(t, h[3].TxID, "04")
// }
