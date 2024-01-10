package hidraw

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/facebookincubator/go-belt/tool/experimental/errmon"
	"github.com/facebookincubator/go-belt/tool/logger"
	"go.uber.org/atomic"
)

type HIDRaw struct {
	serveCount atomic.Int32
	closeOnce  sync.Once
	reader     io.ReadCloser
}

func New(path string) (*HIDRaw, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("unable to open '%s': %w", path, err)
	}

	return NewFromReader(file), nil
}

func NewFromReader(r io.ReadCloser) *HIDRaw {
	return &HIDRaw{
		reader: r,
	}
}

func (h *HIDRaw) Close() error {
	var err error
	h.closeOnce.Do(func() {
		err = h.reader.Close()
	})
	return err
}

func (h *HIDRaw) Serve(ctx context.Context, callback func(HIDEvent) error) error {
	if r := h.serveCount.Add(1); r != 1 {
		return fmt.Errorf("Serve could be used only once, but was already called %d times", r)
	}
	ctx, cancelFn := context.WithCancel(ctx)
	defer cancelFn()
	go func() {
		<-ctx.Done()
		err := h.Close()
		errmon.ObserveErrorCtx(ctx, err)
	}()

	for {
		var hidEvent HIDEvent
		logger.FromCtx(ctx).Debugf("waiting for a HID event")
		n, err := h.reader.Read(hidEvent[:])
		if err != nil {
			return fmt.Errorf("cannot read from hidraw: %w", err)
		}
		if n != len(hidEvent) {
			return fmt.Errorf("the received a HID event is of invalid length: received %d, expected %d: %X", len(hidEvent), n, hidEvent[:n])
		}
		logger.FromCtx(ctx).Debugf("received a HID event %X", hidEvent[:])

		err = callback(hidEvent)
		if err != nil {
			logger.FromCtx(ctx).Errorf("received an error from the callback: %w", err)
			return err
		}
	}
}
