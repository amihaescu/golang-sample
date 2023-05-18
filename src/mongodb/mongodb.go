package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"sample-golang-project/config"
	"sample-golang-project/model"
	"time"
)

type DeviceRepoImp struct {
	cfg        *config.Configuration
	client     *mongo.Client
	collection *mongo.Collection
	logger     *zap.SugaredLogger
}

func New(config *config.Configuration, logger *zap.SugaredLogger) (*DeviceRepoImp, error) {
	credentials := options.Credential{
		Username: config.MongoUser,
		Password: config.MongoPass,
	}
	client, err := mongo.NewClient(options.Client().ApplyURI(config.MongoURI).SetAuth(credentials))
	if err != nil {
		return nil, err
	}
	collection := client.Database(config.MongoDB).Collection(config.MongoColl)
	c := &DeviceRepoImp{
		cfg:        config,
		client:     client,
		collection: collection,
		logger:     logger,
	}
	return c, nil
}

func (c *DeviceRepoImp) Start(ctx context.Context) error {
	return c.client.Connect(ctx)
}

func (c *DeviceRepoImp) Stop(ctx context.Context) error {
	return c.client.Disconnect(ctx)
}

func (c *DeviceRepoImp) Save(ctx context.Context, device *model.Device) (interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	result, err := c.collection.InsertOne(ctx, device)
	if err != nil {
		return nil, err
	}

	c.logger.Info("saved device with name ", device.Name)
	return map[string]interface{}{"id": result.InsertedID}, nil
}
