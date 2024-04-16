package xhttp

import (
	"net"
	"time"
)

func WrapperTCP(network, address string, forward *net.Dialer) (net.Conn, error) {
	var connect net.Conn
	var err error
	connect, err = forward.Dial(network, address)
	if err != nil {
		return nil, err
	}
	return connect, nil

}

func WrapperTcpWithTimeout(network, address string, timeout time.Duration) (net.Conn, error) {
	dialer := &net.Dialer{Timeout: timeout}
	return WrapperTCP(network, address, dialer)
}
