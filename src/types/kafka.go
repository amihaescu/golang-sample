package types

import (
	"context"
)

type DevicePublisher interface {
	StartPublish(ctx context.Context) error
}
