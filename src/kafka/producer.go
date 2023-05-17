package kafka

import (
	"context"
	"log"
	"sample-golang-project/config"
	"time"

	"github.com/segmentio/kafka-go"
)

func PublishMessage(message []byte) {
	cfg := config.NewConfig()

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  cfg.KafkaBrokers,
		Topic:    cfg.KafkaTopic,
		Balancer: &kafka.LeastBytes{},
	})

	defer writer.Close()

	err := writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(time.Now().Format(time.RFC3339)),
			Value: message,
		},
	)
	if err != nil {
		log.Printf("Failed to publish message: %s", err)
		return
	}
	log.Println("Message published to Kafka topic.")
}
