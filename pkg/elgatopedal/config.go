package elgatopedal

import (
	"bytes"

	"gopkg.in/yaml.v2"
)

type PedalKeyCode struct {
	Press  *int `yaml:"press"`
	OnUp   *int `yaml:"on_up"`
	OnDown *int `yaml:"on_down"`
}

type Config struct {
	HIDRAWPath   string         `yaml:"hidraw_path"`
	PedalKeyCode []PedalKeyCode `yaml:"pedal_keycode"`
}

func (cfg *Config) Bytes() []byte {
	var buf bytes.Buffer
	err := yaml.NewEncoder(&buf).Encode(cfg)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func (cfg *Config) String() string {
	return string(cfg.Bytes())
}
