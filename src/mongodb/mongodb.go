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

type ControllerRepoImp struct {
	cfg        *config.Configuration
	client     *mongo.Client
	collection *mongo.Collection
	logger     *zap.SugaredLogger
}

func New(config *config.Configuration, logger *zap.SugaredLogger) (*ControllerRepoImp, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return nil, err
	}
	collection := client.Database(config.MongoDB).Collection(config.MongoColl)
	c := &ControllerRepoImp{
		cfg:        config,
		client:     client,
		collection: collection,
		logger:     logger,
	}
	return c, nil
}

func (c *ControllerRepoImp) Start(ctx context.Context) error {
	return c.client.Connect(ctx)
}

func (c *ControllerRepoImp) Stop(ctx context.Context) error {
	return c.client.Disconnect(ctx)
}

func (c *ControllerRepoImp) Save(ctx context.Context, controller *model.Controller) (interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	result, err := c.collection.InsertOne(ctx, controller)
	if err != nil {
		return nil, err
	}

	c.logger.Info("saved controller with name ", controller.Name)
	return map[string]interface{}{"id": result.InsertedID}, nil
}
