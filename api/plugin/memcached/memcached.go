package memcached

import (
	"Predator/pkg/config"
	"Predator/pkg/utils"
	"Predator/pkg/xhttp"
	"fmt"
	"strings"
	"time"
)

func MemcachedScan(info *config.HostInfo) (err error) {
	realhost := fmt.Sprintf("%s:%v", info.Host, info.Ports)
	client, err := xhttp.WrapperTcpWithTimeout("tcp", realhost, time.Duration(config.Timeout)*time.Second)
	defer func() {
		if client != nil {
			client.Close()
		}
	}()
	if err == nil {
		err = client.SetDeadline(time.Now().Add(time.Duration(config.Timeout) * time.Second))
		if err == nil {
			_, err = client.Write([]byte("stats\n"))
			if err == nil {
				rev := make([]byte, 1024)
				n, err := client.Read(rev)
				if err == nil {
					if strings.Contains(string(rev[:n]), "STAT") {
						result := fmt.Sprintf("[!] Memcached %s unauthorized", realhost)
						utils.LogSuccess(result)
					}
				}
			}
		}
	}
	return err
}
