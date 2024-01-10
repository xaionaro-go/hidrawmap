package command

import (
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/facebookincubator/go-belt"
	"github.com/facebookincubator/go-belt/tool/logger"
	"github.com/facebookincubator/go-belt/tool/logger/implementation/zap"
	"github.com/spf13/cobra"
	"github.com/xaionaro-go/hidrawmap/pkg/elgatopedal"
	"go.uber.org/config"
)

type Config = elgatopedal.Config

func ptr[T any](v T) *T {
	return &v
}

func ExampleConfig() Config {
	return Config{
		HIDRAWPath: "/dev/hidraw3",
		PedalKeyCode: []elgatopedal.PedalKeyCode{
			{Press: ptr(241)},
			{OnDown: ptr(248), OnUp: ptr(240)},
			{Press: ptr(242)},
		},
	}
}

var logLevel = logger.LevelWarning
var Root = &cobra.Command{
	Use:  "elgatopedal",
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := logger.CtxWithLogger(cmd.Context(), zap.Default().WithLevel(logLevel))
		defer belt.Flush(ctx)

		cfgPath, err := cmd.Flags().GetString("config")
		if err != nil {
			panic(err)
		}
		cfgPath = expandPath(cfgPath)
		var cfg Config
		cfgFile, err := os.Open(cfgPath)
		if err == nil {
			provider, err := config.NewYAML(config.Source(cfgFile))
			if err != nil {
				panic(err)
			}
			provider.Get(config.Root).Populate(&cfg)
			cfgFile.Close()
		} else {
			cfg = ExampleConfig()
			err := ioutil.WriteFile(cfgPath, cfg.Bytes(), 0644)
			if err != nil {
				panic(err)
			}
		}

		pedal, err := elgatopedal.New(cfg)
		if err != nil {
			panic(err)
		}
		err = pedal.Serve(ctx)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	Root.Flags().String("config", "~/.elgatopedal.yaml", "")
	Root.Flags().Var(&logLevel, "log-level", "")
}

func expandPath(path string) string {
	switch {
	case path == "~":
		return homeDir()
	case strings.HasPrefix(path, "~/"):
		return filepath.Join(homeDir(), path[2:])
	}
	return path
}

func homeDir() string {
	usr, _ := user.Current()
	return usr.HomeDir
}
