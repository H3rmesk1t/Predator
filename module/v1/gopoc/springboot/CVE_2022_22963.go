package springboot

import (
	"Predator/pkg/utils"
	"strings"
)

func CVE_2021_22963(u string) bool {
	var header = make(map[string]string, 1)
	header["Content-Type"] = "application/x-www-form-urlencoded"
	if req, err := utils.HttpRequset(u+"/functionRouter", "POST", "test", false, header); err == nil {
		if req.StatusCode == 500 && strings.Contains(req.Body, "Internal Server Error") {
			return true
		}
	}

	return false
}
