package txpool

import (
	"fmt"

	"github.com/pactus-project/pactus/types/amount"
)

// ConfigError is returned when the transaction pool configuration is invalid.
type ConfigError struct {
	Reason string
}

func (e ConfigError) Error() string {
	return e.Reason
}

// InvalidFeeError indicates that a transaction fee is below the required minimum.
type InvalidFeeError struct {
	MinimumFee amount.Amount
}

func (e InvalidFeeError) Error() string {
	return fmt.Sprintf("transaction fee is below the minimum of %s", e.MinimumFee)
}
