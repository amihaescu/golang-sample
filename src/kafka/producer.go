package kafka

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"sample-golang-project/config"
	"sample-golang-project/model"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

type DevicePublisher struct {
	input  chan *model.Device
	cfg    *config.Configuration
	writer *kafka.Conn
	logger *zap.SugaredLogger
	wg     *sync.WaitGroup
}

func NewDevicePublisher(ctx context.Context, cfg *config.Configuration, logger *zap.SugaredLogger, input chan *model.Device) *DevicePublisher {
	leader, err := kafka.DialLeader(ctx, "tcp", cfg.KafkaBrokers[0], cfg.KafkaTopic, 0)
	if err != nil {

	}
	return &DevicePublisher{
		input:  input,
		cfg:    cfg,
		writer: leader,
		logger: logger,
		wg:     &sync.WaitGroup{},
	}
}

func (p *DevicePublisher) StartPublish(ctx context.Context) error {
	for controller := range p.input {
		p.send(controller)
	}
	return nil
}

func (p *DevicePublisher) send(controller *model.Device) {
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
	if err != nil {
		p.logger.Errorf("failed to sent message to kafka topic")
	}
	p.logger.Info("Message published to Kafka topic.")
}
