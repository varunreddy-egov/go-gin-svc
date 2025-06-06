package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	MessageBrokerType        string // Added
	KafkaBroker              string
	KafkaTopic               string
	KafkaConsumerGroupID     string
	ClamAVHost               string
	ClamAVPort               int
	ClamAVDialTimeoutSeconds int
	ClamAVChunkSizeKB        int
	ClamAVMaxFileSizeMB      int
	MinioEndpoint            string
	MinioAccessKey           string
	MinioSecretKey           string
	StagingBucket            string
	CleanBucket              string
	QuarantineBucket         string
	UseSSL                   bool
)

func Init() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or error loading it, relying on system env vars.")
	}

	MessageBrokerType = getEnv("MESSAGE_BROKER_TYPE", "kafka") // Added
	KafkaBroker = getEnv("KAFKA_BROKER", "localhost:9092")
	KafkaTopic = getEnv("KAFKA_TOPIC", "minio-events")
	KafkaConsumerGroupID = getEnv("KAFKA_CONSUMER_GROUP_ID", "filestore-antivirus-group")
	ClamAVHost = getEnv("CLAMAV_HOST", "localhost")
	ClamAVPort = getEnvAsInt("CLAMAV_PORT", 3310)
	ClamAVDialTimeoutSeconds = getEnvAsInt("CLAMAV_DIAL_TIMEOUT_SECONDS", 10)
	ClamAVChunkSizeKB = getEnvAsInt("CLAMAV_CHUNK_SIZE_KB", 32)
	ClamAVMaxFileSizeMB = getEnvAsInt("CLAMAV_MAX_FILE_SIZE_MB", 1024*1024)
	MinioEndpoint = getEnv("MINIO_ENDPOINT", "localhost:9000")
	MinioAccessKey = getEnv("MINIO_ACCESS_KEY", "minioadmin")
	MinioSecretKey = getEnv("MINIO_SECRET_KEY", "minioadmin")
	StagingBucket = getEnv("STAGING_BUCKET", "staging")
	CleanBucket = getEnv("CLEAN_BUCKET", "clean")
	QuarantineBucket = getEnv("QUARANTINE_BUCKET", "quarantine")
	UseSSL = getEnvAsBool("USE_SSL", false)
}

func getEnv(key string, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvAsInt(name string, defaultVal int) int {
	valStr := os.Getenv(name)
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return defaultVal
}

func getEnvAsBool(key string, defaultVal bool) bool {
	valStr := os.Getenv(key)
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}
	return defaultVal
}
