package metric

import "github.com/pactus-project/pactus/sync/bundle/message"

type Counter struct {
	Bytes   int64
	Bundles int64
}

type Metric struct {
	TotalInvalid    Counter
	TotalSent       Counter
	TotalReceived   Counter
	MessageSent     map[message.Type]*Counter
	MessageReceived map[message.Type]*Counter
}

func NewMetric() Metric {
	return Metric{
		MessageSent:     make(map[message.Type]*Counter),
		MessageReceived: make(map[message.Type]*Counter),
	}
}

func (m *Metric) UpdateSentMetric(msgType message.Type, bytes int64) {
	m.TotalSent.Bundles++
	m.TotalSent.Bytes += bytes

	_, ok := m.MessageSent[msgType]
	if !ok {
		m.MessageSent[msgType] = &Counter{}
	}

	m.MessageSent[msgType].Bundles++
	m.MessageSent[msgType].Bytes += bytes
}

func (m *Metric) UpdateReceivedMetric(msgType message.Type, bytes int64) {
	m.TotalReceived.Bundles++
	m.TotalReceived.Bytes += bytes

	_, ok := m.MessageReceived[msgType]
	if !ok {
		m.MessageReceived[msgType] = &Counter{}
	}

	m.MessageReceived[msgType].Bundles++
	m.MessageReceived[msgType].Bytes += bytes
}

func (m *Metric) UpdateInvalidMetric(bytes int64) {
	m.TotalInvalid.Bundles++
	m.TotalInvalid.Bytes += bytes
}
