package config

import (
	"os"
	"strings"
)

type Configuration struct {
	KafkaBrokers  []string
	KafkaTopic    string
	KafkaCapacity int
	MongoURI      string
	MongoDB       string
	MongoColl     string
	MongoUser     string
	MongoPass     string
	ServerPort    string
}

func NewConfig() *Configuration {

	return &Configuration{
		KafkaBrokers:  strings.Split(getEnv("SAMPLE_KAFKA_BROKERS", "localhost:9092"), ","),
		KafkaTopic:    getEnv("SAMPLE_KAFKA_TOPIC", "messages"),
		KafkaCapacity: 10,
		MongoURI:      getEnv("SAMPLE_MONGODB_URL", "mongodb://localhost:27017"),
		MongoDB:       getEnv("SAMPLE_MONGODB_NAME", "message_store"),
		MongoUser:     getEnv("SAMPLE_MONGODB_USER", "sample"),
		MongoPass:     getEnv("SAMPLE_MONGODB_USER", "sample"),
		MongoColl:     getEnv("SAMPLE_MONGODB_COL", "message"),
		ServerPort:    getEnv("SAMPLE_SERVER_PORT", ":8080"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}
