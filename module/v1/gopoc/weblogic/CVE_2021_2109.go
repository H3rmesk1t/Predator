package weblogic

import (
	"Predator/pkg/utils"
	"strings"
)

func CVE_2021_2109(url string) bool {
	if req, err := utils.HttpRequset(url+"/console/css/%252e%252e%252f/consolejndi.portal", "GET", "", false, nil); err == nil {
		if req.StatusCode == 200 && strings.Contains(req.Body, "Weblogic") {
			return true
		}
	}
	return false
}
