package ssh

import (
	"Predator/pkg/config"
	"Predator/pkg/utils"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net"
	"strings"
	"time"
)

func SshScan(info *config.HostInfo) (tamperer error) {
	if config.IsBrutePass {
		return
	}
	for _, user := range config.DefaultUsers["ssh"] {
		for _, pass := range config.DefaultPasswords {
			pass = strings.Replace(pass, "{user}", user, -1)
			flag, err := sshConn(info, user, pass)
			if flag == true && err == nil {
				return err
			}
			if config.SshKey != "" {
				return err
			}
		}
	}
	return tamperer
}

func sshConn(info *config.HostInfo, user string, pass string) (flag bool, err error) {
	flag = false
	Host, Port, Username, Password := info.Host, info.Ports, user, pass
	var Auth []ssh.AuthMethod
	if config.SshKey != "" {
		pemBytes, err := ioutil.ReadFile(config.SshKey)
		if err != nil {
			return false, errors.New("read key failed" + err.Error())
		}
		signer, err := ssh.ParsePrivateKey(pemBytes)
		if err != nil {
			return false, errors.New("parse key failed" + err.Error())
		}
		Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}
	} else {
		Auth = []ssh.AuthMethod{ssh.Password(Password)}
	}

	sshClientConfig := &ssh.ClientConfig{
		User:    Username,
		Auth:    Auth,
		Timeout: time.Duration(config.Timeout) * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%v:%v", Host, Port), sshClientConfig)
	if err == nil {
		defer client.Close()
		session, err := client.NewSession()
		if err == nil {
			defer session.Close()
			flag = true
			var result string
			result = fmt.Sprintf("[!] SSH %v:%v:%v %v", Host, Port, Username, Password)
			if config.SshKey != "" {
				result = fmt.Sprintf("[!] SSH %v:%v sshkey correct", Host, Port)
			}
			utils.LogSuccess(result)
		}
	}

	return flag, err
}
