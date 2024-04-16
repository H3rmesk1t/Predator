package smb

import (
	"Predator/pkg/config"
	"Predator/pkg/utils"
	"fmt"
	"github.com/hirochachacha/go-smb2"
	"net"
	"os"
	"strings"
	"time"
)

func SmbScanTwo(info *config.HostInfo) (tmperr error) {
	if config.IsBrutePass {
		return nil
	}

	hasPrint := false
	hash := config.HashBytes
	for _, user := range config.DefaultUsers["smb"] {
	PASS:
		for _, pass := range config.DefaultPasswords {
			pass = strings.Replace(pass, "{user}", user, -1)
			flag, err, flag2 := smb2Con(info, user, pass, hash, hasPrint)
			if flag2 {
				hasPrint = true
			}
			if flag == true {
				var result string
				if config.Domain != "" {
					result = fmt.Sprintf("[!] SMB2 %v:%v:%v\\%v ", info.Host, info.Ports, config.Domain, user)
				} else {
					result = fmt.Sprintf("[!] SMB2 %v:%v:%v ", info.Host, info.Ports, user)
				}
				if len(hash) > 0 {
					result += "hash: " + config.Hash
				} else {
					result += pass
				}
				utils.LogSuccess(result)
				return err
			}
			if len(config.Hash) > 0 {
				break PASS
			}
		}
	}
	return tmperr
}

func smb2Con(info *config.HostInfo, user string, pass string, hash []byte, hasprint bool) (flag bool, err error, flag2 bool) {
	conn, err := net.DialTimeout("tcp", info.Host+":445", time.Duration(config.Timeout)*time.Second)
	if err != nil {
		return
	}
	defer conn.Close()
	initiator := smb2.NTLMInitiator{
		User:   user,
		Domain: config.Domain,
	}
	if len(hash) > 0 {
		initiator.Hash = hash
	} else {
		initiator.Password = pass
	}
	d := &smb2.Dialer{
		Initiator: &initiator,
	}

	s, err := d.Dial(conn)
	if err != nil {
		return
	}
	defer s.Logoff()
	names, err := s.ListSharenames()
	if err != nil {
		return
	}
	if !hasprint {
		var result string
		if config.Domain != "" {
			result = fmt.Sprintf("[*] SMB2-shares %v:%v:%v\\%v ", info.Host, info.Ports, config.Domain, user)
		} else {
			result = fmt.Sprintf("[*] SMB2-shares %v:%v:%v ", info.Host, info.Ports, user)
		}
		if len(hash) > 0 {
			result += "hash: " + config.Hash
		} else {
			result += pass
		}
		result = fmt.Sprintf("%v shares: %v", result, names)
		utils.LogSuccess(result)
		flag2 = true
	}
	fs, err := s.Mount("C$")
	if err != nil {
		return
	}
	defer fs.Umount()
	path := `Windows\win.ini`
	f, err := fs.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		return
	}
	defer f.Close()
	flag = true

	return
}
