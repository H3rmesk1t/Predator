package ftp

import (
	"Predator/pkg/config"
	"Predator/pkg/utils"
	"fmt"
	"github.com/jlaffaye/ftp"
	"strings"
	"time"
)

func FtpScan(info *config.HostInfo) (tmperr error) {
	if config.IsBrutePass {
		return
	}
	flag, err := ftpConn(info, "anonymous", "")
	if flag && err == nil {
		return err
	}

	for _, user := range config.DefaultUsers["ftp"] {
		for _, pass := range config.DefaultPasswords {
			pass = strings.Replace(pass, "{user}", user, -1)
			flag, err := ftpConn(info, user, pass)
			if flag && err == nil {
				return err
			}
		}
	}
	return tmperr
}

func ftpConn(info *config.HostInfo, user string, pass string) (flag bool, err error) {
	flag = false
	Host, Port, Username, Password := info.Host, info.Ports, user, pass
	conn, err := ftp.DialTimeout(fmt.Sprintf("%v:%v", Host, Port), time.Duration(config.Timeout)*time.Second)
	if err == nil {
		err = conn.Login(Username, Password)
		if err == nil {
			flag = true
			result := fmt.Sprintf("[!] FTP %v:%v:%v %v", Host, Port, Username, Password)
			dirs, err := conn.List("")
			if err == nil {
				if len(dirs) > 0 {
					for i := 0; i < len(dirs); i++ {
						if len(dirs[i].Name) > 50 {
							result += "\n	[->]" + dirs[i].Name[:50]
						} else {
							result += "\n   [->]" + dirs[i].Name
						}
						if i == 5 {
							break
						}
					}
				}
			}
			utils.LogSuccess(result)
		}
	}
	return flag, err
}
