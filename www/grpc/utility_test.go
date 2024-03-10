package grpc

import (
	"context"
	"testing"

	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/stretchr/testify/assert"
)

func TestCalculateFee(t *testing.T) {
	td := setup(t, nil)
	conn, client := td.utilityClient(t)

	t.Run("Not fixed amount", func(t *testing.T) {
		amount := td.RandAmount()
		res, err := client.CalculateFee(context.Background(),
			&pactus.CalculateFeeRequest{
				Amount:      amount,
				PayloadType: pactus.PayloadType_TRANSFER_PAYLOAD,
				FixedAmount: false,
			})
		assert.NoError(t, err)
		assert.Equal(t, res.Amount, amount)
	})

	t.Run("Fixed amount", func(t *testing.T) {
		amount := td.RandAmount()
		res, err := client.CalculateFee(context.Background(),
			&pactus.CalculateFeeRequest{
				Amount:      amount,
				PayloadType: pactus.PayloadType_TRANSFER_PAYLOAD,
				FixedAmount: true,
			})
		assert.NoError(t, err)
		assert.LessOrEqual(t, res.Amount+res.Fee, amount)
	})

	assert.Nil(t, conn.Close(), "Error closing connection")
	td.StopServer()
}

func TestGetCalculateFee(t *testing.T) {
	td := setup(t, nil)
	conn, client := td.utilityClient(t)

	amount := td.RandAmount()
	res, err := client.CalculateFee(context.Background(),
		&pactus.CalculateFeeRequest{
			Amount:      amount,
			PayloadType: pactus.PayloadType_TRANSFER_PAYLOAD,
		})
	assert.NoError(t, err)
	assert.Equal(t, amount/10000, res.Fee)

	assert.Nil(t, conn.Close(), "Error closing connection")
	td.StopServer()
}
