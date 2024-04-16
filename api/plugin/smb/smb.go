package smb

import (
	"Predator/pkg/config"
	"Predator/pkg/utils"
	"errors"
	"fmt"
	"github.com/stacktitan/smb/smb"
	"strings"
	"time"
)

func SmbScan(info *config.HostInfo) (tmperr error) {
	if config.IsBrutePass {
		return nil
	}
	for _, user := range config.DefaultUsers["smb"] {
		for _, pass := range config.DefaultPasswords {
			pass = strings.Replace(pass, "{user}", user, -1)
			flag, err := doWithTimeOut(info, user, pass)
			if flag == true && err == nil {
				var result string
				if config.Domain != "" {
					result = fmt.Sprintf("[!] SMB %v:%v:%v\\%v %v", info.Host, info.Ports, config.Domain, user, pass)
				} else {
					result = fmt.Sprintf("[!] SMB %v:%v:%v %v", info.Host, info.Ports, user, pass)
				}
				utils.LogSuccess(result)
				return err
			}
		}
	}
	return tmperr
}

func doWithTimeOut(info *config.HostInfo, user string, pass string) (flag bool, err error) {
	signal := make(chan struct{})
	go func() {
		flag, err = smbConn(info, user, pass, signal)
	}()
	select {
	case <-signal:
		return flag, err
	case <-time.After(time.Duration(config.Timeout) * time.Second):
		return false, errors.New("time out")
	}
}

func smbConn(info *config.HostInfo, user string, pass string, signal chan struct{}) (flag bool, err error) {
	flag = false
	Host, Username, Password := info.Host, user, pass
	options := smb.Options{
		Host:        Host,
		Port:        445,
		User:        Username,
		Password:    Password,
		Domain:      config.Domain,
		Workstation: "",
	}

	session, err := smb.NewSession(options, false)
	if err == nil {
		session.Close()
		if session.IsAuthenticated {
			flag = true
		}
	}
	signal <- struct{}{}
	return flag, err
}
