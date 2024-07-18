package executor

// ExecutionLevel determines the level of strictness during execution.
type ExecutionLevel int

const (
	// In Commit level the transaction is executed without any validation.
	Commit ExecutionLevel = 0

	// Check is less restrictive than Validation, allowing transactions to remain
	// in the pool even if they're currently invalid, as they may become valid
	// in future blocks.
	Check ExecutionLevel = 2

	// Validate is the strictest level, ensuring full compliance with rules and
	// conditions at consensus time.
	Validate ExecutionLevel = 1
)
