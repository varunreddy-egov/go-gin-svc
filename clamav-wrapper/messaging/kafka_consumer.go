package messaging

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"

	"clamav-wrapper/config"
	"clamav-wrapper/models"
)

// KafkaConsumer implements the MessageConsumer interface for Apache Kafka.
// It handles the connection to Kafka, message consumption, and deserialization.
type KafkaConsumer struct {
	// Reader is the underlying Kafka reader from the segmentio/kafka-go library.
	Reader *kafka.Reader
}

// NewKafkaConsumer creates and configures a new KafkaConsumer.
// It initializes a Kafka reader based on global configuration settings
// (KafkaBroker, KafkaTopic, KafkaConsumerGroupID).
func NewKafkaConsumer() *KafkaConsumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{config.KafkaBroker}, // Assumes config.KafkaBroker is properly initialized
		Topic:   config.KafkaTopic,           // Assumes config.KafkaTopic is properly initialized
		GroupID: config.KafkaConsumerGroupID, // Assumes config.KafkaConsumerGroupID is properly initialized
	})
	return &KafkaConsumer{Reader: r}
}

// StartConsumer begins consuming messages from the Kafka topic.
// It continuously reads messages, deserializes them into models.FileEvent,
// and then passes them to the provided handler function.
// This method will block until a read error occurs or the context is cancelled
// (which can be triggered by closing the reader via the Close method).
func (kc *KafkaConsumer) StartConsumer(handler func(event models.FileEvent) error) error {
	for {
		m, err := kc.Reader.ReadMessage(context.Background())
		if err != nil {
			// If the reader is closed, ReadMessage will return an error.
			// This is the expected way to stop the consumer.
			log.Printf("Kafka read error: %v. This might indicate the consumer is closing.", err)
			return err // Return error to signal consumer stop or failure
		}

		var event models.FileEvent
		if err := json.Unmarshal(m.Value, &event); err != nil {
			log.Printf("Invalid event message: %v. Skipping message.", err)
			continue // Skip malformed messages
		}

		// Call the handler with the received event
		if err := handler(event); err != nil {
			log.Printf("Error processing event: %v. Continuing to next message.", err)
			// Decide if we should continue or stop based on the handler's error.
			// For now, logging and continuing.
		}
	}
}

// Close gracefully shuts down the Kafka consumer by closing the underlying Kafka reader.
// This will cause the StartConsumer loop to exit.
func (kc *KafkaConsumer) Close() error {
	if kc.Reader != nil {
		log.Println("Closing Kafka consumer reader.")
		return kc.Reader.Close()
	}
	return nil
}
