package redis

import (
	"Predator/pkg/config"
	"Predator/pkg/utils"
	"Predator/pkg/xhttp"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

var (
	dbfilename string
	dir        string
)

func RedisScan(info *config.HostInfo) (tmperr error) {
	flag, err := redisUnauth(info)
	if flag == true && err == nil {
		return err
	}
	if config.IsBrutePass {
		return
	}
	for _, pass := range config.DefaultPasswords {
		pass = strings.Replace(pass, "{user}", "redis", -1)
		flag, err := redisConn(info, pass)
		if flag == true && err == nil {
			return err
		}
	}
	return tmperr
}

func redisConn(info *config.HostInfo, pass string) (flag bool, err error) {
	flag = false
	realhost := fmt.Sprintf("%s:%v", info.Host, info.Ports)
	conn, err := xhttp.WrapperTcpWithTimeout("tcp", realhost, time.Duration(config.Timeout)*time.Second)
	if err != nil {
		return flag, err
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)
	err = conn.SetReadDeadline(time.Now().Add(time.Duration(config.Timeout) * time.Second))
	if err != nil {
		return flag, err
	}
	_, err = conn.Write([]byte(fmt.Sprintf("auth %s\r\n", pass)))
	if err != nil {
		return flag, err
	}
	reply, err := readreply(conn)
	if err != nil {
		return flag, err
	}
	if strings.Contains(reply, "+OK") {
		flag = true
		dbfilename, dir, err = getconfig(conn)
		if err != nil {
			result := fmt.Sprintf("[!] Redis %s %s", realhost, pass)
			utils.LogSuccess(result)
			return flag, err
		} else {
			result := fmt.Sprintf("[!] Redis %s %s file:%s/%s", realhost, pass, dir, dbfilename)
			utils.LogSuccess(result)
		}
	}
	return flag, err
}

func redisUnauth(info *config.HostInfo) (flag bool, err error) {
	flag = false
	realhost := fmt.Sprintf("%s:%v", info.Host, info.Ports)
	conn, err := xhttp.WrapperTcpWithTimeout("tcp", realhost, time.Duration(config.Timeout)*time.Second)
	if err != nil {
		return flag, err
	}
	defer conn.Close()
	err = conn.SetReadDeadline(time.Now().Add(time.Duration(config.Timeout) * time.Second))
	if err != nil {
		return flag, err
	}
	_, err = conn.Write([]byte("info\r\n"))
	if err != nil {
		return flag, err
	}
	reply, err := readreply(conn)
	if err != nil {
		return flag, err
	}
	if strings.Contains(reply, "redis_version") {
		flag = true
		dbfilename, dir, err = getconfig(conn)
		if err != nil {
			result := fmt.Sprintf("[!] Redis %s unauthorized", realhost)
			utils.LogSuccess(result)
			return flag, err
		} else {
			result := fmt.Sprintf("[!] Redis %s unauthorized file:%s/%s", realhost, dir, dbfilename)
			utils.LogSuccess(result)
		}
	}
	return flag, err
}

func readreply(conn net.Conn) (string, error) {
	conn.SetReadDeadline(time.Now().Add(time.Second))
	bytes, err := io.ReadAll(conn)
	if len(bytes) > 0 {
		err = nil
	}
	return string(bytes), err
}

func getconfig(conn net.Conn) (dbfilename string, dir string, err error) {
	_, err = conn.Write([]byte("CONFIG GET dbfilename\r\n"))
	if err != nil {
		return
	}
	text, err := readreply(conn)
	if err != nil {
		return
	}
	text1 := strings.Split(text, "\r\n")
	if len(text1) > 2 {
		dbfilename = text1[len(text1)-2]
	} else {
		dbfilename = text1[0]
	}
	_, err = conn.Write([]byte("CONFIG GET dir\r\n"))
	if err != nil {
		return
	}
	text, err = readreply(conn)
	if err != nil {
		return
	}
	text1 = strings.Split(text, "\r\n")
	if len(text1) > 2 {
		dir = text1[len(text1)-2]
	} else {
		dir = text1[0]
	}
	return
}
