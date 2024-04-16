package v2

import (
	"Predator/module/v2/core/icmp"
	"Predator/module/v2/core/ping"
	"Predator/module/v2/core/tcp"
	"Predator/module/v2/core/udp"
	"Predator/module/v2/lib"
	"Predator/pkg/config"
	"os"
)

func Execute() {
	lib.InitLog()
	lib.RecEnvInfo()

	if _, err := os.Stat(config.SpyOutput); os.IsNotExist(err) {
		file, err := os.Create(config.SpyOutput)
		if err != nil {
			panic(err)
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
			}
		}(file)
	} else if err != nil {
		panic(err)
	}

	switch config.SpyModule {
	case "icmp":
		icmp.Spy()
	case "ping":
		ping.Spy()
	case "tcp":
		tcp.Spy()
	case "udp":
		udp.Spy()
	default:
		lib.Log.Error("spy type does not exist")
		os.Exit(1)
	}
}
