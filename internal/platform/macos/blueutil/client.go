package blueutil

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"time"

	"github.com/fumihumi/bt-manage/internal/core"
)

type Client struct {
	Exec    ExecPort
	Bin     string
	Verbose bool
	Logger  io.Writer
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

func (c Client) logf(format string, args ...any) {
	if !c.Verbose {
		return
	}
	w := c.Logger
	if w == nil {
		return
	}
	fmt.Fprintf(w, format, args...)
}

func (c Client) List(ctx context.Context) ([]core.Device, error) {
	if _, err := lookPath(c.bin()); err != nil {
		return nil, core.ErrDependencyMissing{Dependency: c.bin()}
	}
	start := time.Now()
	c.logf("blueutil: start=%s %s --paired --format json\n", start.Format("15:04:05.000"), c.bin())
	stdout, _, err := c.execPort().Run(ctx, c.bin(), "--paired", "--format", "json")
	c.logf("blueutil: done  start=%s elapsed=%s\n", start.Format("15:04:05.000"), time.Since(start).Truncate(time.Millisecond))
	if err != nil {
		return nil, c.mapExecErr(err)
	}
	return parseDeviceListJSON(stdout)
}

func (c Client) Connect(ctx context.Context, address string) error {
	if _, err := lookPath(c.bin()); err != nil {
		return core.ErrDependencyMissing{Dependency: c.bin()}
	}
	addr := denormalizeAddress(address)
	start := time.Now()
	c.logf("blueutil: start=%s %s --connect %s\n", start.Format("15:04:05.000"), c.bin(), addr)
	_, _, err := c.execPort().Run(ctx, c.bin(), "--connect", addr)
	c.logf("blueutil: done  start=%s --connect %s elapsed=%s\n", start.Format("15:04:05.000"), addr, time.Since(start).Truncate(time.Millisecond))
	if err != nil {
		return c.mapExecErr(err)
	}
	return nil
}

func (c Client) Disconnect(ctx context.Context, address string) error {
	if _, err := lookPath(c.bin()); err != nil {
		return core.ErrDependencyMissing{Dependency: c.bin()}
	}
	addr := denormalizeAddress(address)
	start := time.Now()
	c.logf("blueutil: start=%s %s --disconnect %s\n", start.Format("15:04:05.000"), c.bin(), addr)
	_, _, err := c.execPort().Run(ctx, c.bin(), "--disconnect", addr)
	c.logf("blueutil: done  start=%s --disconnect %s elapsed=%s\n", start.Format("15:04:05.000"), addr, time.Since(start).Truncate(time.Millisecond))
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
