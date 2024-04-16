package v1

import (
	"Predator/pkg/config"
	"Predator/pkg/utils"
	"bytes"
	"fmt"
	"golang.org/x/net/icmp"
	"net"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	AliveHosts []string
	ExistHosts = make(map[string]struct{})
	livewg     sync.WaitGroup
)

func CheckLive(hostslist []string, isPing bool) []string {
	chanHosts := make(chan string, len(hostslist))
	go func() {
		for ip := range chanHosts {
			if _, ok := ExistHosts[ip]; !ok && IsContain(hostslist, ip) {
				ExistHosts[ip] = struct{}{}
				if config.Silent == false {
					fmt.Printf("[icmp] Target %-15s is alive.\n", ip)
				} else {
					fmt.Printf("[ping] Target %-15s is alive.\n", ip)
				}
				AliveHosts = append(AliveHosts, ip)
			}
			livewg.Done()
		}
	}()

	if isPing == true {
		RunPing(hostslist, chanHosts)
	} else {
		conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
		if err == nil {
			RunIcmpOne(hostslist, conn, chanHosts)
		} else {
			utils.LogError(err)
			conn, err := net.DialTimeout("ip4:icmp", "127.0.0.1", 5*time.Second)
			defer func() {
				if conn != nil {
					conn.Close()
				}
			}()

			if err == nil {
				RunIcmpTwo(hostslist, chanHosts)
			} else {
				utils.LogError(err)
				fmt.Println("[-] The current user permissions unable to send icmp packets, start ping...")
				RunPing(hostslist, chanHosts)
			}
		}
	}

	livewg.Wait()
	close(chanHosts)

	if len(hostslist) > 1024 {
		arrTop, arrLen := ArrayCountValueTop(AliveHosts, config.LiveTop, true)
		for i := 0; i < len(arrTop); i++ {
			output := fmt.Sprintf("[*] LiveTop %-16s 段存活数量为: %d", arrTop[i]+".0.0/16", arrLen[i])
			utils.LogSuccess(output)
		}
	}
	if len(hostslist) > 256 {
		arrTop, arrLen := ArrayCountValueTop(AliveHosts, config.LiveTop, false)
		for i := 0; i < len(arrTop); i++ {
			output := fmt.Sprintf("[*] LiveTop %-16s 段存活数量为: %d", arrTop[i]+".0/24", arrLen[i])
			utils.LogSuccess(output)
		}
	}

	return AliveHosts
}

func RunPing(hostsList []string, chanHosts chan string) {
	var wg sync.WaitGroup

	limiter := make(chan struct{}, 64)
	for _, host := range hostsList {
		wg.Add(1)
		limiter <- struct{}{}
		go func(host string) {
			if ExecPing(host) {
				livewg.Add(1)
				chanHosts <- host
			}
			<-limiter
			wg.Done()
		}(host)
	}
	wg.Wait()
}

func ExecPing(ip string) bool {
	var command *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		command = exec.Command("module", "/c", "ping -n 1 -w 1 "+ip+" && echo true || echo false")
	case "darwin":
		command = exec.Command("/bin/bash", "-c", "ping -c 1 -W 1 "+ip+" && echo true || echo false")
	default:
		command = exec.Command("/bin/bash", "-c", "ping -c 1 -w 1 "+ip+" && echo true || echo false")
	}

	buffer := bytes.Buffer{}
	command.Stdout = &buffer
	err := command.Start()
	if err != nil {
		return false
	}
	if err = command.Wait(); err != nil {
		return false
	} else {
		if strings.Contains(buffer.String(), "true") && strings.Count(buffer.String(), ip) > 2 {
			return true
		} else {
			return false
		}
	}
}

func RunIcmpOne(hostslist []string, conn *icmp.PacketConn, chanHosts chan string) {
	flag := false

	go func() {
		for {
			if flag == true {
				return
			}
			msg := make([]byte, 100)
			_, sourceIP, _ := conn.ReadFrom(msg)
			if sourceIP != nil {
				livewg.Add(1)
				chanHosts <- sourceIP.String()
			}
		}
	}()

	for _, host := range hostslist {
		dst, _ := net.ResolveIPAddr("ip", host)
		IcmpByte := makemsg(host)
		conn.WriteTo(IcmpByte, dst)
	}
	start := time.Now()
	for {
		if len(AliveHosts) == len(hostslist) {
			break
		}
		since := time.Since(start)
		var wait time.Duration
		switch {
		case len(hostslist) <= 256:
			wait = time.Second * 3
		default:
			wait = time.Second * 5
		}
		if since > wait {
			break
		}
	}
	flag = true
	err := conn.Close()
	if err != nil {
		return
	}
}

func RunIcmpTwo(hostslist []string, chanHosts chan string) {
	num := 1024
	if len(hostslist) < num {
		num = len(hostslist)
	}

	var wg sync.WaitGroup
	limiter := make(chan struct{}, num)
	for _, host := range hostslist {
		wg.Add(1)
		limiter <- struct{}{}
		go func(host string) {
			if icmpalive(host) {
				livewg.Add(1)
				chanHosts <- host
			}
			<-limiter
			wg.Done()
		}(host)
	}
	wg.Wait()
	close(limiter)
}

func icmpalive(host string) bool {
	startTime := time.Now()
	conn, err := net.DialTimeout("ip4:icmp", host, 5*time.Second)
	if err != nil {
		return false
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
		}
	}(conn)
	if err := conn.SetDeadline(startTime.Add(5 * time.Second)); err != nil {
		return false
	}
	msg := makemsg(host)
	if _, err := conn.Write(msg); err != nil {
		return false
	}

	receive := make([]byte, 64)
	if _, err := conn.Read(receive); err != nil {
		return false
	}

	return true
}

func makemsg(host string) []byte {
	msg := make([]byte, 40)
	id0, id1 := genIdentifier(host)
	msg[0] = 8
	msg[1] = 0
	msg[2] = 0
	msg[3] = 0
	msg[4], msg[5] = id0, id1
	msg[6], msg[7] = genSequence(1)
	check := checkSum(msg[0:40])
	msg[2] = byte(check >> 8)
	msg[3] = byte(check & 255)
	return msg
}

func checkSum(msg []byte) uint16 {
	sum := 0
	length := len(msg)
	for i := 0; i < length-1; i += 2 {
		sum += int(msg[i])*256 + int(msg[i+1])
	}
	if length%2 == 1 {
		sum += int(msg[length-1]) * 256
	}
	sum = (sum >> 16) + (sum & 0xffff)
	sum = sum + (sum >> 16)
	answer := uint16(^sum)
	return answer
}

func genSequence(v int16) (byte, byte) {
	ret1 := byte(v >> 8)
	ret2 := byte(v & 255)
	return ret1, ret2
}

func genIdentifier(host string) (byte, byte) {
	return host[0], host[1]
}

func ArrayCountValueTop(arrInit []string, length int, flag bool) (arrTop []string, arrLen []int) {
	if len(arrInit) == 0 {
		return
	}
	arrMap1 := make(map[string]int)
	arrMap2 := make(map[string]int)
	for _, value := range arrInit {
		line := strings.Split(value, ".")
		if len(line) == 4 {
			if flag {
				value = fmt.Sprintf("%s.%s", line[0], line[1])
			} else {
				value = fmt.Sprintf("%s.%s.%s", line[0], line[1], line[2])
			}
		}
		if arrMap1[value] != 0 {
			arrMap1[value]++
		} else {
			arrMap1[value] = 1
		}
	}
	for k, v := range arrMap1 {
		arrMap2[k] = v
	}

	i := 0
	for range arrMap1 {
		var maxCountKey string
		var maxCountVal = 0
		for key, val := range arrMap2 {
			if val > maxCountVal {
				maxCountVal = val
				maxCountKey = key
			}
		}
		arrTop = append(arrTop, maxCountKey)
		arrLen = append(arrLen, maxCountVal)
		i++
		if i >= length {
			return
		}
		delete(arrMap2, maxCountKey)
	}
	return
}
