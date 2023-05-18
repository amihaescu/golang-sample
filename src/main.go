package main

import (
	"context"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"sample-golang-project/api"
	"sample-golang-project/config"
	"sample-golang-project/kafka"
	"sample-golang-project/model"
	"sample-golang-project/mongodb"
	"sample-golang-project/server"
)

func main() {
	ctx := context.Background()
	cfg := config.NewConfig()

	basicLogger, err := zap.NewProduction()
	if err != nil {
		panic("failed to create logger")
	}
	logger := basicLogger.Sugar()

	deviceRepository, err := mongodb.New(cfg, logger)
	if err != nil {
		logger.Error("failed to create mongo repository")
	}

	restToKafka := make(chan *model.Device, cfg.KafkaCapacity)
	kafkaToEs := make(chan *model.Device, cfg.KafkaCapacity)

	publisher := kafka.NewDevicePublisher(ctx, cfg, logger, restToKafka)
	listener := kafka.NewListener(ctx, cfg, logger, kafkaToEs)

	router := mux.NewRouter()
	endpoints := api.NewDeviceController(router, deviceRepository, restToKafka)

	s := server.New(endpoints, deviceRepository, router, publisher, cfg, logger, listener)

	s.Start(ctx)

}
