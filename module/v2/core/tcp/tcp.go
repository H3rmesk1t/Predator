package tcp

import (
	"Predator/module/v2/core/spy"
	"Predator/module/v2/lib"
	"Predator/pkg/config"
	"errors"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

var (
	timeout time.Duration
	ports   []int
)

func send(netloc string) (string, error) {
	conn, err := net.DialTimeout("tcp", netloc, timeout)
	if err != nil {
		return "", errors.New(err.Error() + " STEP1:CONNECT")
	}
	defer conn.Close()

	data := ""
	_, err = io.WriteString(conn, data)
	if err != nil {
		return "", errors.New(err.Error() + " STEP2:WRITE")
	}

	_ = conn.SetReadDeadline(time.Now().Add(timeout))
	size := 0
	buf := make([]byte, size)
	length, err := conn.Read(buf)
	if err != nil && err.Error() != "EOF" {
		return "", errors.New(err.Error() + " STEP3:READ")
	}
	if length == 0 {
		return "", errors.New("STEP3:response is empty")
	}
	return string(buf[:length]), nil
}

func isOpen(netloc string) bool {
	result, err := send(netloc)
	if err == nil {
		return true
	}
	if len(result) > 0 {
		return true
	}
	if strings.Contains(err.Error(), "STEP1") {
		return false
	} else {
		return true
	}
}

func tcpCheck(ip string) bool {
	for _, port := range ports {
		netloc := net.JoinHostPort(ip, strconv.Itoa(port))
		if isOpen(netloc) {
			lib.Log.Debugf("%s open", netloc)
			return true
		}
	}
	return false
}

func Spy() {
	lib.Log.Info("use tcp protocol to spy")
	timeout = time.Duration(config.SpyTimeout) * time.Millisecond
	ports = lib.SetPort(config.SpyTcpPort)
	spy.Spy(tcpCheck)
}
