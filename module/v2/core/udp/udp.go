package udp

import (
	"Predator/module/v2/core/spy"
	"Predator/module/v2/lib"
	"Predator/pkg/config"
	"net"
	"strconv"
	"time"
)

var (
	timeout time.Duration
	ports   []int
)

func send(addr string) bool {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return false
	}

	conn, err := net.DialTimeout("udp", udpAddr.String(), timeout)
	if err != nil {
		return false
	}
	defer conn.Close()

	packet := []byte{0xFF, 0xFF, 0xFF, 0xFF}
	if len(packet) > 0 {
		_, err = conn.Write(packet)
		if err != nil {
			return false
		}
	}

	_ = conn.SetReadDeadline(time.Now().Add(timeout))
	data := make([]byte, 1)
	_, err = conn.Read(data)
	if err != nil {
		return false
	}
	return true
}

func udpCheck(ip string) bool {
	for _, port := range ports {
		netloc := net.JoinHostPort(ip, strconv.Itoa(port))
		if send(netloc) {
			lib.Log.Debugf("%s open", netloc)
			return true
		}
	}
	return false
}

func Spy() {
	lib.Log.Info("use udp protocol to spy")
	timeout = time.Duration(config.SpyTimeout) * time.Millisecond
	ports = lib.SetPort(config.SpyUdpPort)
	spy.Spy(udpCheck)
}
