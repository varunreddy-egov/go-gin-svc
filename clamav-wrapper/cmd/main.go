package main

import (
	"log"
	"net/url"

	"clamav-wrapper/clamav"
	"clamav-wrapper/config"
	"clamav-wrapper/messaging"
	"clamav-wrapper/minio"
	"clamav-wrapper/models"
)

// processFileEvent handles the processing of a single file event.
func processFileEvent(event models.FileEvent) error {
	if len(event.Records) == 0 {
		log.Printf("No records in event")
		// Consider returning an error or just logging
		return nil // Or an error indicating no records
	}

	// Assuming we only process the first record, as in the original logic
	record := event.Records[0]
	bucketName := record.S3.Bucket.Name
	objectKey, err := url.QueryUnescape(record.S3.Object.Key)
	if err != nil {
		log.Printf("Invalid object key: %v", err)
		return err // Return error to indicate failure
	}

	log.Printf("Processing file: %s from bucket: %s", objectKey, bucketName)

	file, size, err := minio.GetFileStreamWithSize(bucketName, objectKey)
	if err != nil {
		log.Printf("Failed to get file from MinIO: %v", err)
		return err
	}
	defer file.Close()

	isClean, err := clamav.Scan(file, size)
	if err != nil {
		log.Printf("ClamAV scan error: %v", err)
		return err
	}

	targetBucket := config.CleanBucket
	if !isClean {
		targetBucket = config.QuarantineBucket
		log.Printf("File %s in bucket %s is infected. Moving to quarantine.", objectKey, bucketName)
	} else {
		log.Printf("File %s in bucket %s is clean. Moving to clean bucket.", objectKey, bucketName)
	}

	if err := minio.CopyObject(bucketName, targetBucket, objectKey); err != nil {
		log.Printf("Failed to move file to %s bucket: %v", targetBucket, err)
		return err
	}

	if err := minio.DeleteObject(bucketName, objectKey); err != nil {
		log.Printf("Failed to delete original file %s from bucket %s: %v", objectKey, bucketName, err)
		return err
	}

	log.Printf("File %s processed and moved to %s bucket successfully.", objectKey, targetBucket)
	return nil
}

func main() {
	config.Init()
	minio.Init() // MinIO client needs to be initialized

	var consumer messaging.MessageConsumer
	var err error

	log.Printf("Using message broker type: %s", config.MessageBrokerType)

	switch config.MessageBrokerType {
	case "kafka":
		// Create the Kafka consumer
		consumer = messaging.NewKafkaConsumer()
	default:
		log.Fatalf("Unsupported message broker type: %s", config.MessageBrokerType)
		// No need to return here as Fatalf will exit
	}

	// Defer the closing of the consumer right after it's successfully created.
	if consumer != nil {
		defer func() {
			log.Println("Closing consumer...")
			if err := consumer.Close(); err != nil {
				log.Printf("Error closing consumer: %v", err)
			}
		}()
	}

	log.Println("Starting consumer...")
	// Start the consumer with the defined handler function
	// This will block and continuously process messages.
	// If StartConsumer returns an error, the deferred Close will still be called.
	err = consumer.StartConsumer(processFileEvent)
	if err != nil {
		log.Fatalf("Consumer error: %v", err)
	}
}
