package types

import (
	"context"
	"sample-golang-project/model"
)

type ControllerRepository interface {
	Save(ctx context.Context, controller *model.Controller) (interface{}, error)
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}
