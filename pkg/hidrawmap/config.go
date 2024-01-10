package hidrawmap

import (
	"bytes"

	"gopkg.in/yaml.v2"
)

type Config struct {
	HIDRAWPath  string             `yaml:"hidraw_path"`
	Assignments map[string]KeyCode `yaml:"assignments"`
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
