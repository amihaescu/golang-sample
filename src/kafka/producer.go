package kafka

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"sample-golang-project/config"
	"sample-golang-project/model"
	"time"

	"github.com/segmentio/kafka-go"
)

type ControllerPublisher struct {
	input  chan *model.Controller
	cfg    *config.Configuration
	writer *kafka.Conn
	logger *zap.SugaredLogger
}

func NewControllerPublisher(ctx context.Context, cfg *config.Configuration, logger *zap.SugaredLogger) *ControllerPublisher {
	leader, err := kafka.DialLeader(ctx, "tcp", cfg.KafkaBrokers[0], cfg.KafkaTopic, 0)
	if err != nil {

	}
	return &ControllerPublisher{
		input:  make(chan *model.Controller, cfg.KafkaCapacity),
		cfg:    cfg,
		writer: leader,
		logger: logger,
	}
}

func (p *ControllerPublisher) Listen(ctx context.Context) {
	go func() {
		for controller := range p.input {
			p.send(controller)
		}
	}()
}

func (p *ControllerPublisher) GetChannel() chan *model.Controller {
	return p.input
}

func (p *ControllerPublisher) send(controller *model.Controller) {
	marshal, err := json.Marshal(controller)
	if err != nil {
		p.logger.Error("failed to serialize message ", err)
		return
	}
	_, err = p.writer.WriteMessages(
		kafka.Message{
			Time:  time.Now(),
			Value: marshal,
		},
	)
	p.logger.Info("Message published to Kafka topic.")
}
