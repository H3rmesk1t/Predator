package main

import (
	v1 "Predator/module/v1"
	v2 "Predator/module/v2"
	"Predator/pkg/common"
	"Predator/pkg/config"
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	var Info config.HostInfo
	common.Help(&Info)

	if !config.ModuleFlag {
		v2.Execute()
	} else {
		common.Process(&Info)
		v1.Start(Info)
		fmt.Print(time.Since(start))
	}
}
