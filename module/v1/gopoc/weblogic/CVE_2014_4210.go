package weblogic

import (
	"Predator/pkg/utils"
)

func CVE_2014_4210(url string) bool {
	if req, err := utils.HttpRequset(url+"/uddiexplorer/SearchPublicRegistries.jsp", "GET", "", false, nil); err == nil {
		if req.StatusCode == 200 {
			return true
		}
	}
	return false
}
