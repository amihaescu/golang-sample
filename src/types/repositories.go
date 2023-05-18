package types

import (
	"context"
	"sample-golang-project/model"
)

type DeviceRepository interface {
	Save(ctx context.Context, device *model.Device) (interface{}, error)
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}
