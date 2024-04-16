package utils

import (
	"Predator/pkg/config"
	"fmt"
	"github.com/fatih/color"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

var Num int64
var End int64
var Results = make(chan *string)
var LogSucTime int64
var LogErrTime int64
var LogWG sync.WaitGroup

type JsonText struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func init() {
	log.SetOutput(io.Discard)
	LogSucTime = time.Now().Unix()
	go SaveLog()
}

func LogSuccess(result string) {
	LogWG.Add(1)
	LogSucTime = time.Now().Unix()
	Results <- &result
}

func SaveLog() {
	for result := range Results {
		if !config.Silent {
			if config.NoColor {
				fmt.Println(*result)
			} else {
				if strings.HasPrefix(*result, "[+]") {
					color.HiGreen(*result)
				} else if strings.HasPrefix(*result, "[>]") {
					color.HiCyan(*result)
				} else if strings.HasPrefix(*result, "[~]") {
					color.HiBlue(*result)
				} else if strings.HasPrefix(*result, "[!]") {
					color.HiRed(*result)
				} else if strings.HasPrefix(*result, "[*] WebScan") {
					color.HiYellow(*result)
				} else {
					fmt.Println(*result)
				}
			}
		}
		if config.IsSave {
			WriteFile(*result, config.OutputFile)
		}
		LogWG.Done()
	}
}

func WriteFile(result string, filename string) {
	fl, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("[-] Open %s error, %v\n", filename, err)
		return
	}
	_, err = fl.Write([]byte(result + "\n"))
	if err != nil {
		fmt.Printf("[-] Write %s error, %v\n", filename, err)
	}
	err = fl.Close()
	if err != nil {
		fmt.Printf("[-] Write %s error, %v\n", filename, err)
	}
}

func LogError(errinfo interface{}) {
	if config.WaitTime == 0 {
		fmt.Printf("[-] finish %v/%v %v \n", End, Num, errinfo)
	} else if (time.Now().Unix()-LogSucTime) > config.WaitTime && (time.Now().Unix()-LogErrTime) > config.WaitTime {
		fmt.Printf("[-] finish %v/%v %v \n", End, Num, errinfo)
		LogErrTime = time.Now().Unix()
	}
}

func CheckErrs(err error) bool {
	if err == nil {
		return false
	}
	errs := []string{
		"closed by the remote host", "too many connections",
		"i/o timeout", "EOF", "A connection attempt failed",
		"established connection failed", "connection attempt failed",
		"Unable to read", "is not allowed to connect to this",
		"no pg_hba.conf entry",
		"No connection could be made",
		"invalid packet size",
		"bad connection",
	}
	for _, key := range errs {
		if strings.Contains(strings.ToLower(err.Error()), strings.ToLower(key)) {
			return true
		}
	}
	return false
}
