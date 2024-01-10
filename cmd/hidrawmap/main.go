package main

import (
	"github.com/xaionaro-go/hidrawmap/cmd/hidrawmap/command"
)

func main() {
	if err := command.Root.Execute(); err != nil {
		panic(err)
	}
}
