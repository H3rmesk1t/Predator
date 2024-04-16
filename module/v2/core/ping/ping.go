package ping

import (
	"Predator/module/v2/core/spy"
	"Predator/module/v2/lib"
	"Predator/pkg/config"
	"strconv"
)

var (
	times   string
	timeout string
)

func pingCheck(ip string) bool {
	if lib.IsPing(ip, times, timeout) {
		return true
	} else {
		return false
	}
}

func Spy() {
	lib.Log.Info("use ping command to spy")
	times = strconv.Itoa(config.SpyTimes)
	timeout = strconv.Itoa(config.SpyTimeout)
	spy.Spy(pingCheck)
}
