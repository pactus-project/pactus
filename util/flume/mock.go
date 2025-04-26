package flume

// MockPipeline implements the Pipeline interface for testing
type MockPipeline[T any] struct {
	ch chan T
}

// MockingPipeline creates a new mock pipeline instance
func MockingPipeline[T any]() *MockPipeline[T] {
	return &MockPipeline[T]{
		ch: make(chan T, 100),
	}
}

func (m *MockPipeline[T]) Name() string {
	return "mocked pipeline"
}

func (m *MockPipeline[T]) Send(data T) {
	m.ch <- data
}

func (m *MockPipeline[T]) RegisterReceiver(fn func(T)) {
}

func (m *MockPipeline[T]) Close() {
}

func (m *MockPipeline[T]) IsClosed() bool {
	return false
}

func (m *MockPipeline[T]) Channel() <-chan T {
	return m.ch
}
