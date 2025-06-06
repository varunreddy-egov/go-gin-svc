// Package messaging provides interfaces and implementations for message consuming.
// It abstracts the underlying message broker technology (e.g., Kafka, RabbitMQ),
// allowing the application to interact with different messaging systems through a
// common interface.
package messaging

import (
	"clamav-wrapper/models" // Import the actual FileEvent model
)

// MessageConsumer defines the interface for a message consumer.
// It provides a way to start consuming messages and to gracefully close the consumer.
type MessageConsumer interface {
	// StartConsumer begins listening for messages from the configured message broker.
	// It takes a handler function as an argument, which is called for each message received.
	// The method should block until an unrecoverable error occurs or the consumer is closed.
	StartConsumer(handler func(event models.FileEvent) error) error

	// Close gracefully shuts down the message consumer, releasing any resources.
	// It should ensure that any in-flight message processing is completed if possible.
	Close() error
}
