package blueutil

import (
	"context"

	"github.com/fumihumi/bt-manage/internal/core"
)

// Client is a placeholder. Phase 4 will implement real blueutil integration.
// It currently returns an empty device list.
type Client struct{}

func (c Client) List(ctx context.Context) ([]core.Device, error) {
	return []core.Device{}, nil
}

func (c Client) Connect(ctx context.Context, address string) error {
	return nil
}

func (c Client) Disconnect(ctx context.Context, address string) error {
	return nil
}
