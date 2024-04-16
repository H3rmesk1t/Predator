package v1

import (
	"Predator/pkg/common"
	"Predator/pkg/config"
	"Predator/pkg/utils"
	"Predator/pkg/xhttp"
	"fmt"
	"math/rand"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

func Start(Info config.HostInfo) {
	if Info.Host != "" && config.HostFile == "" {
		config.OutputFile = generateFileName(Info.Host)
	}

	utils.LogSuccess("[>] Scan result is saved to " + config.OutputFile)
	utils.LogSuccess("[>] Start alive scan...")
	hosts, err := common.ParseIP(Info.Host, config.HostFile, config.NoHosts)
	if err != nil {
		fmt.Println("[-] Can not find hosts")
		return
	}

	err = xhttp.Init()
	if err != nil {
		fmt.Println("[-] Can not initialization request")
		return
	}
	var ch = make(chan struct{}, config.Thread)
	var wg = sync.WaitGroup{}

	web := strconv.Itoa(config.PORTList["web"])
	ms17010 := strconv.Itoa(config.PORTList["ms17010"])

	if len(hosts) > 0 || len(config.HostPort) > 0 {
		if config.NoPing == false && len(hosts) > 1 || config.ScanType == "icmp" {
			hosts = CheckLive(hosts, config.Ping)
			fmt.Println("[*] Icmp alive hosts len is:", len(hosts))
		}
		if config.ScanType == "icmp" {
			utils.LogWG.Wait()
			return
		}
		common.GCMemory()

		var AlivePorts []string
		if config.ScanType == "webscan" {
			AlivePorts = common.NoPortScan(hosts, config.Ports)
		} else if config.ScanType == "netbios" {
			config.Ports = "139"
			AlivePorts = common.NoPortScan(hosts, config.Ports)
		} else if len(hosts) > 0 {
			AlivePorts = common.PortScan(hosts, config.Ports, config.Timeout)
			fmt.Println("[*] Alive Ports len is:", len(AlivePorts))
			if config.ScanType == "portscan" {
				utils.LogWG.Wait()
				return
			}
		}

		if len(config.HostPort) > 0 {
			AlivePorts = append(AlivePorts, config.HostPort...)
			AlivePorts = utils.RemoveDuplicateOfStringArray(AlivePorts)
			config.HostPort = nil
			fmt.Println("[*] Alive Ports len is:", len(AlivePorts))
		}
		common.GCMemory()

		var severports []string
		for _, port := range config.PORTList {
			severports = append(severports, strconv.Itoa(port))
		}

		fmt.Println()
		utils.LogSuccess("[>] Start info scan...")
		for _, targetIP := range AlivePorts {
			Info.Host, Info.Ports = strings.Split(targetIP, ":")[0], strings.Split(targetIP, ":")[1]
			if config.ScanType == "all" || config.ScanType == "main" {
				switch {
				case Info.Ports == "445":
					AddScan(ms17010, Info, &ch, &wg)
				case IsContain(severports, Info.Ports):
					AddScan(Info.Ports, Info, &ch, &wg)
				default:
					AddScan(web, Info, &ch, &wg)
				}
			} else {
				scanType := strconv.Itoa(config.PORTList[config.ScanType])
				AddScan(scanType, Info, &ch, &wg)
			}
		}
	}
	common.GCMemory()

	for _, u := range config.Urls {
		config.Url = u
		AddScan(web, Info, &ch, &wg)
	}
	common.GCMemory()
	wg.Wait()
	utils.LogWG.Wait()
	close(utils.Results)

	fmt.Printf("[*] %v/%v tasks has been done, consuming time: ", utils.End, utils.Num)
}

var Mutex = &sync.Mutex{}

func AddScan(scanType string, info config.HostInfo, ch *chan struct{}, wg *sync.WaitGroup) {
	*ch <- struct{}{}
	wg.Add(1)
	go func() {
		Mutex.Lock()
		utils.Num += 1
		Mutex.Unlock()
		ScanFunc(&scanType, &info)
		Mutex.Lock()
		utils.End += 1
		Mutex.Unlock()
		wg.Done()
		<-*ch
	}()
}

func ScanFunc(name *string, info *config.HostInfo) {
	f := reflect.ValueOf(PluginList[*name])
	in := []reflect.Value{reflect.ValueOf(info)}
	f.Call(in)
}

func IsContain(items []string, item string) bool {
	for _, i := range items {
		if i == item {
			return true
		}
	}
	return false
}

func generateFileName(host string) string {
	reg := regexp.MustCompile(`\b(?:[0-9]{1,3}\.){3}[0-9]{1,3}\b`)
	suffix := reg.FindString(host)

	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	s := make([]rune, 5)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}

	return fmt.Sprintf("%s-%s.txt", suffix, string(s))
}
