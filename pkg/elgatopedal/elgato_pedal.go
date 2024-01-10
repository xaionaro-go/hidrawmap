package elgatopedal

import (
	"context"
	"fmt"

	"github.com/bendahl/uinput"
	"github.com/facebookincubator/go-belt/tool/logger"
	"github.com/xaionaro-go/hidrawmap/pkg/hidraw"
)

type ElgatoPedal struct {
	Config      Config
	PedalIsDown []bool
}

func New(cfg Config) (*ElgatoPedal, error) {
	if len(cfg.PedalKeyCode) > 4 {
		return nil, fmt.Errorf("too many pedals, supporting 4 pedals maximum, but received %d", len(cfg.PedalKeyCode))
	}
	return &ElgatoPedal{
		Config:      cfg,
		PedalIsDown: make([]bool, len(cfg.PedalKeyCode)),
	}, nil
}

type KeyCode int

func (h *ElgatoPedal) Serve(ctx context.Context) error {
	hidRaw, err := hidraw.New(h.Config.HIDRAWPath)
	if err != nil {
		return fmt.Errorf("unable to initialize hidraw: %w", err)
	}
	defer hidRaw.Close()

	const uinputPath = "/dev/uinput"
	kb, err := uinput.CreateKeyboard(uinputPath, []byte("hidrawmap"))
	if err != nil {
		return fmt.Errorf("cannot initialize virtual keyboard using uinput '%s': %w", uinputPath, err)
	}

	err = hidRaw.Serve(ctx, func(hidEvent hidraw.HIDEvent) error {
		pedalIsDown := make([]bool, len(h.PedalIsDown))

		// 01000300LLCCRR00
		for pedalIdx := 0; pedalIdx < len(h.Config.PedalKeyCode); pedalIdx++ {
			pedalIsDown[pedalIdx] = hidEvent[4+pedalIdx] != 0

			if pedalIsDown[pedalIdx] == h.PedalIsDown[pedalIdx] {
				continue
			}

			keyCodes := h.Config.PedalKeyCode[pedalIdx]
			if pedalIsDown[pedalIdx] {
				if keyCodes.Press != nil {
					logger.FromCtx(ctx).Debugf("sending KeyDown for keycode %d", *keyCodes.Press)
					kb.KeyDown(*keyCodes.Press)
				}
				if keyCodes.OnDown != nil {
					logger.FromCtx(ctx).Debugf("sending KeyPress for keycode %d", *keyCodes.OnDown)
					kb.KeyPress(*keyCodes.OnDown)
				}
			} else {
				if keyCodes.Press != nil {
					logger.FromCtx(ctx).Debugf("sending KeyUp for keycode %d", *keyCodes.Press)
					kb.KeyUp(*keyCodes.Press)
				}
				if keyCodes.OnUp != nil {
					logger.FromCtx(ctx).Debugf("sending KeyPress for keycode %d", *keyCodes.OnUp)
					kb.KeyPress(*keyCodes.OnUp)
				}
			}
		}

		h.PedalIsDown = pedalIsDown
		return nil
	})
	if err != nil {
		return fmt.Errorf("unable to serve hidraw reader: %w", err)
	}
	return nil
}
