package common

import (
	"Predator/pkg/config"
	"Predator/pkg/utils"
	"encoding/hex"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func Process(Info *config.HostInfo) {
	processUser()
	processPassword()
	processUrl()
	processPort(Info)
	checkParams(Info)
	checkScanType()
}

func processUser() {
	if config.Username == "" && config.UserFile == "" {
		return
	}

	var Usernames []string
	if config.Username != "" {
		Usernames = strings.Split(config.Username, ",")
	}

	if config.UserFile != "" {
		users, err := utils.ReadFile(config.UserFile)
		if err == nil {
			for _, user := range users {
				if user != "" {
					Usernames = append(Usernames, user)
				}
			}
		}
	}

	Usernames = utils.RemoveDuplicateOfStringArray(Usernames)
	for name := range config.DefaultUsers {
		config.DefaultUsers[name] = Usernames
	}
}

func processPassword() {
	if config.Password == "" && config.PasswordFile == "" {
		return
	}

	var Passwords []string
	if config.Password != "" {
		Passwords = strings.Split(config.Password, ",")
	}

	if config.PasswordFile != "" {
		pass, err := utils.ReadFile(config.PasswordFile)
		if err == nil {
			for _, p := range pass {
				if p != "" {
					Passwords = append(Passwords, p)
				}
			}
		}
	}

	config.DefaultPasswords = utils.RemoveDuplicateOfStringArray(Passwords)
}

func processUrl() {
	if config.Url == "" && config.UrlFile == "" {
		return
	}

	if config.Url != "" {
		Urls := strings.Split(config.Url, ",")
		tempUrls := make(map[string]struct{})
		for _, u := range Urls {
			if _, ok := tempUrls[u]; !ok {
				tempUrls[u] = struct{}{}
				str, err := url.Parse(u)
				if err == nil && str.Scheme != "" && str.Host != "" {
					config.Urls = append(config.Urls, u)
				}
			}
		}
	}

	if config.UrlFile != "" {
		Urls, err := utils.ReadFile(config.UrlFile)
		if err == nil {
			tempUrls := make(map[string]struct{})
			for _, u := range Urls {
				if _, ok := tempUrls[u]; !ok {
					tempUrls[u] = struct{}{}
					str, err := url.Parse(u)
					if err == nil && str.Scheme != "" && str.Host != "" {
						config.Urls = append(config.Urls, u)
					}
				}
			}
		}
	}
}

func processPort(Info *config.HostInfo) {
	if config.PortFile != "" {
		ports, err := utils.ReadFile(config.PortFile)
		if err != nil {
			newport := ""
			for _, port := range ports {
				if port != "" {
					newport += port + ","
				}
			}
			Info.Ports = newport
		}
	}
}

func checkParams(Info *config.HostInfo) {
	if Info.Host == "" && config.HostFile == "" && config.Url == "" && config.UrlFile == "" {
		fmt.Println("[-] target is none...")
		flag.Usage()
		os.Exit(0)
	}

	if config.BruteThread <= 0 {
		config.BruteThread = 5
	}

	if config.NoSave == true {
		config.IsSave = false
	}

	if config.Ports == config.DefaultPorts {
		config.Ports += "," + config.DefaultWebPorts
	}
	if config.AddPorts != "" {
		config.AddPorts = strings.TrimSuffix(config.AddPorts, ",")
		if strings.HasSuffix(config.Ports, ",") {
			config.Ports += config.AddPorts
		} else {
			config.Ports += "," + config.AddPorts
		}
	}
	config.Ports = utils.RemoveDuplicateOfString(config.Ports)

	if config.AddUsers != "" {
		user := strings.Split(config.AddUsers, ",")
		for name := range config.DefaultUsers {
			config.DefaultUsers[name] = append(config.DefaultUsers[name], user...)
			config.DefaultUsers[name] = utils.RemoveDuplicateOfStringArray(config.DefaultUsers[name])
		}
	}

	if config.AddPasswords != "" {
		pass := strings.Split(config.AddPasswords, ",")
		config.DefaultPasswords = append(config.DefaultPasswords, pass...)
		config.DefaultPasswords = utils.RemoveDuplicateOfStringArray(config.DefaultPasswords)
	}

	if config.Hash != "" && len(config.Hash) != 32 {
		fmt.Println("[-] Hash is error, hash length must be 32")
		os.Exit(0)
	} else if config.Hash != "" {
		var err error
		config.HashBytes, err = hex.DecodeString(config.Hash)
		if err != nil {
			fmt.Println("[-] Hash is error, hex decode error")
			os.Exit(0)
		}
	}
}

func checkScanType() {
	if _, ok := config.PORTList[config.ScanType]; !ok {
		fmt.Println("[-] The specified scan type does not exist")
		fmt.Println("-m")
		for name := range config.PORTList {
			fmt.Println("	[" + name + "]")
		}
		os.Exit(0)
	}

	if (config.ScanType != "all") && (config.Ports == utils.RemoveDuplicateOfString(config.DefaultPorts+","+config.DefaultWebPorts)) {
		switch config.ScanType {
		case "wmi":
			config.Ports = "135"
		case "smb":
			config.Ports = "445"
		case "netbios":
			config.Ports = "135,137,139,445"
		case "web":
			config.Ports = config.DefaultWebPorts
		case "ms17010":
			config.Ports = "445"
		case "smbghost":
			config.Ports = "445"
		case "portscan":
			config.Ports = utils.RemoveDuplicateOfString(config.DefaultPorts + "," + config.DefaultWebPorts)
		case "main":
			config.Ports = config.DefaultPorts
		default:
			port, _ := config.PORTList[config.ScanType]
			config.Ports = strconv.Itoa(port)
		}
		fmt.Println("[-] -m ", config.ScanType, ". Start scan:", config.Ports)
	}
}
