package weblogic

import (
	"Predator/pkg/utils"
)

func CVE_2018_2894(url string) bool {
	if req, err := utils.HttpRequset(url+"/ws_utc/begin.do", "GET", "", false, nil); err == nil {
		if req2, err2 := utils.HttpRequset(url+"/ws_utc/config.do", "GET", "", false, nil); err2 == nil {
			if req.StatusCode == 200 || req2.StatusCode == 200 {
				return true
			}
		}
	}
	return false
}
