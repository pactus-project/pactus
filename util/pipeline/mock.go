package pipeline

import "context"

// MockPipeline implements the Pipeline interface for testing.
type MockPipeline[T any] struct {
	Pipeline[T]
}

// MockingPipeline creates a new mock pipeline instance.
func MockingPipeline[T any]() *MockPipeline[T] {
	return &MockPipeline[T]{
		Pipeline: New[T](context.TODO(), "test", 100),
	}
}
