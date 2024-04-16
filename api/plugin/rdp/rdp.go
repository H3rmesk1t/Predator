package rdp

import (
	"Predator/pkg/config"
	"Predator/pkg/utils"
	"Predator/pkg/xhttp"
	"errors"
	"fmt"
	"github.com/tomatome/grdp/core"
	"github.com/tomatome/grdp/glog"
	"github.com/tomatome/grdp/protocol/nla"
	"github.com/tomatome/grdp/protocol/pdu"
	"github.com/tomatome/grdp/protocol/rfb"
	"github.com/tomatome/grdp/protocol/sec"
	"github.com/tomatome/grdp/protocol/t125"
	"github.com/tomatome/grdp/protocol/tpkt"
	"github.com/tomatome/grdp/protocol/x224"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type BruteList struct {
	user string
	pass string
}

func RdpScan(info *config.HostInfo) (tamperer error) {
	if config.IsBrutePass {
		return
	}

	var wg sync.WaitGroup
	var signal bool
	var num = 0
	var all = len(config.DefaultUsers["rdp"]) * len(config.DefaultPasswords)
	var mutex sync.Mutex

	brList := make(chan BruteList)
	port, _ := strconv.Atoi(info.Ports)

	for i := 0; i < config.BruteThread; i++ {
		wg.Add(1)
		go worker(info.Host, config.Domain, port, &wg, brList, &signal, &num, all, &mutex, config.Timeout)
	}

	for _, user := range config.DefaultUsers["rdp"] {
		for _, pass := range config.DefaultPasswords {
			pass = strings.Replace(pass, "{user}", user, -1)
			brList <- BruteList{user, pass}
		}
	}
	close(brList)
	go func() {
		wg.Wait()
		signal = true
	}()

	for !signal {
	}

	return tamperer
}

func worker(host, domain string, port int, wg *sync.WaitGroup, brlist chan BruteList, signal *bool, num *int, all int, mutex *sync.Mutex, timeout int64) {
	defer wg.Done()
	for one := range brlist {
		if *signal == true {
			return
		}
		go incrNum(num, mutex)
		user, pass := one.user, one.pass
		flag, err := rdpConn(host, domain, user, pass, port, timeout)
		if flag == true && err == nil {
			var result string
			if domain != "" {
				result = fmt.Sprintf("[!] RDP %v:%v:%v\\%v %v", host, port, domain, user, pass)
			} else {
				result = fmt.Sprintf("[!] RDP %v:%v:%v %v", host, port, user, pass)
			}
			utils.LogSuccess(result)
			*signal = true
			return
		}
		//else {
		//	errLog := fmt.Sprintf("[-] (%v/%v) rdp %v:%v %v %v %v", *num, all, host, port, user, pass, err)
		//	utils.LogError(errLog)
		//}
	}
}

func incrNum(num *int, mutex *sync.Mutex) {
	mutex.Lock()
	*num = *num + 1
	mutex.Unlock()
}

func rdpConn(ip, domain, user, password string, port int, timeout int64) (bool, error) {
	target := fmt.Sprintf("%s:%d", ip, port)
	g := newClient(target, glog.NONE)
	err := g.Login(domain, user, password, timeout)

	if err == nil {
		return true, nil
	}

	return false, err
}

type Client struct {
	Host string // ip:port
	tpkt *tpkt.TPKT
	x224 *x224.X224
	mcs  *t125.MCSClient
	sec  *sec.Client
	pdu  *pdu.Client
	vnc  *rfb.RFB
}

func newClient(host string, logLevel glog.LEVEL) *Client {
	glog.SetLevel(logLevel)
	logger := log.New(os.Stdout, "", 0)
	glog.SetLogger(logger)
	return &Client{
		Host: host,
	}
}

func (g *Client) Login(domain, user, pwd string, timeout int64) error {
	conn, err := xhttp.WrapperTcpWithTimeout("tcp", g.Host, time.Duration(timeout)*time.Second)
	if err != nil {
		return fmt.Errorf("[dial err] %v", err)
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)
	glog.Info(conn.LocalAddr().String())

	g.tpkt = tpkt.New(core.NewSocketLayer(conn), nla.NewNTLMv2(domain, user, pwd))
	g.x224 = x224.New(g.tpkt)
	g.mcs = t125.NewMCSClient(g.x224)
	g.sec = sec.NewClient(g.mcs)
	g.pdu = pdu.NewClient(g.sec)

	g.sec.SetUser(user)
	g.sec.SetPwd(pwd)
	g.sec.SetDomain(domain)

	g.tpkt.SetFastPathListener(g.sec)
	g.sec.SetFastPathListener(g.pdu)
	g.pdu.SetFastPathSender(g.tpkt)

	err = g.x224.Connect()
	if err != nil {
		return fmt.Errorf("[x224 connect err] %v", err)
	}
	glog.Info("wait connect ok")
	wg := &sync.WaitGroup{}
	breakFlag := false
	wg.Add(1)

	g.pdu.On("error", func(e error) {
		err = e
		glog.Error("error", e)
		g.pdu.Emit("done")
	})
	g.pdu.On("close", func() {
		err = errors.New("close")
		glog.Info("on close")
		g.pdu.Emit("done")
	})
	g.pdu.On("success", func() {
		err = nil
		glog.Info("on success")
		g.pdu.Emit("done")
	})
	g.pdu.On("ready", func() {
		glog.Info("on ready")
		g.pdu.Emit("done")
	})
	g.pdu.On("update", func(rectangles []pdu.BitmapData) {
		glog.Info("on update:", rectangles)
	})
	g.pdu.On("done", func() {
		if breakFlag == false {
			breakFlag = true
			wg.Done()
		}
	})
	wg.Wait()
	return err
}
