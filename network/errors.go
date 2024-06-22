package network

import "fmt"

// ConfigError is returned when the config is not valid with a descriptive Reason message.
type ConfigError struct {
	Reason string
}

func (e ConfigError) Error() string {
	return e.Reason
}

// NotSubscribedError is returned when the peer is not subscribed to a
// specific topic.
type NotSubscribedError struct {
	TopicID TopicID
}

func (e NotSubscribedError) Error() string {
	return fmt.Sprintf("not subscribed to the '%s' topic", e.TopicID.String())
}

// InvalidTopicError is returned when the Pub-Sub topic is invalid.
type InvalidTopicError struct {
	TopicID TopicID
}

func (e InvalidTopicError) Error() string {
	return fmt.Sprintf("invalid topic: '%s'", e.TopicID.String())
}

// LibP2PError is returned when an underlying libp2p operation encounters an error.
type LibP2PError struct {
	Err error
}

func (e LibP2PError) Error() string {
	return fmt.Sprintf("libp2p error: %s", e.Err.Error())
}

// PeerStoreError is returned when an loading or saving permanent peer-store encounters an error.
type PeerStoreError struct {
	Err error
}

func (e PeerStoreError) Error() string {
	return fmt.Sprintf("libp2p error: %s", e.Err.Error())
}
