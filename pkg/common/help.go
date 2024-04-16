package common

import (
	"Predator/pkg/config"
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
)

type pFlagSet struct {
	*flag.FlagSet
	cmdCommand string
}

func init() {
	go func() {
		for {
			GCMemory()
			time.Sleep(10 * time.Second)
		}
	}()
}

func GCMemory() {
	runtime.GC()
	debug.FreeOSMemory()
}

func Help(Info *config.HostInfo) {
	scanCmd := &pFlagSet{
		FlagSet:    flag.NewFlagSet("scan", flag.ExitOnError),
		cmdCommand: "comprehensive scanning of internal and external network assets.",
	}
	scanCmd.StringVar(&Info.Host, "h", "", "select target address of the host to scan, for example: 192.168.1.1 | 192.168.1.1-255 | 192.168.1.1,192.168.1.2")
	scanCmd.StringVar(&config.HostFile, "hf", "", "use host file, -hf ip.txt")
	scanCmd.StringVar(&config.NoHosts, "nh", "", "set the hosts no scan, -hn 192.168.1.1/24")
	scanCmd.StringVar(&config.Ports, "p", config.DefaultPorts, "select port to scan, for example: 22 | 1-65535 | 22,80,3306")
	scanCmd.StringVar(&config.NoPorts, "np", "", "the ports no scan, -pn 445")
	scanCmd.StringVar(&config.AddPorts, "ap", "", "add port base DefaultPorts, -pa 3389")
	scanCmd.StringVar(&config.PortFile, "pf", "", "use port File, -pf port.txt")
	scanCmd.StringVar(&config.Url, "u", "", "url, -u url")
	scanCmd.StringVar(&config.UrlFile, "uf", "", "urlfile, -uf url.txt")
	scanCmd.StringVar(&config.Username, "user", "", "username, -user xxx")
	scanCmd.StringVar(&config.AddUsers, "auser", "", "add a user base DefaultUsers, -auser xxx")
	scanCmd.StringVar(&config.UserFile, "userf", "", "use username file, -userf user.txt")
	scanCmd.StringVar(&config.Password, "pwd", "", "password, -pwd xxx")
	scanCmd.StringVar(&config.AddPasswords, "apwd", "", "add a password base DefaultPasswords, -apwd xxx")
	scanCmd.StringVar(&config.PasswordFile, "pwdf", "", "use password file, -pwdf pwd.txt")
	scanCmd.BoolVar(&config.NoPoc, "nopoc", false, "not to scan web vul, -nopoc")
	scanCmd.StringVar(&config.PocInfo.PocName, "pocname", "", "use the pocs these contain pocname, -pocname weblogic")
	scanCmd.IntVar(&config.PocNum, "num", 25, "pocs rate, -num xxx")
	scanCmd.StringVar(&config.PocFile, "pocf", "", "use pocs file path, -pf pocs.txt")
	scanCmd.StringVar(&config.SshKey, "sk", "", "sshkey file (id_rsa), -sk xxx")
	scanCmd.StringVar(&config.Domain, "dc", "", "smb domain, -dc xxx")
	scanCmd.StringVar(&config.Hash, "hash", "", "set hash, -hash xxx")
	scanCmd.BoolVar(&config.DnsLog, "dns", false, "using dnslog pocs, -dns")
	scanCmd.BoolVar(&config.IsCheckFastjson, "fastjson", false, "check fastjson vuln, -fastjson")
	scanCmd.BoolVar(&config.IsCheckLog4j2, "log4", false, "check log4j2 vuln, -log4")
	scanCmd.BoolVar(&config.IsCheckSpringBoot, "springboot", false, "check springboot vuln. -springboot")
	scanCmd.StringVar(&config.ScanType, "m", "all", "select scan type, -m ssh")
	scanCmd.BoolVar(&config.Ping, "ping", false, "use ping replace icmp, -ping")
	scanCmd.BoolVar(&config.NoPing, "noping", false, "not to ping, -noping")
	scanCmd.Int64Var(&config.Timeout, "time", 5, "set timeout when scan, -time xxx")
	scanCmd.Int64Var(&config.WaitTime, "debug", 60, "every time to LogErr, -debug xxx")
	scanCmd.Int64Var(&config.WebTimeout, "wt", 3, "set web timeout, -wt xxx")
	scanCmd.BoolVar(&config.IsBrutePass, "nobp", false, "not to brute password, -bp")
	scanCmd.IntVar(&config.BruteThread, "br", 5, "brute threads when scan, -br")
	scanCmd.IntVar(&config.Thread, "thread", 600, "thread nums when scan, -thread xxx")
	scanCmd.IntVar(&config.LiveTop, "top", 10, "show live len top, -top xxx")
	scanCmd.BoolVar(&config.Silent, "s", false, "silent scan, -s")
	scanCmd.BoolVar(&config.NoColor, "nc", false, "no color cli, -nc")
	scanCmd.StringVar(&config.Cookie, "cookie", "", "set pocs cookie, -cookie rememberMe=login")
	scanCmd.BoolVar(&config.NoSave, "nosave", false, "not to save output log, -nosave")
	scanCmd.StringVar(&config.OutputFile, "o", "output.txt", "output results to file, -o xxx")

	spyCmd := &pFlagSet{
		FlagSet:    flag.NewFlagSet("spy", flag.ExitOnError),
		cmdCommand: "quickly detect reachable network segments on the internal network.",
	}
	spyCmd.StringVar(&config.SpyModule, "sm", "", "icmp protocol to spy, use icmp/ping/tcp/udp")
	spyCmd.StringVar(&config.SpyTcpPort, "stp", "21,22,23,80,135,139,443,445,3389,8080", "specify tcp port to spy (default: 21, 22, 23, 80, 135, 139, 443, 445, 3389, 8080)")
	spyCmd.StringVar(&config.SpyUdpPort, "sup", "53,123,137,161,520,523,1645,1701,1900,5353", "specify udp port to spy (default: 53, 123, 137, 161, 520, 523, 1645, 1701, 1900, 5353)")
	spyCmd.StringVar(&config.SpyCidr, "sc", "", "specify detection CIDR (for example: 172.16.0.0/12)")
	spyCmd.StringVar(&config.SpyOutput, "so", "spy.txt", "the path to save the surviving network segment results (default: spy.txt)")
	spyCmd.StringVar(&config.SpyEnd, "se", "", "specify the last number of the IP (default: 1, 254, 2, 255)")
	spyCmd.IntVar(&config.SpyTimes, "st", 1, "number of echo request messages be sent (default: 1)")
	spyCmd.IntVar(&config.SpyRandom, "sr", 1, "the number of digits at the end of the IP random number (default: 1)")
	spyCmd.IntVar(&config.SpyThread, "sth", 50, "number of concurrencies (default: 50)")
	spyCmd.IntVar(&config.SpyTimeout, "sto", 100, "packet sending timeout in milliseconds (default: 100)")
	spyCmd.BoolVar(&config.SpyRapid, "sx", false, "rapid detection mode (default: false)")
	spyCmd.BoolVar(&config.SpySpecial, "ssp", false, "whether to detect special intranets (default: false)")
	spyCmd.BoolVar(&config.SpyForce, "sf", false, "force detection of all generated IPs (default: false)")
	spyCmd.BoolVar(&config.SpySilent, "ssi", false, "only show surviving network segments (default: false)")
	spyCmd.BoolVar(&config.SpyDebug, "sd", false, "show debugging information (default: false)")

	subcommands := map[string]*pFlagSet{
		scanCmd.Name(): scanCmd,
		spyCmd.Name():  spyCmd,
	}

	useAge := func() {
		fmt.Printf("Usage: Predator Command\n\n")
		for _, sub := range subcommands {
			if sub.Name() != "" {
				v := strings.ToUpper(sub.Name()[:1]) + sub.Name()[1:]
				fmt.Printf("%s : %s\n", v, sub.cmdCommand)
				sub.PrintDefaults()
				fmt.Println()
			}
		}
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		fmt.Println("Error, parameters used in the module must be specified")
		useAge()
	}

	cmd := subcommands[os.Args[1]]
	if cmd == nil {
		useAge()
	}

	if os.Args[1] == "scan" {
		config.ModuleFlag = true
	} else if os.Args[1] == "spy" {
		config.ModuleFlag = false
	}

	err := cmd.Parse(os.Args[2:])
	if err != nil {
		return
	}
}
