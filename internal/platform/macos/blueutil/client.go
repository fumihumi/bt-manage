package blueutil

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
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
	stdout, stderr, err := c.execPort().Run(ctx, c.bin(), "--paired", "--format", "json")
	c.logf("blueutil: done  start=%s elapsed=%s\n", start.Format("15:04:05.000"), time.Since(start).Truncate(time.Millisecond))
	if err != nil {
		return nil, c.mapExecErrWithStderr(err, stderr)
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
	_, stderr, err := c.execPort().Run(ctx, c.bin(), "--connect", addr)
	c.logf("blueutil: done  start=%s --connect %s elapsed=%s\n", start.Format("15:04:05.000"), addr, time.Since(start).Truncate(time.Millisecond))
	if err != nil {
		return c.mapExecErrWithStderr(err, stderr)
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
	_, stderr, err := c.execPort().Run(ctx, c.bin(), "--disconnect", addr)
	c.logf("blueutil: done  start=%s --disconnect %s elapsed=%s\n", start.Format("15:04:05.000"), addr, time.Since(start).Truncate(time.Millisecond))
	if err != nil {
		return c.mapExecErrWithStderr(err, stderr)
	}
	return nil
}

func (c Client) Pair(ctx context.Context, address string, pin string) error {
	if _, err := lookPath(c.bin()); err != nil {
		return core.ErrDependencyMissing{Dependency: c.bin()}
	}
	addr := denormalizeAddress(address)
	args := []string{"--pair", addr}
	if pin != "" {
		args = append(args, pin)
	}
	start := time.Now()
	c.logf("blueutil: start=%s %s %s\n", start.Format("15:04:05.000"), c.bin(), strings.Join(args, " "))
	_, stderr, err := c.execPort().Run(ctx, c.bin(), args...)
	c.logf("blueutil: done  start=%s --pair %s elapsed=%s\n", start.Format("15:04:05.000"), addr, time.Since(start).Truncate(time.Millisecond))
	if err != nil {
		return c.mapExecErrWithStderr(err, stderr)
	}
	return nil
}

func (c Client) Unpair(ctx context.Context, address string) error {
	if _, err := lookPath(c.bin()); err != nil {
		return core.ErrDependencyMissing{Dependency: c.bin()}
	}
	addr := denormalizeAddress(address)
	start := time.Now()
	c.logf("blueutil: start=%s %s --unpair %s\n", start.Format("15:04:05.000"), c.bin(), addr)
	_, stderr, err := c.execPort().Run(ctx, c.bin(), "--unpair", addr)
	c.logf("blueutil: done  start=%s --unpair %s elapsed=%s\n", start.Format("15:04:05.000"), addr, time.Since(start).Truncate(time.Millisecond))
	if err != nil {
		return c.mapExecErrWithStderr(err, stderr)
	}
	return nil
}

func (c Client) Inquiry(ctx context.Context, durationSeconds int) ([]core.Device, error) {
	if _, err := lookPath(c.bin()); err != nil {
		return nil, core.ErrDependencyMissing{Dependency: c.bin()}
	}
	if durationSeconds <= 0 {
		durationSeconds = 10
	}
	start := time.Now()
	c.logf("blueutil: start=%s %s --inquiry %d --format json\n", start.Format("15:04:05.000"), c.bin(), durationSeconds)
	stdout, stderr, err := c.execPort().Run(ctx, c.bin(), "--inquiry", strconv.Itoa(durationSeconds), "--format", "json")
	c.logf("blueutil: done  start=%s --inquiry %d elapsed=%s\n", start.Format("15:04:05.000"), durationSeconds, time.Since(start).Truncate(time.Millisecond))
	if err != nil {
		return nil, c.mapExecErrWithStderr(err, stderr)
	}
	return parseDeviceListJSON(stdout)
}

func (c Client) WaitConnect(ctx context.Context, address string, timeoutSeconds int) error {
	if _, err := lookPath(c.bin()); err != nil {
		return core.ErrDependencyMissing{Dependency: c.bin()}
	}
	addr := denormalizeAddress(address)
	args := []string{"--wait-connect", addr}
	if timeoutSeconds > 0 {
		args = append(args, strconv.Itoa(timeoutSeconds))
	}
	start := time.Now()
	c.logf("blueutil: start=%s %s %s\n", start.Format("15:04:05.000"), c.bin(), strings.Join(args, " "))
	_, stderr, err := c.execPort().Run(ctx, c.bin(), args...)
	c.logf("blueutil: done  start=%s --wait-connect %s elapsed=%s\n", start.Format("15:04:05.000"), addr, time.Since(start).Truncate(time.Millisecond))
	if err != nil {
		return c.mapExecErrWithStderr(err, stderr)
	}
	return nil
}

func (c Client) IsConnected(ctx context.Context, address string) (bool, error) {
	if _, err := lookPath(c.bin()); err != nil {
		return false, core.ErrDependencyMissing{Dependency: c.bin()}
	}
	addr := denormalizeAddress(address)
	start := time.Now()
	c.logf("blueutil: start=%s %s --is-connected %s\n", start.Format("15:04:05.000"), c.bin(), addr)
	stdout, stderr, err := c.execPort().Run(ctx, c.bin(), "--is-connected", addr)
	c.logf("blueutil: done  start=%s --is-connected %s elapsed=%s\n", start.Format("15:04:05.000"), addr, time.Since(start).Truncate(time.Millisecond))
	if err != nil {
		return false, c.mapExecErrWithStderr(err, stderr)
	}
	v := strings.TrimSpace(string(stdout))
	return v == "1", nil
}

func (c Client) ConnectedDevices(ctx context.Context) ([]core.Device, error) {
	if _, err := lookPath(c.bin()); err != nil {
		return nil, core.ErrDependencyMissing{Dependency: c.bin()}
	}
	start := time.Now()
	c.logf("blueutil: start=%s %s --connected --format json\n", start.Format("15:04:05.000"), c.bin())
	stdout, stderr, err := c.execPort().Run(ctx, c.bin(), "--connected", "--format", "json")
	c.logf("blueutil: done  start=%s --connected elapsed=%s\n", start.Format("15:04:05.000"), time.Since(start).Truncate(time.Millisecond))
	if err != nil {
		return nil, c.mapExecErrWithStderr(err, stderr)
	}
	return parseDeviceListJSON(stdout)
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

func (c Client) mapExecErrWithStderr(err error, stderr []byte) error {
	base := c.mapExecErr(err)
	msg := strings.TrimSpace(string(stderr))
	if msg == "" {
		return base
	}
	return fmt.Errorf("%w: %s", base, msg)
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
