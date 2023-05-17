package config

import (
	"os"
	"strings"
)

type Configuration struct {
	KafkaBrokers []string
	KafkaTopic   string
	MongoURI     string
	MongoDB      string
	MongoColl    string
	ServerPort   string
}

func NewConfig() *Configuration {

	return &Configuration{
		KafkaBrokers: strings.Split(getEnv("SAMPLE_KAFKA_BROKERS", "localhost:9092"), ","),
		KafkaTopic:   getEnv("SAMPLE_KAFKA_TOPIC", "messages"),
		MongoURI:     getEnv("SAMPLE_MONGODB_URL", "mongodb://localhost:27017"),
		MongoDB:      getEnv("SAMPLE_MONGODB_NAME", "message_store"),
		MongoColl:    getEnv("SAMPLE_MONGODB_COL", "message"),
		ServerPort:   getEnv("SAMPLE_SERVER_PORT", ":8080"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}
