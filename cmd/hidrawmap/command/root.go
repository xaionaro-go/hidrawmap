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
	"github.com/xaionaro-go/hidrawmap/pkg/hidrawmap"
	"go.uber.org/config"
)

type Config = hidrawmap.Config

func ExampleConfig() Config {
	return Config{
		HIDRAWPath: "/dev/hidraw3",
		Assignments: map[string]hidrawmap.KeyCode{
			"0100030001000000": 241,
			"0100030000010000": 248,
			"0100030000000100": 242,
		},
	}
}

var logLevel = logger.LevelWarning
var Root = &cobra.Command{
	Use:  "hidrawmap",
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

		err = hidrawmap.New(cfg).Serve(ctx)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	Root.Flags().String("config", "~/.hidrawmap.yaml", "")
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
