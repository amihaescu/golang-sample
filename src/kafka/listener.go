package kafka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"sample-golang-project/config"
	"sample-golang-project/model"
	"sync"
)

type DeviceListener struct {
	cfg    *config.Configuration
	reader *kafka.Reader
	logger *zap.SugaredLogger
	output chan *model.Device
	ctx    context.Context
	wg     *sync.WaitGroup
}

func NewListener(ctx context.Context, cfg *config.Configuration, logger *zap.SugaredLogger, output chan *model.Device) *DeviceListener {
	return &DeviceListener{
		logger: logger,
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: cfg.KafkaBrokers,
			Topic:   cfg.KafkaTopic,
			GroupID: "consumer-group-id",
		}),
		cfg:    cfg,
		ctx:    ctx,
		output: output,
		wg:     &sync.WaitGroup{},
	}
}

func (c *DeviceListener) StartListen() error {
	for i := 0; i < c.cfg.KafkaCapacity; i++ {
		c.wg.Add(1)
		go func(workerID int) {
			for {
				m, err := c.reader.ReadMessage(c.ctx)
				if err != nil {
					c.logger.Errorf("failed to read mssage from kafka topic %s", err)
					continue
				}
				var device = &model.Device{}
				if err := json.Unmarshal(m.Value, device); err != nil {
					c.logger.Errorf("failed to deserialize message %v", err)
					continue
				}
				c.logger.Infof("got device from kafka topic %v", device)
				c.output <- device
			}
		}(i)

	}
	return nil
}

func (c *DeviceListener) Stop() error {
	return c.reader.Close()
}
