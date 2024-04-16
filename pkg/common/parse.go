package common

import (
	"Predator/pkg/config"
	"Predator/pkg/utils"
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var ParseIPErr = errors.New("[-] Host parsing error\n" +
	"format: \n" +
	"192.168.1.1\n" +
	"192.168.1.1/8\n" +
	"192.168.1.1/16\n" +
	"192.168.1.1/24\n" +
	"192.168.1.1,192.168.1.2\n" +
	"192.168.1.1-192.168.255.255\n" +
	"192.168.1.1-255")

func ParseIP(host string, filename string, nohosts ...string) (hosts []string, err error) {
	if filename == "" && strings.Contains(host, ":") {
		temp := strings.Split(host, ":")
		if len(temp) == 2 {
			host = temp[0]
			hosts = ParseIPs(host)
			config.Ports = temp[1]
		}
	} else {
		hosts = ParseIPs(host)
		if filename != "" {
			var fileHosts []string
			fileHosts, _ = ReadIPFile(filename)
			hosts = append(hosts, fileHosts...)
		}
	}

	if len(nohosts) > 0 {
		noHost := nohosts[0]
		if noHost != "" {
			noHosts := ParseIPs(noHost)
			if len(noHosts) > 0 {
				temp := map[string]struct{}{}
				for _, host := range hosts {
					temp[host] = struct{}{}
				}

				for _, host := range noHosts {
					delete(temp, host)
				}

				var newDatas []string
				for host := range temp {
					newDatas = append(newDatas, host)
				}
				hosts = newDatas
				sort.Strings(hosts)
			}
		}
	}

	hosts = utils.RemoveDuplicateOfStringArray(hosts)
	if len(hosts) == 0 && len(config.HostPort) == 0 && host != "" && filename != "" {
		err = ParseIPErr
	}
	return
}

func ParseIPs(ip string) (hosts []string) {
	if strings.Contains(ip, ",") {
		IPList := strings.Split(ip, ",")
		var ips []string
		for _, ip := range IPList {
			ips = HandleIP(ip)
			hosts = append(hosts, ips...)
		}
	} else {
		hosts = HandleIP(ip)
	}
	return hosts
}

func HandleIP(ip string) []string {
	reg := regexp.MustCompile(`[a-zA-Z]+`)
	switch {
	case ip == "192":
		return HandleIP("192.168.0.0/8")
	case ip == "172":
		return HandleIP("172.16.0.0/12")
	case ip == "10":
		return HandleIP("10.0.0.0/8")
	case strings.HasSuffix(ip, "/8"):
		return HandleIP8(ip)
	case strings.Contains(ip, "/"):
		return HandleIP2(ip)
	case reg.MatchString(ip):
		return []string{ip}
	case strings.Contains(ip, "-"):
		return HandleIP1(ip)
	default:
		testIP := net.ParseIP(ip)
		if testIP == nil {
			return nil
		}
		return []string{ip}
	}
}

func HandleIP1(ip string) []string {
	IPRange := strings.Split(ip, "-")
	testIP := net.ParseIP(IPRange[0])
	var AllIP []string
	if len(IPRange[1]) < 4 {
		Range, err := strconv.Atoi(IPRange[1])
		if testIP == nil || Range > 255 || err != nil {
			return nil
		}
		SplitIP := strings.Split(IPRange[0], ".")
		ip1, err1 := strconv.Atoi(SplitIP[3])
		ip2, err2 := strconv.Atoi(IPRange[1])
		PrefixIP := strings.Join(SplitIP[0:3], ".")
		if ip1 > ip2 || err1 != nil || err2 != nil {
			return nil
		}
		for i := ip1; i <= ip2; i++ {
			AllIP = append(AllIP, PrefixIP+"."+strconv.Itoa(i))
		}
	} else {
		SplitIP1 := strings.Split(IPRange[0], ".")
		SplitIP2 := strings.Split(IPRange[1], ".")
		if len(SplitIP1) != 4 || len(SplitIP2) != 4 {
			return nil
		}
		start, end := [4]int{}, [4]int{}
		for i := 0; i < 4; i++ {
			ip1, err1 := strconv.Atoi(SplitIP1[i])
			ip2, err2 := strconv.Atoi(SplitIP2[i])
			if ip1 > ip2 || err1 != nil || err2 != nil {
				return nil
			}
			start[i], end[i] = ip1, ip2
		}
		startNum := start[0]<<24 | start[1]<<16 | start[2]<<8 | start[3]
		endNum := end[0]<<24 | end[1]<<16 | end[2]<<8 | end[3]
		for num := startNum; num <= endNum; num++ {
			ip := strconv.Itoa((num>>24)&0xff) + "." + strconv.Itoa((num>>16)&0xff) + "." + strconv.Itoa((num>>8)&0xff) + "." + strconv.Itoa((num)&0xff)
			AllIP = append(AllIP, ip)
		}
	}
	return AllIP
}

func HandleIP2(host string) (hosts []string) {
	_, ipNet, err := net.ParseCIDR(host)
	if err != nil {
		return
	}
	hosts = HandleIP1(IPRange(ipNet))
	return
}

func HandleIP8(ip string) []string {
	realIP := ip[:len(ip)-2]
	testIP := net.ParseIP(realIP)

	if testIP == nil {
		return nil
	}

	IPRange := strings.Split(ip, ".")[0]
	var AllIP []string
	for a := 0; a <= 255; a++ {
		for b := 0; b <= 255; b++ {
			AllIP = append(AllIP, fmt.Sprintf("%s.%d.%d.%d", IPRange, a, b, 1))
			AllIP = append(AllIP, fmt.Sprintf("%s.%d.%d.%d", IPRange, a, b, 2))
			AllIP = append(AllIP, fmt.Sprintf("%s.%d.%d.%d", IPRange, a, b, 4))
			AllIP = append(AllIP, fmt.Sprintf("%s.%d.%d.%d", IPRange, a, b, 5))
			AllIP = append(AllIP, fmt.Sprintf("%s.%d.%d.%d", IPRange, a, b, RandInt(6, 25)))
			AllIP = append(AllIP, fmt.Sprintf("%s.%d.%d.%d", IPRange, a, b, RandInt(26, 55)))
			AllIP = append(AllIP, fmt.Sprintf("%s.%d.%d.%d", IPRange, a, b, RandInt(56, 75)))
			AllIP = append(AllIP, fmt.Sprintf("%s.%d.%d.%d", IPRange, a, b, RandInt(76, 100)))
			AllIP = append(AllIP, fmt.Sprintf("%s.%d.%d.%d", IPRange, a, b, RandInt(101, 125)))
			AllIP = append(AllIP, fmt.Sprintf("%s.%d.%d.%d", IPRange, a, b, RandInt(126, 150)))
			AllIP = append(AllIP, fmt.Sprintf("%s.%d.%d.%d", IPRange, a, b, RandInt(151, 175)))
			AllIP = append(AllIP, fmt.Sprintf("%s.%d.%d.%d", IPRange, a, b, RandInt(176, 200)))
			AllIP = append(AllIP, fmt.Sprintf("%s.%d.%d.%d", IPRange, a, b, RandInt(201, 225)))
			AllIP = append(AllIP, fmt.Sprintf("%s.%d.%d.%d", IPRange, a, b, RandInt(226, 253)))
			AllIP = append(AllIP, fmt.Sprintf("%s.%d.%d.%d", IPRange, a, b, 254))
		}
	}
	return AllIP
}

func IPRange(c *net.IPNet) string {
	start := c.IP.String()
	mask := c.Mask
	bcst := make(net.IP, len(c.IP))
	copy(bcst, c.IP)
	for i := 0; i < len(mask); i++ {
		ipIdx := len(bcst) - i - 1
		bcst[ipIdx] = c.IP[ipIdx] | ^mask[len(mask)-i-1]
	}
	end := bcst.String()
	return fmt.Sprintf("%s-%s", start, end)
}

func RandInt(min, max int) int {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Intn(max-min) + min
}

func ParsePort(ports string) (scanPorts []int) {
	if ports == "" {
		return
	}
	slices := strings.Split(ports, ",")
	for _, port := range slices {
		port = strings.TrimSpace(port)
		if port == "" {
			continue
		}

		if config.PortGroup[port] != "" {
			port = config.PortGroup[port]
			scanPorts = append(scanPorts, ParsePort(port)...)
			continue
		}
		upper := port
		if strings.Contains(port, "-") {
			ranges := strings.Split(port, "-")
			if len(ranges) < 2 {
				continue
			}

			startPort, _ := strconv.Atoi(ranges[0])
			endPort, _ := strconv.Atoi(ranges[1])
			if startPort < endPort {
				port = ranges[0]
				upper = ranges[1]
			} else {
				port = ranges[1]
				upper = ranges[0]
			}
		}
		start, _ := strconv.Atoi(port)
		end, _ := strconv.Atoi(upper)
		for i := start; i <= end; i++ {
			if i > 65535 || i < 1 {
				continue
			}
			scanPorts = append(scanPorts, i)
		}
	}
	scanPorts = utils.RemoveDuplicatePort(scanPorts)
	return scanPorts
}

func ReadIPFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("[-] Open %s error, %v", filename, err)
		os.Exit(0)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	var content []string
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			text := strings.Split(line, ":")
			if len(text) == 2 {
				port := strings.Split(text[1], " ")[0]
				num, err := strconv.Atoi(port)
				if err != nil || (num < 1 || num > 65535) {
					continue
				}
				hosts := ParseIPs(text[0])
				for _, host := range hosts {
					config.HostPort = append(config.HostPort, fmt.Sprintf("%s:%s", host, port))
				}
			} else {
				host := ParseIPs(line)
				content = append(content, host...)
			}
		}
	}
	return content, nil
}
