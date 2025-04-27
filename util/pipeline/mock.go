package pipeline

// MockPipeline implements the Pipeline interface for testing.
type MockPipeline[T any] struct {
	ch       chan T
	receiver func(T)
}

// MockingPipeline creates a new mock pipeline instance.
func MockingPipeline[T any]() *MockPipeline[T] {
	return &MockPipeline[T]{
		ch: make(chan T, 100),
	}
}

func (*MockPipeline[T]) Name() string {
	return "mocked pipeline"
}

func (m *MockPipeline[T]) Send(data T) {
	if m.receiver != nil {
		m.receiver(data)
	}

	if m.ch != nil {
		m.ch <- data
	}
}

func (m *MockPipeline[T]) RegisterReceiver(fn func(T)) {
	m.receiver = fn
}

func (*MockPipeline[T]) Close() {
}

func (*MockPipeline[T]) IsClosed() bool {
	return false
}

func (m *MockPipeline[T]) UnsafeGetChannel() <-chan T {
	return m.ch
}
