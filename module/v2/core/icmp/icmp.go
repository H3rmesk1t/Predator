package icmp

import (
	"Predator/module/v2/core/spy"
	"Predator/module/v2/lib"
	"Predator/pkg/config"
	"github.com/go-ping/ping"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

var (
	times   int
	timeout time.Duration
	rg      = runtime.GOOS
)

func checkPermission() {
	if rg == "linux" {
		cmd := exec.Command("cat", "/proc/sys/net/ipv4/ping_group_range")
		buf, _ := cmd.Output()
		str := string(buf)
		if !strings.Contains(str, "2147483647") {
			lib.Log.Error("you must manually execute the command to grant the right to send icmp package")
			lib.Log.Error("try sudo sysctl -w net.ipv4.ping_group_range=\"0 2147483647\", or you can try to use the PingSpy module")
			os.Exit(1)
		}
	}
}

func icmpCheck(ip string) bool {
	pinger, err := ping.NewPinger(ip)
	if err != nil {
		lib.Log.Error(err)
		return false
	}
	if rg == "windows" {
		pinger.SetPrivileged(true)
	}
	pinger.Count = times
	pinger.Timeout = timeout
	err = pinger.Run()
	if err != nil {
		lib.Log.Error(err)
	}
	stats := pinger.Statistics()
	if stats.PacketsRecv > 0 {
		return true
	}
	return false
}

func Spy() {
	checkPermission()
	times = config.SpyTimes
	timeout = time.Duration(config.SpyTimeout) * time.Millisecond
	spy.Spy(icmpCheck)
}
