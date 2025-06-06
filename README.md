# ClamAV Wrapper Service

This service acts as a wrapper around ClamAV, primarily designed to scan files received from a message queue (e.g., Kafka) that originate from an S3-compatible storage (e.g., MinIO). Infected files are moved to a quarantine bucket, while clean files are moved to a clean bucket.

## Features

*   Consumes file event messages from a configured message broker.
*   Fetches files from S3-compatible storage (MinIO).
*   Scans files using ClamAV.
*   Moves files to appropriate buckets (clean/quarantine) based on scan results.
*   Configurable via environment variables.

## Configuration

The application is configured using environment variables. An `.env` file can also be used for local development.

### General Configuration
*   `MESSAGE_BROKER_TYPE`: Specifies the type of message broker to use.
    *   Currently supported: `kafka` (default)
*   `MINIO_ENDPOINT`: MinIO server endpoint (e.g., `localhost:9000`).
*   `MINIO_ACCESS_KEY`: MinIO access key.
*   `MINIO_SECRET_KEY`: MinIO secret key.
*   `STAGING_BUCKET`: The S3 bucket where new files are initially uploaded and events are triggered from.
*   `CLEAN_BUCKET`: The S3 bucket to move files to if they are scanned and found clean.
*   `QUARANTINE_BUCKET`: The S3 bucket to move files to if they are scanned and found infected.
*   `USE_SSL`: Set to `true` if MinIO connection should use SSL. Defaults to `false`.

### ClamAV Configuration
*   `CLAMAV_HOST`: Hostname for the ClamAV daemon (e.g., `localhost`).
*   `CLAMAV_PORT`: Port number for the ClamAV daemon (e.g., `3310`).
*   `CLAMAV_DIAL_TIMEOUT_SECONDS`: Timeout in seconds for connecting to ClamAV.
*   `CLAMAV_CHUNK_SIZE_KB`: Size of chunks (in KB) for streaming files to ClamAV.
*   `CLAMAV_MAX_FILE_SIZE_MB`: Maximum file size (in MB) to scan.

### Kafka Configuration (if `MESSAGE_BROKER_TYPE=kafka`)
*   `KAFKA_BROKER`: Kafka broker address (e.g., `localhost:9092`).
*   `KAFKA_TOPIC`: Kafka topic to consume messages from (e.g., `minio-events`).
*   `KAFKA_CONSUMER_GROUP_ID`: Kafka consumer group ID (e.g., `clamav-wrapper-group`).