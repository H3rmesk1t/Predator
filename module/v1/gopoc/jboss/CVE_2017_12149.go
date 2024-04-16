package jboss

import (
	"Predator/pkg/utils"
)

func CVE_2017_12149(url string) bool {
	if req, err := utils.HttpRequset(url+"/invoker/readonly", "GET", "", false, nil); err == nil {
		if req.StatusCode == 500 {
			return true
		}
	}
	return false
}
