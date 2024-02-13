package wallet2

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDB(t *testing.T) {
	someDB, err := newDB("wallet.db")
	// someDB, err := newDB(":memory:")

	assert.Nil(t, err)
	assert.NotNil(t, someDB)

	someDB.CreateTables()
	err = someDB.InsertIntoAddress("test", "test", "test", "test")
	log.Println(err)
	tran := &Transaction{
		TxID:        "ssf",
		BlockHeight: 2,
		BlockTime:   3,
		PayloadType: "sf",
		Data:        "sf",
		Description: "sfd",
		Amount:      3,
		Status:      1,
	}
	err = someDB.InsertIntoTransaction(tran)
	log.Println(err)
}
