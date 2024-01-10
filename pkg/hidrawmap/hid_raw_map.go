package hidrawmap

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/bendahl/uinput"
	"github.com/facebookincubator/go-belt/tool/logger"
	"github.com/xaionaro-go/hidrawmap/pkg/hidraw"
)

type HIDRawMap struct {
	Config Config
}

func New(cfg Config) *HIDRawMap {
	return &HIDRawMap{
		Config: cfg,
	}
}

type KeyCode int

func (h *HIDRawMap) Serve(ctx context.Context) error {
	hidRaw, err := hidraw.New(h.Config.HIDRAWPath)
	if err != nil {
		return fmt.Errorf("unable to initialize hidraw: %w", err)
	}
	defer hidRaw.Close()

	assignments := map[hidraw.HIDEvent]KeyCode{}
	for hidRawEventHEX, keyCode := range h.Config.Assignments {
		hidRawEventBytes, err := hex.DecodeString(hidRawEventHEX)
		if err != nil {
			return fmt.Errorf("unable to decode hex '%s': %w", hidRawEventHEX, err)
		}

		if len(hidRawEventBytes) != len(hidraw.HIDEvent{}) {
			return fmt.Errorf(
				"invalid length of the HID event '%s', should be %d, but got %d",
				hidRawEventHEX,
				len(hidraw.HIDEvent{}),
				len(hidRawEventBytes),
			)
		}

		var hidRawEvent hidraw.HIDEvent
		copy(hidRawEvent[:], hidRawEventBytes)
		assignments[hidRawEvent] = keyCode
	}

	const uinputPath = "/dev/uinput"
	kb, err := uinput.CreateKeyboard(uinputPath, []byte("hidrawmap"))
	if err != nil {
		return fmt.Errorf("cannot initialize virtual keyboard using uinput '%s': %w", uinputPath, err)
	}

	err = hidRaw.Serve(ctx, func(hidEvent hidraw.HIDEvent) error {
		keyCode, ok := assignments[hidEvent]
		if !ok {
			logger.FromCtx(ctx).Debugf("no keycode assigned to %X, skipping", hidEvent[:])
			return nil
		}
		logger.FromCtx(ctx).Debugf("received keycode %d, simulating keypress", keyCode)
		err := kb.KeyPress(int(keyCode))
		if err != nil {
			return fmt.Errorf("unable to press key with keyCode %d: %w", keyCode, err)
		}
		logger.FromCtx(ctx).Debugf("keypress simulating completed")
		return nil
	})
	if err != nil {
		return fmt.Errorf("unable to serve hidraw reader: %w", err)
	}
	return nil
}
