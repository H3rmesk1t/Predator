package telnet

import (
	"Predator/pkg/config"
	"Predator/pkg/utils"
	"fmt"
	"net"
	"strings"
	"time"
)

func TelnetScan(info *config.HostInfo) (tmperr error) {
	if config.IsBrutePass {
		return
	}

	for _, user := range config.DefaultUsers["telnet"] {
		for _, pass := range config.DefaultPasswords {
			pass = strings.Replace(pass, "{user}", user, -1)
			flag, err := telnetConn(info, user, pass)
			if flag && err == nil {
				return err
			}
		}
	}
	return tmperr
}

func telnetConn(info *config.HostInfo, user string, pass string) (flag bool, err error) {
	flag = false
	Host, Port, Username, Password := info.Host, info.Ports, user, pass
	address := fmt.Sprintf("%v:%v", Host, Port)
	conn, err := net.DialTimeout("tcp", address, time.Duration(config.Timeout)*time.Second)
	if err == nil {
		defer func(conn net.Conn) {
			err := conn.Close()
			if err != nil {
			}
		}(conn)
		if telnetProtocolHandshake(conn, Username, Password) {
			result := fmt.Sprintf("[!] Telnet:%v:%v:%v %v", Host, Port, Username, Password)
			utils.LogSuccess(result)
			flag = true
		}
	}
	return flag, err
}

func telnetProtocolHandshake(conn net.Conn, Username string, Password string) bool {
	var buf [4096]byte
	var n int
	n, err := conn.Read(buf[0:])
	if nil != err {
		return false
	}

	buf[0] = 0xff
	buf[1] = 0xfc
	buf[2] = 0x25
	buf[3] = 0xff
	buf[4] = 0xfe
	buf[5] = 0x01
	n, err = conn.Write(buf[0:6])
	if nil != err {
		return false
	}

	n, err = conn.Read(buf[0:])
	if nil != err {
		return false
	}

	buf[0] = 0xff
	buf[1] = 0xfe
	buf[2] = 0x03
	buf[3] = 0xff
	buf[4] = 0xfc
	buf[5] = 0x27
	n, err = conn.Write(buf[0:6])
	if nil != err {
		return false
	}

	n, err = conn.Read(buf[0:])
	if nil != err {
		return false
	}

	n, err = conn.Write([]byte(Username + "\r\n"))
	if nil != err {
		return false
	}
	time.Sleep(time.Millisecond * 500)

	n, err = conn.Read(buf[0:])
	if nil != err {
		return false
	}

	n, err = conn.Write([]byte(Password + "\r\n"))
	if nil != err {
		return false
	}
	time.Sleep(time.Millisecond * 2000)
	n, err = conn.Read(buf[0:])
	if nil != err {
		return false
	}
	if strings.Contains(string(buf[0:n]), "Login Failed") {
		return false
	}

	buf[0] = 0xff
	buf[1] = 0xfc
	buf[2] = 0x18

	n, err = conn.Write(buf[0:3])
	if nil != err {
		return false
	}
	n, err = conn.Read(buf[0:])
	if nil != err {
		return false
	}
	return true
}
