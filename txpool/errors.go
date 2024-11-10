package txpool

import (
	"fmt"

	"github.com/pactus-project/pactus/types/amount"
)

// ConfigError is returned when the txPool configuration is invalid.
type ConfigError struct {
	Reason string
}

func (e ConfigError) Error() string {
	return e.Reason
}

// InvalidFeeError indicates that the transaction fee is below the minimum required.
type InvalidFeeError struct {
	MinimumFee amount.Amount
}

func (e InvalidFeeError) Error() string {
	return fmt.Sprintf("transaction fee is below the minimum of %s", e.MinimumFee)
}
