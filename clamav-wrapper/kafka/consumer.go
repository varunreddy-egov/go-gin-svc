package kafka

import (
	"context"
	"encoding/json"
	"log"
	"net/url"

	"github.com/segmentio/kafka-go"

	"clamav-wrapper/clamav"
	"clamav-wrapper/config"
	"clamav-wrapper/minio"
	"clamav-wrapper/models"
)

func StartConsumer() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{config.KafkaBroker},
		Topic:   config.KafkaTopic,
		GroupID: config.KafkaConsumerGroupID,
	})

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Kafka read error: %v", err)
			continue
		}

		var event models.FileEvent
		if err := json.Unmarshal(m.Value, &event); err != nil {
			log.Printf("Invalid event: %v", err)
			continue
		}

		if len(event.Records) == 0 {
			log.Printf("No records in event")
			continue
		}

		bucketName := event.Records[0].S3.Bucket.Name
		objectKey, err := url.QueryUnescape(event.Records[0].S3.Object.Key)
		if err != nil {
			log.Fatalf("Invalid object key: %v", err)
		}

		log.Printf("Processing file: %s from bucket: %s", objectKey, bucketName)

		file, size, err := minio.GetFileStreamWithSize(bucketName, objectKey)

		if err != nil {
			log.Printf("Failed to get file from MinIO: %v", err)
			continue
		}

		isClean, err := clamav.Scan(file, size)
		file.Close()
		if err != nil {
			log.Printf("ClamAV scan error: %v", err)
			continue
		}

		targetBucket := config.CleanBucket
		if !isClean {
			targetBucket = config.QuarantineBucket
		}

		if err := minio.CopyObject(bucketName, targetBucket, objectKey); err != nil {
			log.Printf("Failed to move file: %v", err)
			continue
		}

		_ = minio.DeleteObject(bucketName, objectKey)
		log.Printf("File moved to %s bucket", targetBucket)
	}
}
