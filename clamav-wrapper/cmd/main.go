package main

import (
	"clamav-wrapper/config"
	"clamav-wrapper/kafka"
	"clamav-wrapper/minio"
)

func main() {
	config.Init()
	minio.Init()
	kafka.StartConsumer()
}
