package main

import (
	"github.com/xaionaro-go/hidrawmap/cmd/elgatopedal/command"
)

func main() {
	if err := command.Root.Execute(); err != nil {
		panic(err)
	}
}
