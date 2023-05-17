package types

import (
	"context"
	"sample-golang-project/model"
)

type ControllerPublisher interface {
	Listen(ctx context.Context)
	GetChannel() chan *model.Controller
}
