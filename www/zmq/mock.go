package zmq

import "github.com/pactus-project/pactus/types/block"

type MockPublisher struct {
	MockAddress   string
	MockTopicName string
	MockHwm       int
}

var _ Publisher = &MockPublisher{}

func MockingPublisher(address, topicName string, hwm int) *MockPublisher {
	return &MockPublisher{
		MockAddress:   address,
		MockTopicName: topicName,
		MockHwm:       hwm,
	}
}

func (m *MockPublisher) Address() string {
	return m.MockAddress
}

func (m *MockPublisher) TopicName() string {
	return m.MockTopicName
}

func (m *MockPublisher) HWM() int {
	return m.MockHwm
}

func (*MockPublisher) onNewBlock(*block.Block) {
}
