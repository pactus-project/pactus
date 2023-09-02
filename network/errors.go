package network

import "fmt"

// NotSubscribedError is returned when the peer is not subscribed to a
// specific topic.
type NotSubscribedError struct {
	TopicID TopicID
}

func (e NotSubscribedError) Error() string {
	return fmt.Sprintf("Not subscribed to the '%s' topic", e.TopicID.String())
}

// InvalidTopicError is returned when the Pub-Sub topic is invalid.
type InvalidTopicError struct {
	TopicID TopicID
}

func (e InvalidTopicError) Error() string {
	return fmt.Sprintf("invalid topic: %s",
		e.TopicID.String())
}

// LibP2PError is returned when an underlying libp2p operation encounters an error.
type LibP2PError struct {
	Err error
}

func (e LibP2PError) Error() string {
	return fmt.Sprintf("libp2p error: %s",
		e.Err.Error())
}
