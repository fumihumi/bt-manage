package blueutil

import (
	"context"
	"errors"
	"fmt"
	"os/exec"

	"github.com/fumihumi/bt-manage/internal/core"
)

type Client struct {
	Exec ExecPort
	Bin  string
}

func (c Client) bin() string {
	if c.Bin != "" {
		return c.Bin
	}
	return "blueutil"
}

func (c Client) execPort() ExecPort {
	if c.Exec != nil {
		return c.Exec
	}
	return OSExec{}
}

func (c Client) List(ctx context.Context) ([]core.Device, error) {
	stdout, _, err := c.execPort().Run(ctx, c.bin(), "--paired", "--format", "json")
	if err != nil {
		return nil, c.mapExecErr(err)
	}
	return parseDeviceListJSON(stdout)
}

func (c Client) Connect(ctx context.Context, address string) error {
	_, _, err := c.execPort().Run(ctx, c.bin(), "--connect", denormalizeAddress(address))
	if err != nil {
		return c.mapExecErr(err)
	}
	return nil
}

func (c Client) Disconnect(ctx context.Context, address string) error {
	_, _, err := c.execPort().Run(ctx, c.bin(), "--disconnect", denormalizeAddress(address))
	if err != nil {
		return c.mapExecErr(err)
	}
	return nil
}

func (c Client) mapExecErr(err error) error {
	var ee *exec.Error
	if errors.As(err, &ee) {
		if errors.Is(ee.Err, exec.ErrNotFound) {
			return core.ErrDependencyMissing{Dependency: c.bin()}
		}
	}
	return fmt.Errorf("blueutil: %w", err)
}

func denormalizeAddress(addr string) string {
	// blueutil は '-' 区切りを受け付ける（':' でも動くが念のため合わせる）
	b := []byte(addr)
	for i := range b {
		if b[i] == ':' {
			b[i] = '-'
		}
	}
	return string(b)
}
